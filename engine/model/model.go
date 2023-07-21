package model

import (
	"context"

	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// type Model interface {
// 	AddWatcher(Watcher)
// 	LoadFrom(db.IDatabase)
// }

type FieldName protoreflect.Name

type Watcher interface {
	// 订阅这个Message的字段变更
	Watch(c context.Context, msg proto.Message)
	// Message字段变更时调用，idx是pb中定义的字段编号
	SetNotify(msg proto.Message, n protowire.Number, name FieldName)
	// Message不希望被订阅时调用
	QuitWatch(c context.Context, msg proto.Message)
}
