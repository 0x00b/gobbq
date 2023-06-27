package frame

import (
	"time"

	"github.com/0x00b/gobbq/engine/entity"
)

type Player struct {
	clientID          entity.EntityID
	isReady           bool
	isOnline          bool
	loadingProgress   int32
	lastHeartbeatTime int64
	sendFrameCount    uint32
}

func NewPlayer(id entity.EntityID) *Player {
	p := &Player{
		clientID: id,
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
