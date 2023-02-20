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

func New{{$typeName}}Client(client *codec.PacketReadWriter) *{{lowerCamelcase $typeName}} {
	t := &{{lowerCamelcase $typeName}}{client:client}
	return t
}

func New{{$typeName}}(client *codec.PacketReadWriter) *{{lowerCamelcase $typeName}} {
	t := &{{lowerCamelcase $typeName}}{client:client}
	return t
}

type {{lowerCamelcase $typeName}} struct{
	client *codec.PacketReadWriter
}

{{else}}

func Register{{$typeName}}(impl {{$typeName}}) {
	entity.Manager.RegisterEntity(&{{$typeName}}Desc, impl)
}

func New{{$typeName}}Client(client *codec.PacketReadWriter, entity entity.EntityID) *{{lowerCamelcase $typeName}} {
	t := &{{lowerCamelcase $typeName}}{client:client, entity:entity}
	return t
}

func New{{$typeName}}(c entity.Context, client *codec.PacketReadWriter) *{{lowerCamelcase $typeName}}  {
	return New{{$typeName}}WithID(c,*entity.NewEntityID.NewEntityID("{{$.GoPackageName}}.{{$typeName}}"), client)
}

func New{{$typeName}}WithID(c entity.Context, id entity.EntityID, client *codec.PacketReadWriter) *{{lowerCamelcase $typeName}}  {

	err := entity.NewEntity(c, &id)
	if err != nil {
		xlog.Errorln("new entity err")
		return nil
	}
	t := &{{lowerCamelcase $typeName}}{entity: id, client:client}

	return t
}

type {{lowerCamelcase $typeName}} struct{
	entity entity.EntityID

	client *codec.PacketReadWriter
}
{{end -}}



{{range $midx, $m := $s.Methods}}
{{- if $m.ClientStreaming}}
{{else if $m.ServerStreaming}}
{{else}}
func (t *{{lowerCamelcase $typeName}}){{$m.GoName}}(c entity.Context, req *{{$m.GoInput.GoIdent.GoName}}) {{if $m.HasResponse}}(*{{$m.GoOutput.GoIdent.GoName}}, error){{end}}{

	pkt, release := codec.NewPacket()
	defer release()
 
	pkt.Header.Version=      1
	pkt.Header.RequestId=    snowflake.GenUUID()
	pkt.Header.Timeout=      1
	pkt.Header.RequestType=  bbq.RequestType_RequestRequest 
	pkt.Header.ServiceType=  {{if $isSvc}}bbq.ServiceType_Service{{else}}bbq.ServiceType_Entity{{end}} 
	pkt.Header.SrcEntity =  entity.ToPBEntityID(c.EntityID())
	pkt.Header.DstEntity = {{if $isSvc}}&bbq.EntityID{Type: "{{$.GoPackageName}}.{{$typeName}}"}{{else}}entity.ToPBEntityID(t.entity){{end}}
	pkt.Header.Method=       "{{$m.GoName}}" 
	pkt.Header.ContentType=  bbq.ContentType_Proto
	pkt.Header.CompressType= bbq.CompressType_None
	pkt.Header.CheckFlags=   0
	pkt.Header.TransInfo=    map[string][]byte{}
	pkt.Header.ErrCode=      0
	pkt.Header.ErrMsg=       "" 

	// err = entity.HandleCallLocalMethod(pkt, req, itfCallback)
	// if err == nil {
	// 	return nil
	// }

	hdrBytes, err := codec.GetCodec(bbq.ContentType_Proto).Marshal(req)
	if err != nil {
		xlog.Errorln(err)
		return {{if $m.HasResponse}}nil, err{{end}}
	}

	pkt.WriteBody(hdrBytes)

	t.client.WritePacket(pkt)

	{{if $m.HasResponse}}
		// register callback
		chanRsp := make(chan any)
		if c != nil {
			c.RegisterCallback(pkt.Header.RequestId, func(pkt *codec.Packet) {
				rsp := new({{$m.GoOutput.GoIdent.GoName}})
				reqbuf := pkt.PacketBody()
				err := codec.GetCodec(pkt.Header.GetContentType()).Unmarshal(reqbuf, rsp)
				if err != nil {
					chanRsp <- err
					return
				}
				chanRsp <- rsp
				close(chanRsp)
			})
		}
		rsp := <-chanRsp
		if rsp, ok := rsp.(*{{$m.GoOutput.GoIdent.GoName}}); ok {
			return rsp, nil
		}
		return nil, rsp.(error)
	{{end}}

}
{{end -}}
{{end -}}

// {{goComments $typeName $s.Comments}}
type {{$typeName}} interface {
	 {{if $isSvc}}entity.IService{{else}}entity.IEntity{{end}} 

{{range $midx, $m := $s.Methods}}
// {{goComments $m.GoName $m.Comments}}
{{- if $m.ClientStreaming}}
{{else if $m.ServerStreaming}}
{{else}}
	{{$m.GoName}}(c entity.Context, req *{{$m.GoInput.GoIdent.GoName}}){{if $m.HasResponse}}(*{{$m.GoOutput.GoIdent.GoName}}, error){{end}}
{{end -}}
{{end -}}
}


{{range $midx, $m := $s.Methods}}
{{- if $m.ClientStreaming}}
{{else if $m.ServerStreaming}}
{{else}}

func _{{$typeName}}_{{$m.GoName}}_Handler(svc any, ctx entity.Context, in *{{$m.GoInput.GoIdent.GoName}}, interceptor entity.ServerInterceptor){{if $m.HasResponse}}(*{{$m.GoOutput.GoIdent.GoName}},error){{end}} {
	if interceptor == nil {

	{{if $m.HasResponse}}
		return svc.({{$typeName}}).{{$m.GoName}}(ctx, in)
	{{else}}
		svc.({{$typeName}}).{{$m.GoName}}(ctx, in)
		return
	{{end}}

	}
	
	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/{{$.GoPackageName}}.{{$typeName}}/{{$m.GoName}}",
	}

	handler := func(ctx entity.Context, rsp any)  (any, error) {
	{{if $m.HasResponse}}
		return svc.({{$typeName}}).{{$m.GoName}}(ctx, in)
	{{else}}
		svc.({{$typeName}}).{{$m.GoName}}(ctx, in)
		return nil,nil
	{{end}}
	}
 
{{if $m.HasResponse}}
	rsp, err := interceptor(ctx, in, info, handler)
	return rsp.(*{{$m.GoOutput.GoIdent.GoName}}), err
{{else}}
 	interceptor(ctx, in, info, handler)
{{end}}

}

//func _{{$typeName}}_{{$m.GoName}}_Local_Handler(svc any, ctx entity.Context, in any, interceptor entity.ServerInterceptor)(any, error) {
//	{{if $m.HasResponse}}
//		ret := func(rsp *{{$m.GoOutput.GoIdent.GoName}}, err error) {
//			if err != nil {
//				_ = err
//			}
//			callback(ctx, rsp)
//		}
//	{{end}}
//	
//	_{{$typeName}}_{{$m.GoName}}_Handler(svc, ctx, in.(*{{$m.GoInput.GoIdent.GoName}}) {{if $m.HasResponse}}, ret{{end}}, interceptor)
//	
//}

func _{{$typeName}}_{{$m.GoName}}_Remote_Handler(svc any, ctx entity.Context, pkt *codec.Packet, interceptor entity.ServerInterceptor) {
 
	hdr := pkt.Header
	
	in := new({{$m.GoInput.GoIdent.GoName}})
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)
	if err != nil {
		// {{if $m.HasResponse}}nil,{{end}}err
		return
	}


{{if $m.HasResponse}}
	rsp, err := _{{$typeName}}_{{$m.GoName}}_Handler(svc, ctx, in, interceptor)

	npkt, release := codec.NewPacket()
	defer release()

	npkt.Header.Version=      hdr.Version
	npkt.Header.RequestId=    hdr.RequestId
	npkt.Header.Timeout=      hdr.Timeout
	npkt.Header.RequestType=  bbq.RequestType_RequestRespone
	npkt.Header.ServiceType=  hdr.ServiceType
	npkt.Header.SrcEntity=    hdr.DstEntity
	npkt.Header.DstEntity=    hdr.SrcEntity
	npkt.Header.Method=       hdr.Method
	npkt.Header.ContentType=  hdr.ContentType
	npkt.Header.CompressType= hdr.CompressType
	npkt.Header.CheckFlags=   0
	npkt.Header.TransInfo=    hdr.TransInfo

	if err != nil{
		npkt.Header.ErrCode= 1
		npkt.Header.ErrMsg=  err.Error() 
		
		npkt.WriteBody(nil)
	}else{
		rb, err := codec.DefaultCodec.Marshal(rsp)
		if err != nil {
			xlog.Errorln("Marshal(rsp)", err)
			return
		}

		npkt.WriteBody(rb)
	}
	err = pkt.Src.WritePacket(npkt)
	if err != nil {
		xlog.Errorln("WritePacket", err)
		return
	}
{{else}}
	_{{$typeName}}_{{$m.GoName}}_Handler(svc, ctx, in, interceptor)
{{end}}

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
			//LocalHandler:	_{{$typeName}}_{{$m.GoName}}_Local_Handler,
		},
{{end -}}
{{end -}}
	},

	Metadata: "{{$.Name}}",
}


{{end -}}