package mongo

import (
	"context"
	"fmt"

	"github.com/0x00b/gobbq/engine/model"
	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/proto"
)

var _ model.Watcher = &mongoDB{}

type watchMessage struct {
	msg    proto.Message
	fileds map[protowire.Number]model.FieldName
}

// Watch 订阅这个Message的字段变更
func (m *mongoDB) Watch(c context.Context, msg proto.Message) {

	wm, ok := m.watchModels[msg]
	if ok && wm != nil {
		return
	}

	m.watchModels[msg] = &watchMessage{
		msg:    msg,
		fileds: make(map[protowire.Number]model.FieldName),
	}

}

// SetNotify Message字段变更时调用，idx是pb中定义的字段编号
func (m *mongoDB) SetNotify(msg proto.Message, n protowire.Number, name model.FieldName) {

	wm, ok := m.watchModels[msg]
	if !ok || wm == nil {
		fmt.Println("update 3ew:")
		return
	}

	wm.fileds[n] = name
}

// QuitWatch Message不希望被订阅时调用
func (m *mongoDB) QuitWatch(c context.Context, msg proto.Message) {

	wm, ok := m.watchModels[msg]
	if !ok || wm == nil {
		return
	}
	// todo save
	m.updateDirtyField(c, msg)

	m.watchModels[msg] = nil
}

func (m *mongoDB) updateDirtyField(c context.Context, msg proto.Message) error {

	wm, ok := m.watchModels[msg]
	if !ok || wm == nil {
		return nil
	}
	fields := []model.FieldName{}
	for _, v := range wm.fileds {
		fields = append(fields, v)
	}

	// clear dirty

	return m.Update(c, msg, fields)
}
