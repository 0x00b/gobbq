package frame

import (
	"time"

	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/frame/frameproto"
)

type Player struct {
	*frameproto.FrameClient

	UID uint64

	isReady           bool
	isOnline          bool
	loadingProgress   int32
	lastHeartbeatTime int64
	sendFrameCount    uint32

	result int32
}

func NewPlayer(eid entity.EntityID, uid uint64) *Player {
	p := &Player{
		FrameClient: frameproto.NewFrameClientClient(eid),
		UID:         uid,
	}

	return p
}

func (p *Player) IsOnline() bool {
	return p.isOnline
}

func (p *Player) LoadingProgress() int32 {
	return p.loadingProgress
}

func (p *Player) SetLoadingProgress(n int32) {
	p.loadingProgress = n
}

func (p *Player) RefreshHeartbeatTime() {
	p.lastHeartbeatTime = time.Now().Unix()
}

func (p *Player) GetLastHeartbeatTime() int64 {
	return p.lastHeartbeatTime
}

func (p *Player) SetSendFrameCount(c uint32) {
	p.sendFrameCount = c
}

func (p *Player) GetSendFrameCount() uint32 {
	return p.sendFrameCount
}

func (p *Player) Cleanup() {

	p.isReady = false
	p.isOnline = false
}
