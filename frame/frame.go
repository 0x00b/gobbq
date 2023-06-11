package frame

import (
	"time"

	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/frame/frameproto"
	"github.com/0x00b/gobbq/xlog"
)

// FrameSeverEntity
type FrameSeverEntity struct {
	entity.Entity

	entityNum int
	curNum    int
	entities  []entity.EntityID

	frameData []*frameproto.FrameData
}

// Heartbeat
func (f *FrameSeverEntity) Heartbeat(c entity.Context, req *frameproto.HeartbeatReq) error {

	xlog.Info("recv heart beat")

	return nil
}

// Init
func (f *FrameSeverEntity) Init(c entity.Context, req *frameproto.InitReq) (*frameproto.InitRsp, error) {

	f.entityNum = int(req.GetClientNum())
	xlog.Info(c, "Init entityNum:", f.entityNum)

	return &frameproto.InitRsp{}, nil
}

func (f *FrameSeverEntity) broadcastFrameData() {

	xlog.Info("broadcastFrameData", f.frameData)

	for _, v := range f.entities {
		client := frameproto.NewFrameClientEntityClient(v)
		err := client.Frame(f.Context(), &frameproto.FrameReq{
			Data: f.frameData,
		})
		if err != nil {
			xlog.Errorln(err)
		}
	}

}

// Join
func (f *FrameSeverEntity) Join(c entity.Context, req *frameproto.JoinReq) (*frameproto.JoinRsp, error) {

	xlog.Info("recv join", req.String())

	f.entities = append(f.entities, c.SrcEntity())
	f.curNum++

	xlog.Info("recv join", f.curNum, f.entityNum)

	if f.curNum >= f.entityNum {

		for _, v := range f.entities {
			client := frameproto.NewFrameClientEntityClient(v)
			err := client.Start(c, &frameproto.StartReq{})
			if err != nil {
				xlog.Errorln(err)
				// return nil, err
			}
		}
		f.AddTimer(1000/15*time.Millisecond, f.broadcastFrameData)
	}

	return &frameproto.JoinRsp{}, nil
}

// Progress
func (f *FrameSeverEntity) Progress(c entity.Context, req *frameproto.ProgressReq) (*frameproto.ProgressRsp, error) {
	return &frameproto.ProgressRsp{}, nil
}

// Ready
func (f *FrameSeverEntity) Ready(c entity.Context, req *frameproto.ReadyReq) error {
	return nil

}

// Move
func (f *FrameSeverEntity) Move(c entity.Context, req *frameproto.MoveReq) error {

	client := c.SrcEntity()

	f.frameData = append(f.frameData, &frameproto.FrameData{
		CLientID: uint64(client),
		Pos:      req.GetPos(),
		Data:     []*frameproto.InputData{},
	})

	return nil
}

// Input
func (f *FrameSeverEntity) Input(c entity.Context, req *frameproto.InputReq) error {
	client := c.SrcEntity()

	f.frameData = append(f.frameData, &frameproto.FrameData{
		CLientID: uint64(client),
		// Pos:      req.Pos,
		Data: []*frameproto.InputData{req.GetData()},
	})

	return nil

}

// Result
func (f *FrameSeverEntity) Result(c entity.Context, req *frameproto.ResultReq) error {

	return nil
}
