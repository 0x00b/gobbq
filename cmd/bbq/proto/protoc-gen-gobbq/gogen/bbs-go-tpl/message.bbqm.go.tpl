// NOTE:!!
// DON'T EDIT IT!!
//  请勿添加其他内容，包括函数，结构体，变量等等，否则重新生成时会丢失。

package {{$.GoPackageName}}

import (

{{with $.GoImplImportPaths}}
	{{range $idx, $path := $.GoImplImportPaths}} 
	{{- $path.Alias}} "{{$path.ImportPath -}}"
	{{end}}
{{end}}
	"errors"
	"context"

	"github.com/0x00b/gobbq/engine/db"
	"github.com/0x00b/gobbq/engine/model"
	//"github.com/0x00b/gobbq/proto/bbq"
	//pb "google.golang.org/protobuf/proto"
	// {{$.GoPackageName}} "{{$.GoImportPath}}"
)
 

{{range $mi, $m := $.Messages }}
{{with HasBBQFieldTag $m}}


func (m *{{$m.GoIdent.GoName}}) TableName() string {
	return {{$m.GoIdent.GoName}}_TableName
}

type {{$m.GoIdent.GoName}}Model struct {
	{{$m.GoIdent.GoName}}

	watchers []model.Watcher
	db       db.IDatabase
}

const (
	{{$m.GoIdent.GoName}}_TableName string = "{{$m.GoIdent.GoName}}"

{{range $fi, $f := $m.Fields -}}
	{{$m.GoIdent.GoName}}_{{$f.GoName}} model.FieldName = "{{$f.GoName}}"
{{end -}}
)

func (m *{{$m.GoIdent.GoName}}Model) ModelInit(c context.Context, db db.IDatabase) error {
	if m == nil {
		return errors.New("nil model")
	}
	db, err := db.Table({{$m.GoIdent.GoName}}_TableName)
	if err != nil {
		return err
	}
	m.db = db
	m.ModelAddWatcher(c, db)

	return m.ModelLoad(c)
}

func (m *{{$m.GoIdent.GoName}}Model) Destroy(c context.Context) {
	m.ModelStopWatcher(c)
}

func (m *{{$m.GoIdent.GoName}}Model) ModelLoad(c context.Context) error {
	if m == nil {
		return errors.New("nil model")
	}
	err := m.db.Load(c, &m.{{$m.GoIdent.GoName}})
	if err != nil {
		return err
	}
	return nil
}

func (m *{{$m.GoIdent.GoName}}Model) ModelAutoSave(c context.Context) error {
	if m == nil {
		return nil
	}
	if m.db == nil {
		return nil
	}

	return m.db.AutoSave(c, &m.{{$m.GoIdent.GoName}})
}


func (m *{{$m.GoIdent.GoName}}Model) ModelSave(c context.Context) error {
	if m == nil {
		return nil
	}
	if m.db == nil {
		return nil
	}
	
	return m.db.Save(c, &m.{{$m.GoIdent.GoName}})
}

func (m *{{$m.GoIdent.GoName}}Model) ModelAddWatcher(c context.Context, s model.Watcher) {
	if m == nil {
		return
	}
	m.watchers = append(m.watchers, s)
	s.Watch(c, &m.{{$m.GoIdent.GoName}})
}

func (m *{{$m.GoIdent.GoName}}Model) ModelStopWatcher(c context.Context,) {
	if m == nil {
		return
	}
	for _, w := range m.watchers {
		w.QuitWatch(c, &m.{{$m.GoIdent.GoName}})
	}
}

{{ range $fi, $f := $m.Fields }} 
func (m *{{$m.GoIdent.GoName}}Model) Set{{$f.GoName}}(v {{FieldType $f $.GoImportPath}}) {
	if m == nil{
		return
	}
	m.{{$f.GoName}} = v
	
	for _, w := range m.watchers {
		w.SetNotify(&m.{{$m.GoIdent.GoName}}, {{$f.Desc.Number}}, {{$m.GoIdent.GoName}}_{{$f.GoName}})
	}
}

{{end}}
{{end}}
{{end}}