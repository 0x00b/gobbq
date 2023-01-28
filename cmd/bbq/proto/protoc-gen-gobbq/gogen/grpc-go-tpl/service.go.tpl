// NOTE:!!
// DON'T EDIT IT!!
//  请勿添加其他内容，包括函数，结构体，变量等等，否则重新生成时会丢失。

package {{$.GoPackageName}}

import (

{{with $.GoImplImportPaths}}
	{{range $idx, $path := $.GoImplImportPaths}} 
	{{- $path.Alias}} "{{$path.ImportPath -}}"
	{{end}}
{{else}}
	"context" 
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/tool/snowflake"
	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/components/proxy/ex"
	"github.com/0x00b/gobbq/proto"

	// {{$.GoPackageName}} "{{$.GoImportPath}}"
{{end}}
)
 

{{range $sidx, $s := $.Services}}
{{$sName := $s.GoName}}
{{$isSvc := isService $s}}

{{$typeName := concat "" $sName "Entity"}}
{{- if $isSvc}}
	{{$typeName = concat "" $sName "Service"}}

func Register{{$typeName}}(impl {{$typeName}}) {
	entity.Manager.RegisterService(&{{$typeName}}Desc, impl)
}

func New{{$typeName}}() {{$typeName}} {
	t := &{{lowerCamelcase $typeName}}{}
	return t
}

type {{lowerCamelcase $typeName}} struct{
	entity.Service
}

{{else}}

func Register{{$typeName}}(impl {{$typeName}}) {
	entity.Manager.RegisterEntity(&{{$typeName}}Desc, impl)
}

func New{{$typeName}}() {{$typeName}} {
	return New{{$typeName}}WithID(entity.EntityID(snowflake.GenUUID()))
}

func New{{$typeName}}WithID(id entity.EntityID) {{$typeName}} {

	ety := entity.NewEntity(id, {{$typeName}}Desc.TypeName)
	t := &{{lowerCamelcase $typeName}}{entity: ety}

	return t
}

type {{lowerCamelcase $typeName}} struct{
	entity.EntityClient
	entity *proto.Entity
}
{{end -}}



{{range $midx, $m := $s.Methods}}
{{- if $m.ClientStreaming}}
{{else if $m.ServerStreaming}}
{{else}}
func (t *{{lowerCamelcase $typeName}}){{$m.GoName}}(c context.Context, req *{{$m.GoInput.GoIdent.GoName}}) (rsp *{{$m.GoOutput.GoIdent.GoName}},err error){

	pkt := codec.NewPacket()

	hdr := &proto.Header{
		Version:    1,
		RequestId:  "1",
		Timeout:    1,
		Method:     "{{$.GoPackageName}}.{{$typeName}}/{{$m.GoName}}",
		// TransInfo:  map[string][]byte{"xxx": []byte("22222")},
		CallType:   {{- if $isSvc}}proto.CallType_CallService{{else}}proto.CallType_CallEntity{{end -}} ,
		SrcEntity:  nil,
		DstEntity:  {{- if $isSvc}}nil{{else}}t.entity{{end -}} ,
		CheckFlags: codec.FlagDataChecksumIEEE,
	}

	pkt.SetHeader(hdr)

	hdrBytes, err := codec.GetCodec(proto.ContentType_Proto).Marshal(req)
	if err != nil {
 		return nil, err
	}

	pkt.WriteBody(hdrBytes)

	ex.SendProxy(pkt)

	//todo get response

	return nil, nil

}
{{end -}}
{{end -}}

// {{goComments $typeName $s.Comments}}
type {{$typeName}} interface {

{{- if $isSvc}}
	entity.IService
{{else}}
	entity.IEntity
{{end -}}

{{range $midx, $m := $s.Methods}}
// {{goComments $m.GoName $m.Comments}}
{{- if $m.ClientStreaming}}
{{else if $m.ServerStreaming}}
{{else}}
	{{$m.GoName}}(c context.Context, req *{{$m.GoInput.GoIdent.GoName}}) (rsp *{{$m.GoOutput.GoIdent.GoName}},err error)
{{end -}}
{{end -}}
}


{{range $midx, $m := $s.Methods}}
{{- if $m.ClientStreaming}}
{{else if $m.ServerStreaming}}
{{else}}
func _{{$typeName}}_{{$m.GoName}}_Handler(svc interface{}, ctx context.Context, dec func(interface{}) error, interceptor entity.UnaryServerInterceptor) (interface{}, error) {
	in := new({{$m.GoInput.GoIdent.GoName}})
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return svc.({{$typeName}}).{{$m.GoName}}(ctx, in)
	}
	return nil, nil
}
{{end -}}
{{end -}}

var {{$typeName}}Desc = entity.ServiceDesc{
	TypeName:    "{{$.GoPackageName}}.{{$typeName}}",
	HandlerType: (*{{$typeName}})(nil),
	Methods: map[string]entity.MethodDesc{

{{range $midx, $m := $s.Methods}}
{{- if $m.ClientStreaming}}
{{else if $m.ServerStreaming}}
{{else}}
		"{{$m.GoName}}": {
			MethodName: "{{$m.GoName}}",
			Handler:    _{{$typeName}}_{{$m.GoName}}_Handler,
		},
{{end -}}
{{end -}}
	},

	Metadata: "{{$.Name}}",
}


{{end -}}