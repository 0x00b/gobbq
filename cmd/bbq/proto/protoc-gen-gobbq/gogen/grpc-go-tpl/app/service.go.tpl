// NOTE:!!
//  可以在下面的函数中添加具体的实现。
//  请勿添加其他内容，包括函数，结构体，变量等等，否则重新生成时会丢失。

package {{$.GoPackageName}}

import (

{{with $.GoImplImportPaths}}
	{{range $idx, $path := $.GoImplImportPaths}} 
	{{- $path.Alias}} "{{$path.ImportPath -}}"
	{{end}}
{{else}}
	"context" 
	"google.golang.org/grpc"

	{{$.GoPackageName}} "{{$.GoImportPath}}"
{{end}}
)
 
{{range $sidx, $s := $.Services}}
{{$sName := $s.GoName}}
{{$typeName := concat "" $sName "Server"}}
// {{goComments $typeName $s.Comments}}
type {{$typeName}} struct {
	{{$.GoPackageName}}.Unimplemented{{$typeName}}
}

{{$newfunc := concat "" "New" $typeName}}
// {{$newfunc}} 修改此函数，自定义初始化{{$typeName}}
func {{$newfunc}}()*{{$typeName}}{
	{{- $body := goGetMethodBody $.GoRewriter "" $newfunc -}}
	{{- with $body -}}
	{{- $body -}}
	{{- else -}}return &{{$typeName}}{}{{- end -}}
}

// UnaryServerInterceptor 修改此接口，返回只作用于 {{$typeName}}的中间件
func (a *{{$typeName}})UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	{{- $body := goGetMethodBody $.GoRewriter $typeName "UnaryServerInterceptor" -}}
	{{- with $body -}}
	{{- $body -}}
	{{- else -}}return nil{{- end -}}
}

// StreamServerInterceptor 修改此接口，返回只作用于 {{$typeName}}的中间件
func (a *{{$typeName}}) StreamServerInterceptor() grpc.StreamServerInterceptor {
	{{- $body := goGetMethodBody $.GoRewriter $typeName "StreamServerInterceptor" -}}
	{{- with $body -}}
	{{- $body -}}
	{{- else -}}return nil{{- end -}}
}


{{range $midx, $m := $s.Methods}}
// {{goComments $m.GoName $m.Comments}}
{{- if $m.ClientStreaming}}
func (a *{{$typeName}}) {{$m.GoName}}(s {{$.GoPackageName}}.{{$sName}}_{{$m.GoName}}Server) (err error) {
    {{- $body := goGetMethodBody $.GoRewriter $typeName $m.GoName -}}
	{{- with $body -}}
	{{- $body -}}
	{{- else -}}panic("not implemented"){{- end -}}
}
{{else if $m.ServerStreaming}}
func (a *{{$typeName}}) {{$m.GoName}}(
	req *{{$.GoPackageName}}.{{$m.GoInput.GoIdent.GoName}},s {{$.GoPackageName}}.{{$sName}}_{{$m.GoName}}Server) (err error) {
    {{- $body := goGetMethodBody $.GoRewriter $typeName $m.GoName -}}
	{{- with $body -}}
	{{- $body -}}
	{{- else -}}panic("not implemented"){{- end -}}
}
{{else}}
func (a *{{$typeName}}) {{$m.GoName}}(c context.Context, 
    req *{{$.GoPackageName}}.{{$m.GoInput.GoIdent.GoName}}) (rsp *{{$.GoPackageName}}.{{$m.GoOutput.GoIdent.GoName}},err error) {
    {{- $body := goGetMethodBody $.GoRewriter $typeName $m.GoName -}}
	{{- with $body -}}
	{{- $body -}}
	{{- else -}}panic("not implemented"){{- end -}}
}
{{end -}}
{{end -}}
{{end -}}