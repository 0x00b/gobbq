// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.6.1
// source: testpb.proto

package testpb

import (
	_ "github.com/0x00b/gobbq/proto/bbq"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type StartFrameReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *StartFrameReq) Reset() {
	*x = StartFrameReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_testpb_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StartFrameReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StartFrameReq) ProtoMessage() {}

func (x *StartFrameReq) ProtoReflect() protoreflect.Message {
	mi := &file_testpb_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StartFrameReq.ProtoReflect.Descriptor instead.
func (*StartFrameReq) Descriptor() ([]byte, []int) {
	return file_testpb_proto_rawDescGZIP(), []int{0}
}

type StartFrameRsp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FrameSvr uint64 `protobuf:"varint,1,opt,name=FrameSvr,proto3" json:"FrameSvr,omitempty"`
}

func (x *StartFrameRsp) Reset() {
	*x = StartFrameRsp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_testpb_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StartFrameRsp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StartFrameRsp) ProtoMessage() {}

func (x *StartFrameRsp) ProtoReflect() protoreflect.Message {
	mi := &file_testpb_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StartFrameRsp.ProtoReflect.Descriptor instead.
func (*StartFrameRsp) Descriptor() ([]byte, []int) {
	return file_testpb_proto_rawDescGZIP(), []int{1}
}

func (x *StartFrameRsp) GetFrameSvr() uint64 {
	if x != nil {
		return x.FrameSvr
	}
	return 0
}

var File_testpb_proto protoreflect.FileDescriptor

var file_testpb_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x74, 0x65, 0x73, 0x74, 0x70, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06,
	0x74, 0x65, 0x73, 0x74, 0x70, 0x62, 0x1a, 0x09, 0x62, 0x62, 0x71, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x0f, 0x0a, 0x0d, 0x53, 0x74, 0x61, 0x72, 0x74, 0x46, 0x72, 0x61, 0x6d, 0x65, 0x52,
	0x65, 0x71, 0x22, 0x2b, 0x0a, 0x0d, 0x53, 0x74, 0x61, 0x72, 0x74, 0x46, 0x72, 0x61, 0x6d, 0x65,
	0x52, 0x73, 0x70, 0x12, 0x1a, 0x0a, 0x08, 0x46, 0x72, 0x61, 0x6d, 0x65, 0x53, 0x76, 0x72, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x46, 0x72, 0x61, 0x6d, 0x65, 0x53, 0x76, 0x72, 0x32,
	0x4b, 0x0a, 0x05, 0x46, 0x72, 0x61, 0x6d, 0x65, 0x12, 0x3c, 0x0a, 0x0a, 0x53, 0x74, 0x61, 0x72,
	0x74, 0x46, 0x72, 0x61, 0x6d, 0x65, 0x12, 0x15, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x70, 0x62, 0x2e,
	0x53, 0x74, 0x61, 0x72, 0x74, 0x46, 0x72, 0x61, 0x6d, 0x65, 0x52, 0x65, 0x71, 0x1a, 0x15, 0x2e,
	0x74, 0x65, 0x73, 0x74, 0x70, 0x62, 0x2e, 0x53, 0x74, 0x61, 0x72, 0x74, 0x46, 0x72, 0x61, 0x6d,
	0x65, 0x52, 0x73, 0x70, 0x22, 0x00, 0x1a, 0x04, 0xa0, 0xbb, 0x18, 0x01, 0x42, 0x2e, 0x5a, 0x2c,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x30, 0x78, 0x30, 0x30, 0x62,
	0x2f, 0x67, 0x6f, 0x62, 0x62, 0x71, 0x2f, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2f, 0x65,
	0x78, 0x61, 0x6d, 0x70, 0x62, 0x3b, 0x74, 0x65, 0x73, 0x74, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_testpb_proto_rawDescOnce sync.Once
	file_testpb_proto_rawDescData = file_testpb_proto_rawDesc
)

func file_testpb_proto_rawDescGZIP() []byte {
	file_testpb_proto_rawDescOnce.Do(func() {
		file_testpb_proto_rawDescData = protoimpl.X.CompressGZIP(file_testpb_proto_rawDescData)
	})
	return file_testpb_proto_rawDescData
}

var file_testpb_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_testpb_proto_goTypes = []interface{}{
	(*StartFrameReq)(nil), // 0: testpb.StartFrameReq
	(*StartFrameRsp)(nil), // 1: testpb.StartFrameRsp
}
var file_testpb_proto_depIdxs = []int32{
	0, // 0: testpb.Frame.StartFrame:input_type -> testpb.StartFrameReq
	1, // 1: testpb.Frame.StartFrame:output_type -> testpb.StartFrameRsp
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_testpb_proto_init() }
func file_testpb_proto_init() {
	if File_testpb_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_testpb_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StartFrameReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_testpb_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StartFrameRsp); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_testpb_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_testpb_proto_goTypes,
		DependencyIndexes: file_testpb_proto_depIdxs,
		MessageInfos:      file_testpb_proto_msgTypes,
	}.Build()
	File_testpb_proto = out.File
	file_testpb_proto_rawDesc = nil
	file_testpb_proto_goTypes = nil
	file_testpb_proto_depIdxs = nil
}