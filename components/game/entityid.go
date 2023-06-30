package game

import "github.com/0x00b/gobbq/engine/entity"

type defaultIDGener struct {
}

// NewEntityID 如果没有特殊规划，可以使用这个生成entity id
func (g *defaultIDGener) GenID() entity.ID {
	return entity.GenID()
}
