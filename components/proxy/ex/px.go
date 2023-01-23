package ex

import (
	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/engine/entity"
)

type ProxyMap map[entity.EntityID]*codec.PacketReadWriter
