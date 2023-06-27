package frame

import (
	"time"

	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/frame/frameproto"
	"github.com/0x00b/gobbq/xlog"
)

const (
	// MaxReadyTime          int64  = 20            // 准备阶段最长时间，如果超过这个时间没人连进来直接关闭游戏
	// MaxGameFrame          uint32 = 30*60*3 + 100 // 每局最大帧数
	// BroadcastOffsetFrames        = 3             // 每隔多少帧广播一次
	kMaxFrameDataPerMsg  = 60 // 每个消息包最多包含多少个帧数据
	kBadNetworkThreshold = 2  // 这个时间段没有收到心跳包认为他网络很差，不再持续给发包(网络层的读写时间设置的比较长，客户端要求的方案)
)
const (
	Frequency = 15                      // 每秒15帧
	TickTimer = time.Second / Frequency // 心跳Timer
)

// FrameSever
type FrameSever struct {
	entity.Entity

	entityNum int
	curNum    int

	status int

	players map[entity.EntityID]*Player

	logic *lockstep
}

func (f *FrameSever) OnTick() {

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
	f.logic.reset()

	return &frameproto.InitRsp{}, nil
}

func (f *FrameSever) broadcastFrameData() {

	xlog.Info("broadcastFrameData")
	if f.status == 1 {

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

func (f *FrameSever) EntityNotify(wn entity.NotifyInfo) {
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

		for _, v := range f.players {
			client := frameproto.NewFrameClientEntityClient(v.clientID)
			err := client.Start(c, &frameproto.StartReq{})
			if err != nil {
				xlog.Errorln(err)
				// return nil, err
			}
		}

		f.AddTimer(TickTimer, func() {
			f.logic.tick()
			f.broadcastFrameData()
		})

	}

	return &frameproto.JoinRsp{}, nil
}

// Progress
func (f *FrameSever) Progress(c entity.Context, req *frameproto.ProgressReq) (*frameproto.ProgressRsp, error) {
	return &frameproto.ProgressRsp{}, nil
}

// Ready
func (f *FrameSever) Ready(c entity.Context, req *frameproto.ReadyReq) error {
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
