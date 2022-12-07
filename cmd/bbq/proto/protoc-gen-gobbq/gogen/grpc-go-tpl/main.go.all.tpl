package main

import (

{{$imports := goGetImportPaths $.GoRewriter "main.go"}}
{{with $imports}}
	{{range $idx, $path := $imports}}
	{{- $path.Alias}} "{{$path.ImportPath -}}"
	{{end}}
{{else}}
	"context"
	_ "net/http/pprof"

	tgrpc "github.com/tencent/grpc-go-contrib/grpc"
	"google.golang.org/grpc"
{{end -}}

{{range $index, $f := .Files}}
 	{{$alias := concat "_" $f.GoPackageName "app"}}
	{{$exist := goExistImportPaths $imports $alias $f.GoImplPackage}}
	{{if not $exist}}
	{{$alias}} "{{$f.GoImplPackage}}"
	{{end}}
	{{$exist := goExistImportPaths $imports $f.GoPackageName $f.GoImportPath}}
	{{if not $exist}}
	{{$f.GoPackageName}} "{{$f.GoImportPath}}"
	{{end}}
{{end -}}
)

// Init 修改此函数，业务初始化
func Init() {
	{{- $body := goGetMethodBody $.GoRewriter "" "Init" -}}
	{{- with $body -}}
	{{- $body -}}
	{{- end -}}
}


// GlobalUnaryServerInterceptor 修改此函数，返回自定义unary全局中间件
func GlobalUnaryServerInterceptor(ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error)  {
	{{- $body := goGetMethodBody $.GoRewriter "" "GlobalUnaryServerInterceptor" -}}
	{{- with $body -}}
	{{- $body -}}
	{{- else -}} return handler(ctx, req) {{- end -}}
}
 

// NOTE:!!
//  请勿修改下面的代码，或者添加其他内容，否则重新生成时会丢失。

func main() {

	Init()

	s := tgrpc.NewServer(
		tgrpc.WithServerInterceptor(GlobalUnaryServerInterceptor),
		tgrpc.WithServerInterceptor(unaryServerInterceptor),
	)

	var err error
    {{range $index, $f := .Files}}  
        {{range $sidx, $s := $f.Services}}
			{{$comment := goComments "" $s.Comments}}
			{{with $comment}}// {{$s.GoName}} {{$comment}}{{end}}
			{{$f.GoPackageName}}Svc := {{$f.GoPackageName}}_app.New{{$s.GoName}}Server()
			{{$f.GoPackageName}}.Register{{$s.GoName}}Server(s, {{$f.GoPackageName}}Svc) 
			{{with $s.HasHTTPOption}}// {{$s.GoName}} http service
			{{$f.GoPackageName}}Ctx, {{$f.GoPackageName}}Mux, {{$f.GoPackageName}}Client := s.GetGatewayInfo({{$f.GoPackageName}}Svc, &{{$f.GoPackageName}}.{{$s.GoName}}_ServiceDesc) 
			err={{$f.GoPackageName}}.Register{{$s.GoName}}HandlerClient({{$f.GoPackageName}}Ctx, {{$f.GoPackageName}}Mux,exam.NewEchoClient({{$f.GoPackageName}}Client))
			if err != nil {
				panic(err)
			}
			{{end}}
        {{end -}}
    {{end}}

	_ = s.Serve()
}

// UnaryServerInterceptor define
type UnaryServerInterceptor interface {
	UnaryServerInterceptor() grpc.UnaryServerInterceptor
}

func unaryServerInterceptor(ctx context.Context,
	req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	if si, ok := info.Server.(UnaryServerInterceptor); ok {
		usi := si.UnaryServerInterceptor()
		if usi != nil {
			return usi(ctx, req, info, handler)
		}
	}
	return handler(ctx, req)
}
 