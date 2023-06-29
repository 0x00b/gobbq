package frame

import "github.com/0x00b/gobbq/frame/frameproto"

type frameData struct {
	idx  uint32
	cmds []*frameproto.FrameData
}

func newFrameData(index uint32) *frameData {
	f := &frameData{
		idx:  index,
		cmds: make([]*frameproto.FrameData, 0),
	}

	return f
}

type lockstep struct {
	frames     map[uint32]*frameData
	frameCount uint32
}

func newLockstep() *lockstep {
	l := &lockstep{
		frames: make(map[uint32]*frameData),
	}

	return l
}

func (l *lockstep) reset() {
	l.frames = make(map[uint32]*frameData)
	l.frameCount = 0
}

func (l *lockstep) getFrameCount() uint32 {
	return l.frameCount
}

func (l *lockstep) pushCmd(cmd *frameproto.FrameData) bool {
	f, ok := l.frames[l.frameCount]
	if !ok {
		f = newFrameData(l.frameCount)
		l.frames[l.frameCount] = f
	}

	// 检查是否同一帧发来两次同样的操作
	for _, v := range f.cmds {
		if v.CLientID == cmd.CLientID {
			if v.Input == nil || v.Input.OPID == cmd.Input.OPID {
				return false
			}
		}
	}

	f.cmds = append(f.cmds, cmd)

	return true
}

func (l *lockstep) tick() uint32 {
	l.frameCount++
	return l.frameCount
}

// func (l *lockstep) getRangeFrames(from, to uint32) []*frameData {
// 	ret := make([]*frameData, 0, to-from)

// 	for ; from <= to && from <= l.frameCount; from++ {
// 		f, ok := l.frames[from]
// 		if !ok {
// 			continue
// 		}
// 		ret = append(ret, f)
// 	}

// 	return ret
// }

func (l *lockstep) getFrame(idx uint32) *frameData {

	return l.frames[idx]
}
