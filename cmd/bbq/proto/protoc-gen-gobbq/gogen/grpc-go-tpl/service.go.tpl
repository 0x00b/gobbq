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
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/tool/snowflake"
	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/engine/nets"
	"github.com/0x00b/gobbq/proto/bbq"
	"github.com/0x00b/gobbq/xlog"

	// {{$.GoPackageName}} "{{$.GoImportPath}}"
{{end}}
)
 
var _= snowflake.GenUUID()

{{range $sidx, $s := $.Services}}
{{$sName := $s.GoName}}
{{$isSvc := isService $s}}

{{$typeName := concat "" $sName "Entity"}}

{{- if $isSvc}}
	{{$typeName = concat "" $sName "Service"}}

func Register{{$typeName}}(impl {{$typeName}}) {
	entity.Manager.RegisterService(&{{$typeName}}Desc, impl)
}

func New{{$typeName}}Client(client *nets.Client) *{{lowerCamelcase $typeName}} {
	t := &{{lowerCamelcase $typeName}}{client:client}
	return t
}

func New{{$typeName}}(client *nets.Client) *{{lowerCamelcase $typeName}} {
	t := &{{lowerCamelcase $typeName}}{client:client}
	return t
}

type {{lowerCamelcase $typeName}} struct{
	client *nets.Client
}

{{else}}

func Register{{$typeName}}(impl {{$typeName}}) {
	entity.Manager.RegisterEntity(&{{$typeName}}Desc, impl)
}

func New{{$typeName}}Client(client *nets.Client, entity entity.EntityID) *{{lowerCamelcase $typeName}} {
	t := &{{lowerCamelcase $typeName}}{client:client, entity:entity}
	return t
}

func New{{$typeName}}(c *entity.Context, client *nets.Client) *{{lowerCamelcase $typeName}}  {
	return New{{$typeName}}WithID(c, entity.EntityID(snowflake.GenUUID()), client)
}

func New{{$typeName}}WithID(c *entity.Context, id entity.EntityID, client *nets.Client) *{{lowerCamelcase $typeName}}  {

	err := entity.NewEntity(c, &id, {{$typeName}}Desc.TypeName)
	if err != nil {
		xlog.Errorln("new entity err")
		return nil
	}
	t := &{{lowerCamelcase $typeName}}{entity: id, client:client}

	return t
}

type {{lowerCamelcase $typeName}} struct{
	entity entity.EntityID

	client *nets.Client
}
{{end -}}



{{range $midx, $m := $s.Methods}}
{{- if $m.ClientStreaming}}
{{else if $m.ServerStreaming}}
{{else}}
func (t *{{lowerCamelcase $typeName}}){{$m.GoName}}(c *entity.Context, req *{{$m.GoInput.GoIdent.GoName}} {{if $m.HasResponse}}, callback func(c *entity.Context, rsp *{{$m.GoOutput.GoIdent.GoName}}){{end}}) (err error){

	eid := ""
	if c != nil {
		eid = string(c.EntityID())
	}
	pkt, release := codec.NewPacket()
	defer release()
 
	pkt.Header.Version=      1
	pkt.Header.RequestId=    snowflake.GenUUID()
	pkt.Header.Timeout=      1
	pkt.Header.RequestType=  bbq.RequestType_RequestRequest 
	pkt.Header.ServiceType=  {{if $isSvc}}bbq.ServiceType_Service{{else}}bbq.ServiceType_Entity{{end}} 
	pkt.Header.SrcEntity=    eid
	pkt.Header.DstEntity=    {{if $isSvc}}""{{else}}string(t.entity){{end}} 
	pkt.Header.ServiceName=	 "{{$.GoPackageName}}.{{$typeName}}" 
	pkt.Header.Method=       "{{$m.GoName}}" 
	pkt.Header.ContentType=  bbq.ContentType_Proto
	pkt.Header.CompressType= bbq.CompressType_None
	pkt.Header.CheckFlags=   0
	pkt.Header.TransInfo=    map[string][]byte{}
	pkt.Header.ErrCode=      0
	pkt.Header.ErrMsg=       "" 
{{if $m.HasResponse}}
	itfCallback := func(c *entity.Context, rsp any) {
  		callback(c, rsp.(*{{$m.GoOutput.GoIdent.GoName}}))
	}
	_=itfCallback
{{end}}
	// err = entity.HandleCallLocalMethod(pkt, req, itfCallback)
	// if err == nil {
	// 	return nil
	// }

	// if entity.NotMyMethod(err) {

		hdrBytes, err := codec.GetCodec(bbq.ContentType_Proto).Marshal(req)
		if err != nil {
			return err
		}

		pkt.WriteBody(hdrBytes)

		t.client.WritePacket(pkt)

	{{if $m.HasResponse}}
		// register callback
		if (c != nil){
			c.Entity.RegisterCallback(pkt.Header.RequestId, func(c *entity.Context, pkt *codec.Packet) {
				in := new({{$m.GoOutput.GoIdent.GoName}})
				reqbuf := pkt.PacketBody()
				err := codec.GetCodec(pkt.Header.GetContentType()).Unmarshal(reqbuf, in)
				if err != nil {
					return
				}

				callback(c, in)
			})
		}
	{{end}}

	// }

	return err

}
{{end -}}
{{end -}}

// {{goComments $typeName $s.Comments}}
type {{$typeName}} interface {
	entity.IEntity

{{range $midx, $m := $s.Methods}}
// {{goComments $m.GoName $m.Comments}}
{{- if $m.ClientStreaming}}
{{else if $m.ServerStreaming}}
{{else}}
	{{$m.GoName}}(c *entity.Context, req *{{$m.GoInput.GoIdent.GoName}} {{if $m.HasResponse}}, ret func(*{{$m.GoOutput.GoIdent.GoName}}, error){{end}})
{{end -}}
{{end -}}
}


{{range $midx, $m := $s.Methods}}
{{- if $m.ClientStreaming}}
{{else if $m.ServerStreaming}}
{{else}}

func _{{$typeName}}_{{$m.GoName}}_Handler(svc any, ctx *entity.Context, in *{{$m.GoInput.GoIdent.GoName}} {{if $m.HasResponse}}, ret func(rsp *{{$m.GoOutput.GoIdent.GoName}}, err error){{end}}, interceptor entity.ServerInterceptor) {
	if interceptor == nil {
		svc.({{$typeName}}).{{$m.GoName}}(ctx, in ,{{if $m.HasResponse}}ret{{end}})
		return
	}
	
	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/{{$.GoPackageName}}.{{$typeName}}/{{$m.GoName}}",
	}

	handler := func(ctx *entity.Context, rsp any, _ entity.RetFunc) {
		svc.({{$typeName}}).{{$m.GoName}}(ctx, in, {{if $m.HasResponse}}ret{{end}})
	}
 
	interceptor(ctx, in, info, {{if $m.HasResponse}}func(i any, err error) {ret(i.(*{{$m.GoOutput.GoIdent.GoName}}), err)}{{else}}nil{{end}}, handler)
	
}

func _{{$typeName}}_{{$m.GoName}}_Local_Handler(svc any, ctx *entity.Context, in any, callback func(c *entity.Context, rsp any), interceptor entity.ServerInterceptor) {
	{{if $m.HasResponse}}
		ret := func(rsp *{{$m.GoOutput.GoIdent.GoName}}, err error) {
			if err != nil {
				_ = err
			}
			callback(ctx, rsp)
		}
	{{end}}
	
	_{{$typeName}}_{{$m.GoName}}_Handler(svc, ctx, in.(*{{$m.GoInput.GoIdent.GoName}}) {{if $m.HasResponse}}, ret{{end}}, interceptor)
	
}

func _{{$typeName}}_{{$m.GoName}}_Remote_Handler(svc any, ctx *entity.Context, pkt *codec.Packet, interceptor entity.ServerInterceptor) {
 
	hdr := pkt.Header
{{if $m.HasResponse}}
	srcPrw := pkt.Src
	ret:=func(rsp *{{$m.GoOutput.GoIdent.GoName}},err error){ 
		if err != nil {
			xlog.Errorln("err", err)
			return
		}

		npkt, release := codec.NewPacket()
		defer release()

		npkt.Header.Version=      hdr.Version
		npkt.Header.RequestId=    hdr.RequestId
		npkt.Header.Timeout=      hdr.Timeout
		npkt.Header.RequestType=  bbq.RequestType_RequestRespone
		npkt.Header.ServiceType=  hdr.ServiceType
		npkt.Header.SrcEntity=    hdr.DstEntity
		npkt.Header.DstEntity=    hdr.SrcEntity
		npkt.Header.ServiceName=  hdr.ServiceName
		npkt.Header.Method=       hdr.Method
		npkt.Header.ContentType=  hdr.ContentType
		npkt.Header.CompressType= hdr.CompressType
		npkt.Header.CheckFlags=   0
		npkt.Header.TransInfo=    hdr.TransInfo
		npkt.Header.ErrCode= 0
		npkt.Header.ErrMsg=  "" 

		rb, err := codec.DefaultCodec.Marshal(rsp)
		if err != nil {
			xlog.Errorln("Marshal(rsp)", err)
			return
		}

		npkt.WriteBody(rb)

		err = srcPrw.WritePacket(npkt)
		if err != nil {
			xlog.Errorln("WritePacket", err)
			return
		}
	}
{{end}}

	in := new({{$m.GoInput.GoIdent.GoName}})
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)
	if err != nil {
{{- if $m.HasResponse}}
		ret(nil, err)
{{end -}}

	}

	_{{$typeName}}_{{$m.GoName}}_Handler(svc, ctx, in {{if $m.HasResponse}},ret{{end}}, interceptor)

}

{{end -}}
{{end -}}

var {{$typeName}}Desc = entity.EntityDesc{
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