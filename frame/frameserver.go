package frame

import (
	"time"

	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/frame/frameproto"
	"github.com/0x00b/gobbq/xlog"
)

const (
	BroadcastFrameCnt = 3                      // 每秒45帧， 3帧一次下发, 也就是1秒15次下发，大概66毫秒一次
	FrameFrequency    = 15 * BroadcastFrameCnt // 每秒45帧， 3帧一次下发, 也就是1秒15次下发，大概66毫秒一次
	FrameTimer        = 1000 / FrameFrequency  // 心跳Timer

	MaxReadyTime        int64  = 20                         // second,准备阶段最长时间，如果超过这个时间没人连进来直接关闭游戏
	MaxGameFrame        uint32 = 30*60*FrameFrequency + 100 // 每局最大帧数(30分钟)
	MaxFrameDataPerMsg         = 60                         // 每个消息包最多包含多少个帧数据
	BadNetworkThreshold        = 2                          // 这个时间段没有收到心跳包认为他网络很差，不再持续给发包(网络层的读写时间设置的比较长，客户端要求的方案)

)

// GameState 游戏状态
type GameState int

const (
	GameNone  GameState = 0 // 准备阶段
	GameReady GameState = 1 // 准备阶段
	GamRuning GameState = 2 // 战斗中阶段
	GameOver  GameState = 3 // 结束阶段
	GameStop  GameState = 4 // 准备阶段
)

// FrameSever
type FrameSever struct {
	entity.Entity

	startTime              int64
	lastBroadcastFrameTime int64
	clientFrameCnt         uint32

	entityNum int

	status GameState

	// key:Id
	players map[uint64]*Player

	logic *lockstep
}

func (f *FrameSever) OnInit() {
	f.SetTickIntervel(1 * time.Millisecond)
}

func (f *FrameSever) OnTick() {
	if f.status == GameNone {
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

		if now-f.lastBroadcastFrameTime < FrameTimer {
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
		f.Stop()
		return
	}

}

func (f *FrameSever) broadcastFrameData() {

	framesCount := f.logic.getFrameCount()

	if framesCount-f.clientFrameCnt < BroadcastFrameCnt {
		return
	}

	defer func() {
		f.clientFrameCnt = framesCount
	}()

	now := time.Now().Unix()

	for _, p := range f.players {

		// 掉线的
		if !p.IsOnline() {
			continue
		}

		if !p.isReady {
			continue
		}

		// 网络不好的,超过一定时间没收到心跳
		if now-p.GetLastHeartbeatTime() >= BadNetworkThreshold {
			continue
		}

		// 获得这个玩家已经发到哪一帧
		i := p.GetSendFrameCount()

		f.sendFrame(p, i, framesCount)
	}
}

func (f *FrameSever) sendFrame(p *Player, start, end uint32) {

	c := 0

	msg := &frameproto.FrameReq{}

	for i := start; i < end; i++ {

		frameData := f.logic.getFrame(i)
		if nil == frameData && i != (end-1) {
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
		if i == (end-1) || c >= MaxFrameDataPerMsg {
			err := p.Frame(f.Context(), msg)
			if err != nil {
				xlog.Errorln(err)
			}
			c = 0
			msg = &frameproto.FrameReq{}
		}
	}

	p.SetSendFrameCount(end)

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

	p := f.players[req.Id]

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
	f.clientFrameCnt = 0

	return &frameproto.InitRsp{}, nil
}

// Join
func (f *FrameSever) Join(c entity.Context, req *frameproto.JoinReq) (*frameproto.JoinRsp, error) {

	xlog.Info("recv join", req.String())

	p, ok := f.players[req.Id]
	if ok {
		p.isOnline = true
		return &frameproto.JoinRsp{}, nil
	}

	id := c.SrcEntity()
	f.Watch(id)

	p = NewPlayer(id, req.Id)
	p.isOnline = true

	f.players[req.Id] = p

	return &frameproto.JoinRsp{}, nil
}

// Progress 无实际意义,告诉其他人加载进度
func (f *FrameSever) Progress(c entity.Context, req *frameproto.ProgressReq) error {

	for _, v := range f.players {
		if v.Id == req.Id {
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

	p := f.players[req.Id]
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

	p := f.players[req.Id]
	if p == nil {
		return nil
	}

	d := &frameproto.FrameData{
		Id:    uint64(req.Id),
		Input: req.Input,
	}

	f.logic.pushCmd(d)

	return nil

}

// GameOver 上报结束
func (f *FrameSever) GameOver(c entity.Context, req *frameproto.GameOverReq) error {

	p := f.players[req.Id]
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
		// return
	}

	f.sendFrame(p, 0, f.clientFrameCnt)
}
