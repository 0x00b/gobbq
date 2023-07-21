// NOTE:!!
// DON'T EDIT IT!!
//  请勿添加其他内容，包括函数，结构体，变量等等，否则重新生成时会丢失。

package exampb

import (
	empty "github.com/golang/protobuf/ptypes/empty"

	"errors"
	"context"

	"github.com/0x00b/gobbq/engine/db"
	"github.com/0x00b/gobbq/engine/model"
	//"github.com/0x00b/gobbq/proto/bbq"
	//pb "google.golang.org/protobuf/proto"
	// exampb "github.com/0x00b/gobbq/example/exampb"
)

func (m *EchoProperty) TableName() string {
	return EchoProperty_TableName
}

type EchoPropertyModel struct {
	EchoProperty

	watchers []model.Watcher
	db       db.IDatabase
}

const (
	EchoProperty_TableName string = "EchoProperty"

	EchoProperty_Text  model.FieldName = "Text"
	EchoProperty_Test  model.FieldName = "Test"
	EchoProperty_Test2 model.FieldName = "Test2"
	EchoProperty_Test3 model.FieldName = "Test3"
	EchoProperty_Test4 model.FieldName = "Test4"
	EchoProperty_Test5 model.FieldName = "Test5"
	EchoProperty_Test6 model.FieldName = "Test6"
	EchoProperty_TEST7 model.FieldName = "TEST7"
	EchoProperty_Test8 model.FieldName = "Test8"
	EchoProperty_Test9 model.FieldName = "Test9"
)

func (m *EchoPropertyModel) ModelInit(c context.Context, db db.IDatabase) error {
	if m == nil {
		return errors.New("nil model")
	}
	db, err := db.Table(EchoProperty_TableName)
	if err != nil {
		return err
	}
	m.db = db
	m.ModelAddWatcher(c, db)

	return m.ModelLoad(c)
}

func (m *EchoPropertyModel) Destroy(c context.Context) {
	m.ModelStopWatcher(c)
}

func (m *EchoPropertyModel) ModelLoad(c context.Context) error {
	if m == nil {
		return errors.New("nil model")
	}
	err := m.db.Load(c, &m.EchoProperty)
	if err != nil {
		return err
	}
	return nil
}

func (m *EchoPropertyModel) ModelAutoSave(c context.Context) error {
	if m == nil {
		return nil
	}
	if m.db == nil {
		return nil
	}

	return m.db.AutoSave(c, &m.EchoProperty)
}

func (m *EchoPropertyModel) ModelSave(c context.Context) error {
	if m == nil {
		return nil
	}
	if m.db == nil {
		return nil
	}

	return m.db.Save(c, &m.EchoProperty)
}

func (m *EchoPropertyModel) ModelAddWatcher(c context.Context, s model.Watcher) {
	if m == nil {
		return
	}
	m.watchers = append(m.watchers, s)
	s.Watch(c, &m.EchoProperty)
}

func (m *EchoPropertyModel) ModelStopWatcher(c context.Context) {
	if m == nil {
		return
	}
	for _, w := range m.watchers {
		w.QuitWatch(c, &m.EchoProperty)
	}
}

func (m *EchoPropertyModel) SetText(v string) {
	if m == nil {
		return
	}
	m.Text = v

	for _, w := range m.watchers {
		w.SetNotify(&m.EchoProperty, 1, EchoProperty_Text)
	}
}

func (m *EchoPropertyModel) SetTest(v *SayHelloRequest) {
	if m == nil {
		return
	}
	m.Test = v

	for _, w := range m.watchers {
		w.SetNotify(&m.EchoProperty, 2, EchoProperty_Test)
	}
}

func (m *EchoPropertyModel) SetTest2(v []*empty.Empty) {
	if m == nil {
		return
	}
	m.Test2 = v

	for _, w := range m.watchers {
		w.SetNotify(&m.EchoProperty, 3, EchoProperty_Test2)
	}
}

func (m *EchoPropertyModel) SetTest3(v map[int32]string) {
	if m == nil {
		return
	}
	m.Test3 = v

	for _, w := range m.watchers {
		w.SetNotify(&m.EchoProperty, 4, EchoProperty_Test3)
	}
}

func (m *EchoPropertyModel) SetTest4(v []byte) {
	if m == nil {
		return
	}
	m.Test4 = v

	for _, w := range m.watchers {
		w.SetNotify(&m.EchoProperty, 5, EchoProperty_Test4)
	}
}

func (m *EchoPropertyModel) SetTest5(v int64) {
	if m == nil {
		return
	}
	m.Test5 = v

	for _, w := range m.watchers {
		w.SetNotify(&m.EchoProperty, 6, EchoProperty_Test5)
	}
}

func (m *EchoPropertyModel) SetTest6(v int32) {
	if m == nil {
		return
	}
	m.Test6 = v

	for _, w := range m.watchers {
		w.SetNotify(&m.EchoProperty, 7, EchoProperty_Test6)
	}
}

func (m *EchoPropertyModel) SetTEST7(v Enum) {
	if m == nil {
		return
	}
	m.TEST7 = v

	for _, w := range m.watchers {
		w.SetNotify(&m.EchoProperty, 8, EchoProperty_TEST7)
	}
}

func (m *EchoPropertyModel) SetTest8(v int32) {
	if m == nil {
		return
	}
	m.Test8 = v

	for _, w := range m.watchers {
		w.SetNotify(&m.EchoProperty, 9, EchoProperty_Test8)
	}
}

func (m *EchoPropertyModel) SetTest9(v []int32) {
	if m == nil {
		return
	}
	m.Test9 = v

	for _, w := range m.watchers {
		w.SetNotify(&m.EchoProperty, 10, EchoProperty_Test9)
	}
}
