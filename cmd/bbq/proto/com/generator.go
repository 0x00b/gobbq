package com

import (
	"github.com/0x00b/gobbq/cmd/bbq/proto/com/gorewriter/rewrite"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

// Generator generate code
type Generator interface {

	// Generate code
	Generate(tplPath string, proto *Proto) error
}

// Proto is An encoded CodeGeneratorRequest is written to the plugin's stdin.
// https://pkg.go.dev/google.golang.org/protobuf/types/pluginpb#CodeGeneratorRequest
type Proto struct {
	// https://pkg.go.dev/google.golang.org/protobuf@v1.27.1/compiler/protogen#Options.Run
	Plugin *protogen.Plugin

	// All proto file descriptor
	Files []*File

	// .all.tpl 才会有
	GoRewriter *rewrite.Rewriter
}

// File DescriptorProto Describes a complete .proto file.
// https://pkg.go.dev/google.golang.org/protobuf/types/descriptorpb#FileDescriptorProto
type File struct {

	// pluginpb.CodeGeneratorRequest

	Name    *string // file name, relative to root of source tree
	Package *string // e.g. "foo", "foo.bar", etc.
	// Names of files imported by this file.
	Dependency []string
	// Indexes of the public imported files in the dependency list above.
	PublicDependency []int32
	// Indexes of the weak imported files in the dependency list.
	// For Google-internal migration only. Do not use.
	WeakDependency []int32
	Options        *descriptorpb.FileOptions
	// This field contains optional information about the original source code.
	// You may safely remove this entire field without harming runtime
	// functionality of the descriptors -- the information is needed only by
	// development tools.
	SourceCodeInfo *descriptorpb.SourceCodeInfo
	// The syntax of the proto file.
	// The supported values are "proto2" and "proto3".
	Syntax *string

	Desc protoreflect.FileDescriptor

	// Go
	GoImplPackage     string           // 实现本PB协议的package
	GoDescriptorIdent protogen.GoIdent // name of Go variable for the file descriptor
	GoPackageName     string           // name of this file's Go package
	GoImportPath      string           // import path of this file's Go package
	GoImplImportPaths []rewrite.Import // import path of this impelement file's Go package

	GoRewriter *rewrite.Rewriter

	//

	Enums      []*protogen.Enum      // top-level enum declarations
	Messages   []*protogen.Message   // top-level message declarations
	Extensions []*protogen.Extension // top-level extension declarations
	Services   []*Service            // top-level service declarations

	Generate bool // true if we should generate code for this file

	// GeneratedFilenamePrefix is used to construct filenames for generated
	// files associated with this source file.
	//
	// For example, the source file "dir/foo.proto" might have a filename prefix
	// of "dir/foo". Appending ".pb.go" produces an output file of "dir/foo.pb.go".
	GeneratedFilenamePrefix string
}

// Service DescriptorProto Describes a service.
// google.golang.org/protobuf/types/descriptorpb.ServiceDescriptorProto
type Service struct {
	Name *string

	Options *descriptorpb.ServiceOptions
	Desc    protoreflect.ServiceDescriptor

	HasHTTPOption bool

	// Go
	GoName string

	Methods []*Method // service method declarations

	Location protogen.Location   // location of this service
	Comments protogen.CommentSet // comments associated with this service
}

// Method DescriptorProto Describes a method of a service.
// google.golang.org/protobuf/types/descriptorpb.MethodDescriptorProto
type Method struct {
	Name *string
	// Input and output type names.  These are resolved in the same way as
	// FieldDescriptorProto.type_name, but must refer to a message type.
	InputType  *string
	OutputType *string
	Options    *descriptorpb.MethodOptions
	// Identifies if client streams multiple client messages
	ClientStreaming *bool
	// Identifies if server streams multiple server messages
	ServerStreaming *bool

	Desc protoreflect.MethodDescriptor

	// Go
	GoName   string
	GoInput  *protogen.Message
	GoOutput *protogen.Message

	Location protogen.Location   // location of this method
	Comments protogen.CommentSet // comments associated with this method
}
