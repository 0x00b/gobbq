package frame

import (
	"time"

	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/frame/frameproto"
	"github.com/0x00b/gobbq/xlog"
)

const (
	MaxReadyTime         int64  = 20            // 准备阶段最长时间，如果超过这个时间没人连进来直接关闭游戏
	MaxGameFrame         uint32 = 30*60*3 + 100 // 每局最大帧数
	kMaxFrameDataPerMsg         = 60            // 每个消息包最多包含多少个帧数据
	kBadNetworkThreshold        = 2             // 这个时间段没有收到心跳包认为他网络很差，不再持续给发包(网络层的读写时间设置的比较长，客户端要求的方案)
)
const (
	Frequency = 15               // 每秒15帧
	TickTimer = 1000 / Frequency // 心跳Timer
)

// GameState 游戏状态
type GameState int

const (
	GameReady GameState = 0 // 准备阶段
	GamRuning GameState = 1 // 战斗中阶段
	GameOver  GameState = 2 // 结束阶段
	GameStop  GameState = 3 // 停止
)

// FrameSever
type FrameSever struct {
	entity.Entity

	startTime              int64
	lastBroadcastFrameTime int64

	entityNum int
	curNum    int

	status GameState

	// result map[uint64]uint64

	players map[entity.EntityID]*Player

	logic *lockstep
}

func (f *FrameSever) OnTick() {

	now := time.Now().UnixMilli()

	switch f.status {
	case GameReady:
		delta := (now - f.startTime) / 1000
		if delta < MaxReadyTime {
			if f.checkReady() {
				f.doStart()
			}
			return
		}

		if f.getOnlinePlayerCount() > 0 {
			// 大于最大准备时间，只要有在线的，就强制开始
			f.doStart()
			xlog.Warn("[game(%d)] force start game because ready state is timeout ", f.EntityID())
			return
		}

		// 全都没连进来，直接结束
		f.status = GameOver
		xlog.Error("[game(%d)] game over!! nobody ready", f.EntityID())

		return
	case GamRuning:
		if f.checkOver() {
			f.status = GameOver
			xlog.Info("[game(%d)] game over successfully!!", f.EntityID())
			return
		}

		if f.isTimeout() {
			f.status = GameOver
			xlog.Warn("[game(%d)] game timeout", f.EntityID())
			return
		}

		if now-f.lastBroadcastFrameTime < TickTimer {
			return
		}

		f.lastBroadcastFrameTime = now

		f.logic.tick()
		f.broadcastFrameData()

		return
	case GameOver:
		f.doGameOver()
		f.status = GameStop
		xlog.Info("[game(%d)] do game over", f.EntityID())
		return
	case GameStop:
		return
	}

}

// Heartbeat
func (f *FrameSever) Heartbeat(c entity.Context, req *frameproto.HeartbeatReq) error {

	xlog.Info("recv heart beat")

	id := c.SrcEntity()

	p := f.players[id]

	p.RefreshHeartbeatTime()

	return nil
}

// Init
func (f *FrameSever) Init(c entity.Context, req *frameproto.InitReq) (*frameproto.InitRsp, error) {

	xlog.Info(c, "Init entityNum:", f.entityNum)

	f.entityNum = int(req.GetClientNum())
	f.players = make(map[entity.EntityID]*Player, f.entityNum)
	f.logic = newLockstep()

	return &frameproto.InitRsp{}, nil
}

func (f *FrameSever) broadcastFrameData() {

	if f.status == GameStop {

		xlog.Info("stoped")
		return
	}

	// now := time.Now().Unix()

	framesCount := f.logic.getFrameCount()
	for _, p := range f.players {

		// 掉线的
		if !p.IsOnline() {
			continue
		}

		if !p.isReady {
			continue
		}

		// 网络不好的
		// if now-p.GetLastHeartbeatTime() >= kBadNetworkThreshold {
		// 	continue
		// }

		// 获得这个玩家已经发到哪一帧
		i := p.GetSendFrameCount()
		c := 0

		msg := &frameproto.FrameReq{}
		for ; i < framesCount; i++ {
			frameData := f.logic.getFrame(i)
			if nil == frameData && i != (framesCount-1) {
				continue
			}

			fd := &frameproto.Frame{
				FrameID: i,
			}

			if nil != frameData {
				fd.Data = frameData.cmds
			}

			msg.Frames = append(msg.Frames, fd)
			c++

			// 如果是最后一帧或者达到这个消息包能装下的最大帧数，就发送
			if i == (framesCount-1) || c >= kMaxFrameDataPerMsg {
				client := frameproto.NewFrameClientEntityClient(p.clientID)
				err := client.Frame(f.Context(), msg)
				if err != nil {
					xlog.Errorln(err)
				}
				c = 0
				msg = &frameproto.FrameReq{}
			}
		}

		p.SetSendFrameCount(framesCount)

	}

}

func (f *FrameSever) OnNotify(wn entity.NotifyInfo) {
	xlog.Infoln("receive watch client notify...", wn.EntityID)
	// 断线了

	p := f.players[wn.EntityID]

	p.isOnline = false
	for _, v := range f.players {
		if v.IsOnline() {
			return
		}
	}
	// stoped
	f.status = 1
}

// Join
func (f *FrameSever) Join(c entity.Context, req *frameproto.JoinReq) (*frameproto.JoinRsp, error) {

	xlog.Info("recv join", req.String())

	// 目前断线重连这个id会变
	id := c.SrcEntity()

	p, ok := f.players[id]
	if ok {
		p.isOnline = true
		// 重连

		return &frameproto.JoinRsp{}, nil
	}

	f.Watch(id)

	p = NewPlayer(id)
	p.isOnline = true
	p.isReady = true

	f.players[id] = p

	f.curNum++

	xlog.Info("recv join", f.curNum, f.entityNum)

	if f.curNum >= f.entityNum {
		f.doStart()
	}

	return &frameproto.JoinRsp{}, nil
}

// Progress
func (f *FrameSever) Progress(c entity.Context, req *frameproto.ProgressReq) (*frameproto.ProgressRsp, error) {
	return &frameproto.ProgressRsp{}, nil
}

// Ready
func (f *FrameSever) Ready(c entity.Context, req *frameproto.ReadyReq) error {

	f.doReady(f.players[c.SrcEntity()])

	return nil

}

// Input
func (f *FrameSever) Input(c entity.Context, req *frameproto.InputReq) error {
	client := c.SrcEntity()

	d := &frameproto.FrameData{
		CLientID: uint64(client),
		Input:    req.Input,
	}

	f.logic.pushCmd(d)

	return nil

}

// Result
func (f *FrameSever) Result(c entity.Context, req *frameproto.ResultReq) error {

	return nil
}

func (f *FrameSever) doReady(p *Player) {

	if p.isReady {
		return
	}

	p.isReady = true

	// msg := pb_packet.NewPacket(uint8(pb.ID_MSG_Ready), nil)
	// p.SendMessage(msg)
}

func (f *FrameSever) checkReady() bool {
	for _, v := range f.players {
		if !v.isReady {
			return false
		}
	}

	return true
}

func (f *FrameSever) doStart() {

	// f.clientFrameCount = 0
	f.logic.reset()
	for _, v := range f.players {
		v.isReady = true
		v.loadingProgress = 100
	}

	f.startTime = time.Now().UnixMilli()

	for _, v := range f.players {
		client := frameproto.NewFrameClientEntityClient(v.clientID)
		err := client.Start(f.Context(), &frameproto.StartReq{})
		if err != nil {
			xlog.Errorln(err)
			// return nil, err
		}
	}

	f.status = GamRuning

	// f.listener.OnGameStart(f.id)
}

func (f *FrameSever) doGameOver() {

	// f.OnGameOver(f.id)
}

func (f *FrameSever) getOnlinePlayerCount() int {

	i := 0
	for _, v := range f.players {
		if v.IsOnline() {
			i++
		}
	}

	return i
}

func (f *FrameSever) checkOver() bool {
	// 只要有人没发结果并且还在线，就不结束
	for _, v := range f.players {
		if !v.isOnline {
			continue
		}

		// if _, ok := f.result[v.id]; !ok {
		return false
		// }
	}

	return true
}

func (f *FrameSever) isTimeout() bool {
	return f.logic.getFrameCount() > MaxGameFrame

}
