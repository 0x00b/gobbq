package main

import (
	"fmt"

	"github.com/0x00b/gobbq/components/game"
	"github.com/0x00b/gobbq/conf"
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/frame"
	"github.com/0x00b/gobbq/frame/frameproto"
	"github.com/0x00b/gobbq/frame/test/testpb"
	"github.com/0x00b/gobbq/xlog"
	"gopkg.in/natefinch/lumberjack.v2"
)

// FrameService
type FrameService struct {
	entity.Service

	tempFrameSvr entity.EntityID
}

// StartFrame
func (f *FrameService) StartFrame(c entity.Context, req *testpb.StartFrameReq) (*testpb.StartFrameRsp, error) {

	if f.tempFrameSvr.Invalid() {
		return &testpb.StartFrameRsp{FrameSvr: uint64(f.tempFrameSvr)}, nil
	}

	echoClient := frameproto.NewFrameSeverEntity(c)

	_, err := echoClient.Init(c, &frameproto.InitReq{
		ClientNum: 2,
	})

	if err != nil {
		xlog.Println("new frame server:", err)
		return nil, err
	}

	xlog.Println("new frame server:", echoClient)

	f.tempFrameSvr = echoClient.EntityID

	return &testpb.StartFrameRsp{FrameSvr: uint64(echoClient.EntityID)}, nil
}

func main() {

	xlog.Init("info", true, true, &lumberjack.Logger{
		Filename:  "./server.log",
		MaxAge:    7,
		LocalTime: true,
	}, xlog.DefaultLogTag{})

	fmt.Println(conf.C)

	var g = game.NewGame()

	testpb.RegisterFrameService(g.EntityMgr, &FrameService{})
	frameproto.RegisterFrameSeverEntity(g.EntityMgr, &frame.FrameSeverEntity{})

	g.Serve()
}
