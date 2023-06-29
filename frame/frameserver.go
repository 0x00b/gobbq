package frame

import (
	"time"

	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/frame/frameproto"
	"github.com/0x00b/gobbq/xlog"
)

const (
	MaxReadyTime         int64  = 20          // second,准备阶段最长时间，如果超过这个时间没人连进来直接关闭游戏
	MaxGameFrame         uint32 = 30*60 + 100 // 每局最大帧数
	kMaxFrameDataPerMsg         = 60          // 每个消息包最多包含多少个帧数据
	kBadNetworkThreshold        = 2           // 这个时间段没有收到心跳包认为他网络很差，不再持续给发包(网络层的读写时间设置的比较长，客户端要求的方案)
)
const (
	Frequency = 15               // 每秒15帧
	TickTimer = 1000 / Frequency // 心跳Timer
)

// GameState 游戏状态
type GameState int

const (
	GameStop  GameState = 0 // 准备阶段
	GameReady GameState = 1 // 准备阶段
	GamRuning GameState = 2 // 战斗中阶段
	GameOver  GameState = 3 // 结束阶段
)

// FrameSever
type FrameSever struct {
	entity.Entity

	startTime              int64
	lastBroadcastFrameTime int64

	entityNum int

	status GameState

	// key:UID
	players map[uint64]*Player

	logic *lockstep
}

func (f *FrameSever) OnTick() {
	if f.status == GameStop {
		return
	}

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

func (f *FrameSever) broadcastFrameData() {

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
		// i < framesCount 所以最后一帧是不会发出去的
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
				err := p.Frame(f.Context(), msg)
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
	for _, v := range f.players {
		if v.EntityID != wn.EntityID {
			continue
		}

		v.isOnline = false
		v.isReady = false
	}
}

// Heartbeat
func (f *FrameSever) Heartbeat(c entity.Context, req *frameproto.HeartbeatReq) (*frameproto.HeartbeatRsp, error) {

	xlog.Info("recv heart beat")

	p := f.players[req.UID]

	if p != nil {
		p.RefreshHeartbeatTime()
	}

	return &frameproto.HeartbeatRsp{}, nil
}

// Init
func (f *FrameSever) Init(c entity.Context, req *frameproto.InitReq) (*frameproto.InitRsp, error) {

	xlog.Info(c, "Init entityNum:", req.PlayerNum)

	f.entityNum = int(req.PlayerNum)
	f.players = make(map[uint64]*Player, f.entityNum)
	f.logic = newLockstep()
	f.startTime = time.Now().UnixMilli()
	f.status = GameReady

	return &frameproto.InitRsp{}, nil
}

// Join
func (f *FrameSever) Join(c entity.Context, req *frameproto.JoinReq) (*frameproto.JoinRsp, error) {

	xlog.Info("recv join", req.String())

	// 目前断线重连这个id会变

	p, ok := f.players[req.UID]
	if ok {
		p.isOnline = true
		return &frameproto.JoinRsp{}, nil
	}

	id := c.SrcEntity()
	f.Watch(id)

	p = NewPlayer(id, req.UID)
	p.isOnline = true

	f.players[req.UID] = p

	return &frameproto.JoinRsp{}, nil
}

// Progress 无实际意义,告诉其他人加载进度
func (f *FrameSever) Progress(c entity.Context, req *frameproto.ProgressReq) error {

	for _, v := range f.players {
		if v.UID == req.UID {
			continue
		}

		err := v.Progress(f.Context(), req)
		if err != nil {
			xlog.Errorln(err)
			// return nil, err
		}
	}
	return nil
}

func (f *FrameSever) doReady(p *Player) {

	if p.isReady {
		return
	}
	p.isReady = true
}

// Ready
func (f *FrameSever) Ready(c entity.Context, req *frameproto.ReadyReq) error {

	p := f.players[req.UID]
	if p == nil {
		return nil
	}

	f.doReady(p)

	// 重连
	if f.status == GamRuning {
		p.EntityID = c.SrcEntity()
		f.Watch(p.EntityID)
		f.doReconnect(p)
	}

	return nil

}

// Input
func (f *FrameSever) Input(c entity.Context, req *frameproto.InputReq) error {

	p := f.players[req.UID]
	if p == nil {
		return nil
	}

	d := &frameproto.FrameData{
		UID:   uint64(req.UID),
		Input: req.Input,
	}

	f.logic.pushCmd(d)

	return nil

}

// GameOver 上报结束
func (f *FrameSever) GameOver(c entity.Context, req *frameproto.GameOverReq) error {

	p := f.players[req.UID]
	if p == nil {
		return nil
	}

	p.result = req.Result

	return nil
}

func (f *FrameSever) checkReady() bool {

	cnt := 0
	for _, v := range f.players {
		if !v.isReady {
			return false
		}
		cnt++
	}

	return cnt >= f.entityNum
}

func (f *FrameSever) doStart() {

	// f.clientFrameCount = 0
	f.logic.reset()
	for _, v := range f.players {
		v.isReady = true
		v.loadingProgress = 100
	}

	f.status = GamRuning

	for _, v := range f.players {

		err := v.Start(f.Context(), &frameproto.StartReq{})
		if err != nil {
			xlog.Errorln(err)
			// return nil, err
		}
	}
}

func (f *FrameSever) doGameOver() {
	for _, v := range f.players {

		err := v.GameOver(f.Context(), &frameproto.GameOverReq{})
		if err != nil {
			xlog.Errorln(err)
			// return nil, err
		}
	}
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

// 可以按照超过一半人结束就结束游戏
func (f *FrameSever) checkOver() bool {
	// 只要有人没发结果并且还在线，就不结束
	for _, v := range f.players {
		if !v.isOnline {
			continue
		}

		if v.result == 0 {
			return false
		}
	}

	return true
}

func (f *FrameSever) isTimeout() bool {
	return f.logic.getFrameCount() > MaxGameFrame

}

func (f *FrameSever) doReconnect(p *Player) {

	// 先start
	err := p.Start(f.Context(), &frameproto.StartReq{})
	if err != nil {
		xlog.Errorln(err)
		// return nil, err
	}

	framesCount := f.logic.getFrameCount()
	var i uint32 = 0
	c := 0

	frameMsg := &frameproto.FrameReq{}

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

		frameMsg.Frames = append(frameMsg.Frames, fd)
		c++

		// 如果是最后一帧或者达到这个消息包能装下的最大帧数，就发送
		if i == (framesCount-1) || c >= kMaxFrameDataPerMsg {
			err := p.Frame(f.Context(), frameMsg)
			if err != nil {
				xlog.Errorln(err)
			}
			c = 0
			frameMsg = &frameproto.FrameReq{}
		}
	}

	p.SetSendFrameCount(framesCount)

}
