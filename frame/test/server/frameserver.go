package main

import (
	"fmt"
	"net/http"

	_ "net/http/pprof"

	"github.com/0x00b/gobbq/components/game"
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/frame"
	"github.com/0x00b/gobbq/frame/frameproto"
	"github.com/0x00b/gobbq/frame/test/testpb"
	"github.com/0x00b/gobbq/tool/secure"
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

	if !f.tempFrameSvr.Invalid() {
		return &testpb.StartFrameRsp{FrameSvr: uint64(f.tempFrameSvr)}, nil
	}

	echoClient, err := frameproto.NewFrameSever(c)
	if err != nil {
		return nil, err
	}

	_, err = echoClient.Init(c, &frameproto.InitReq{
		PlayerNum: 2,
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

	secure.GO(func() {
		fmt.Println("pprof start...")
		fmt.Println(http.ListenAndServe(":9877", nil))
	})

	xlog.Init("trace", true, true, &lumberjack.Logger{
		Filename:  "./server.log",
		MaxAge:    7,
		LocalTime: true,
	}, xlog.DefaultLogTag{})

	var g = game.NewGame()

	testpb.RegisterFrameService(g.EntityMgr, &FrameService{})
	frameproto.RegisterFrameSeverEntity(g.EntityMgr, &frame.FrameSever{})

	g.Serve()
}
