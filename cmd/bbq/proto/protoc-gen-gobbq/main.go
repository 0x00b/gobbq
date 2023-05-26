package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/0x00b/gobbq/cmd/bbq/proto/com"
	"github.com/0x00b/gobbq/cmd/bbq/proto/protoc-gen-gobbq/gogen"
	options "google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/compiler/protogen"
	protobuf "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

var (
	flags        flag.FlagSet
	plugins      = flags.String("plugins", "grpc", "base framework, grpc/")
	importPrefix = flags.String("import_prefix", "", "prefix to prepend to import paths")
	lang         = flags.String("lang", "go", "language, eg: go")
	tplDir       = flags.String("tpl_dir", "", "tpl dir")
	versionFlag  = flags.Bool("version", false, "print the current version")
)

// Variables set by goreleaser at build time
var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	// flag.Parse()

	if *versionFlag {
		fmt.Printf("Version %v, commit %v, built at %v\n", version, commit, date)
		os.Exit(0)
	}

	importRewriteFunc := func(importPath protogen.GoImportPath) protogen.GoImportPath {
		switch importPath {
		case "context", "fmt", "math":
			return importPath
		}
		if *importPrefix != "" {
			return protogen.GoImportPath(*importPrefix) + importPath
		}
		return importPath
	}

	protogen.Options{
		ParamFunc:         flags.Set,
		ImportRewriteFunc: importRewriteFunc,
	}.Run(func(plugin *protogen.Plugin) error {

		g, err := GetGenerator(plugin)
		if err != nil {
			return err
		}

		return InitAndGenerate(g, plugin)

	})

}

// InitAndGenerate TODO
func InitAndGenerate(g com.Generator, plugin *protogen.Plugin) error {
	req := plugin.Request
	if req == nil {
		return errors.New("pluginpb.CodeGeneratorRequest is nil")
	}

	proto := &com.Proto{
		Plugin: plugin,
		// Files: ,
	}
	for _, f := range req.ProtoFile {
		gf := plugin.FilesByPath[*f.Name]
		if !gf.Generate {
			continue
		}
		file := &com.File{
			Name:                    f.Name,
			Package:                 f.Package,
			Dependency:              f.Dependency,
			PublicDependency:        f.PublicDependency,
			WeakDependency:          f.WeakDependency,
			Options:                 f.Options,
			SourceCodeInfo:          f.SourceCodeInfo,
			Syntax:                  f.Syntax,
			Desc:                    gf.Desc,
			GoDescriptorIdent:       gf.GoDescriptorIdent,
			GoPackageName:           string(gf.GoPackageName),
			GoImportPath:            string(gf.GoImportPath),
			Enums:                   gf.Enums,
			Messages:                gf.Messages,
			Extensions:              gf.Extensions,
			Generate:                gf.Generate,
			GeneratedFilenamePrefix: gf.GeneratedFilenamePrefix,
			// GoImplPackage:           "",
			// Services:                []*com.Service{},
		}
		for _, s := range f.Service {
			gs, err := getService(gf, s)
			if err != nil {
				return err
			}
			// EnumValueDescriptor │ google.protobuf.EnumValueOptions
			// gs.Desc.Options().ProtoReflect().Get(google.protobuf.EnumValueOptions)
			// gs.Desc.Options().(*descriptorpb.ServiceOptions).GetDeprecated()
			svc := &com.Service{
				Name:     s.Name,
				Options:  s.Options,
				Desc:     gs.Desc,
				GoName:   gs.GoName,
				Location: gs.Location,
				Comments: gs.Comments,
				// Methods:  []*com.Method{},
			}
			for _, m := range s.Method {
				gm, err := getMothed(gs, m)
				if err != nil {
					return err
				}
				method := &com.Method{
					Name:            m.Name,
					InputType:       m.InputType,
					OutputType:      m.OutputType,
					Options:         m.Options,
					ClientStreaming: m.ClientStreaming,
					ServerStreaming: m.ServerStreaming,
					Desc:            gm.Desc,
					GoName:          gm.GoName,
					GoInput:         gm.Input,
					GoOutput:        gm.Output,
					Location:        gm.Location,
					Comments:        gm.Comments,
					HasResponse:     hasResponse(gm),
					// MethodBody:      "",
				}
				if protobuf.HasExtension(method.Options, options.E_Http) {
					svc.HasHTTPOption = true
				}
				// _ = g.SetMethodBody(proto, file, svc, method)
				svc.Methods = append(svc.Methods, method)
			}
			file.Services = append(file.Services, svc)
		}
		proto.Files = append(proto.Files, file)
	}
	return g.Generate(*tplDir, proto)
}

func getService(f *protogen.File, s *descriptorpb.ServiceDescriptorProto) (*protogen.Service, error) {
	goName := com.GoCamelCase(*s.Name)
	for _, s := range f.Services {
		if s.GoName == goName {
			return s, nil
		}
	}

	return nil, errors.New("not found service")
}

func getMothed(f *protogen.Service, s *descriptorpb.MethodDescriptorProto) (*protogen.Method, error) {
	goName := com.GoCamelCase(*s.Name)
	for _, s := range f.Methods {
		if s.GoName == goName {
			return s, nil
		}
	}

	return nil, errors.New("not found method")
}

// GetGenerator TODO
func GetGenerator(pp *protogen.Plugin) (com.Generator, error) {

	switch *lang {
	case "go":
		return gogen.NewGoGenerator(".")
	case "ts":
		// todo
	}
	return nil, fmt.Errorf("unkown %s", *lang)
}

func hasResponse(m *protogen.Method) bool {

	ret := true
	if m.Output.GoIdent.GoName == "Empty" && strings.Contains(m.Output.GoIdent.GoImportPath.String(), "golang/protobuf/ptypes/empty") {
		ret = false
	}

	return ret
}
