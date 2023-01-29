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
	"github.com/0x00b/gobbq/proto/bbq"
	"fmt"

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

func New{{$typeName}}() *{{lowerCamelcase $typeName}} {
	t := &{{lowerCamelcase $typeName}}{}
	return t
}

type {{lowerCamelcase $typeName}} struct{
}

{{else}}

func Register{{$typeName}}(impl {{$typeName}}) {
	entity.Manager.RegisterEntity(&{{$typeName}}Desc, impl)
}

func New{{$typeName}}() *{{lowerCamelcase $typeName}}  {
	return New{{$typeName}}WithID(entity.EntityID(snowflake.GenUUID()))
}

func New{{$typeName}}WithID(id entity.EntityID) *{{lowerCamelcase $typeName}}  {

	ety := entity.NewEntity(id, {{$typeName}}Desc.TypeName)
	t := &{{lowerCamelcase $typeName}}{entity: ety}

	return t
}

type {{lowerCamelcase $typeName}} struct{
	entity *bbq.EntityID
}
{{end -}}



{{range $midx, $m := $s.Methods}}
{{- if $m.ClientStreaming}}
{{else if $m.ServerStreaming}}
{{else}}
func (t *{{lowerCamelcase $typeName}}){{$m.GoName}}(c context.Context, req *{{$m.GoInput.GoIdent.GoName}}, callback func(c context.Context, rsp *{{$m.GoOutput.GoIdent.GoName}})) (err error){

	hdr := &bbq.Header{   
		Version:      1,
		RequestId:    "1",
		Timeout:      1,
		RequestType:  bbq.RequestType_RequestRequest,
		ServiceType:  {{- if $isSvc}}bbq.ServiceType_Service{{else}}bbq.ServiceType_Entity{{end -}},
		SrcEntity:    nil,
		DstEntity:   {{- if $isSvc}}nil{{else}}t.entity{{end -}} ,
		Method:      "{{$.GoPackageName}}.{{$typeName}}/{{$m.GoName}}",
		ContentType:  bbq.ContentType_Proto,
		CompressType: bbq.CompressType_None,
		CheckFlags:   0,
		TransInfo:    map[string][]byte{},
		ErrCode:      0,
		ErrMsg:       "",
	}

	itfCallback := func(c context.Context, rsp interface{}) {
		callback(c, rsp.(*{{$m.GoOutput.GoIdent.GoName}}))
	}

	err = entity.HandleCallLocalMethod(c, hdr, req, itfCallback)
	if err == nil {
		return nil
	}

	if entity.NotMyMethod(err) {

		pkt := codec.NewPacket()

		pkt.SetHeader(hdr)

		hdrBytes, err := codec.GetCodec(bbq.ContentType_Proto).Marshal(req)
		if err != nil {
			return err
		}

		pkt.WriteBody(hdrBytes)

		ex.SendProxy(pkt)
		//todo get response

	}

	return err

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
	{{$m.GoName}}(c context.Context, req *{{$m.GoInput.GoIdent.GoName}}, ret func(*{{$m.GoOutput.GoIdent.GoName}}, error))
{{end -}}
{{end -}}
}


{{range $midx, $m := $s.Methods}}
{{- if $m.ClientStreaming}}
{{else if $m.ServerStreaming}}
{{else}}

func _{{$typeName}}_{{$m.GoName}}_Handler(svc interface{}, ctx context.Context, in *{{$m.GoInput.GoIdent.GoName}}, ret func(rsp *{{$m.GoOutput.GoIdent.GoName}}, err error), interceptor entity.ServerInterceptor) {
	if interceptor == nil {
		svc.({{$typeName}}).{{$m.GoName}}(ctx, in, ret)
		return
	}
	
	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/{{$.GoPackageName}}.{{$typeName}}/{{$m.GoName}}",
	}

	handler := func(ctx context.Context, rsp interface{}, _ entity.RetFunc) {
		svc.({{$typeName}}).{{$m.GoName}}(ctx, in, ret)
	}
 
	interceptor(ctx, in, info, func(i interface{}, err error) {ret(i.(*{{$m.GoOutput.GoIdent.GoName}}), err)}, handler)
	return
}

func _{{$typeName}}_{{$m.GoName}}_Local_Handler(svc interface{}, ctx context.Context, in interface{}, callback func(c context.Context, rsp interface{}), interceptor entity.ServerInterceptor) {
	ret := func(rsp *{{$m.GoOutput.GoIdent.GoName}}, err error) {
		if err != nil {
			_ = err
		}
		callback(ctx, rsp)
	}
	_{{$typeName}}_{{$m.GoName}}_Handler(svc, ctx, in.(*{{$m.GoInput.GoIdent.GoName}}), ret, interceptor)
	return
}

func _{{$typeName}}_{{$m.GoName}}_Remote_Handler(svc interface{}, ctx context.Context, pkt *codec.Packet, interceptor entity.ServerInterceptor) {
 
	hdr := pkt.GetHeader()
	dec := func(v interface{}) error {
		reqbuf := pkt.PacketBody()
		err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, v)
		return err
	}
	in := new({{$m.GoInput.GoIdent.GoName}})

	ret:=func(rsp *{{$m.GoOutput.GoIdent.GoName}},err error){ 
		npkt := codec.NewPacket()

		rhdr := &bbq.Header{
			Version:      hdr.Version,
			RequestId:    hdr.RequestId,
			Timeout:      hdr.Timeout,
			RequestType:  hdr.RequestType,
			ServiceType:  hdr.ServiceType,
			SrcEntity:    hdr.DstEntity,
			DstEntity:    hdr.SrcEntity,
			Method:       hdr.Method,
			ContentType:  hdr.ContentType,
			CompressType: hdr.CompressType,
			CheckFlags:   0,
			TransInfo:    hdr.TransInfo,
			ErrCode: 0,
			ErrMsg:  "",
		}
		npkt.SetHeader(rhdr)

		rbyte, err := codec.DefaultCodec.Marshal(rhdr)
		if err != nil {
			fmt.Println("WritePacket", err)
			return
		}
		npkt.WriteBody(rbyte)

		rb, err := codec.DefaultCodec.Marshal(rsp)
		if err != nil {
			fmt.Println("Marshal(rsp)", err)
			return
		}

		npkt.WriteBody(rb)

		err = pkt.Src.WritePacket(npkt)
		if err != nil {
			fmt.Println("WritePacket", err)
			return
		}
	}

	if err := dec(in); err != nil {
		ret(nil, err)
		return
	}

	_{{$typeName}}_{{$m.GoName}}_Handler(svc, ctx, in, ret, interceptor)
	return
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
			MethodName: 	"{{$m.GoName}}",
			Handler:    	_{{$typeName}}_{{$m.GoName}}_Remote_Handler,
			LocalHandler:	_{{$typeName}}_{{$m.GoName}}_Local_Handler,
		},
{{end -}}
{{end -}}
	},

	Metadata: "{{$.Name}}",
}


{{end -}}