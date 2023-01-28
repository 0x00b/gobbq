// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.6.1
// source: proto.proto

package proto

import (
	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
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

type ContentType int32

const (
	ContentType_Proto ContentType = 0
)

// Enum value maps for ContentType.
var (
	ContentType_name = map[int32]string{
		0: "Proto",
	}
	ContentType_value = map[string]int32{
		"Proto": 0,
	}
)

func (x ContentType) Enum() *ContentType {
	p := new(ContentType)
	*p = x
	return p
}

func (x ContentType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ContentType) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_proto_enumTypes[0].Descriptor()
}

func (ContentType) Type() protoreflect.EnumType {
	return &file_proto_proto_enumTypes[0]
}

func (x ContentType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ContentType.Descriptor instead.
func (ContentType) EnumDescriptor() ([]byte, []int) {
	return file_proto_proto_rawDescGZIP(), []int{0}
}

type CompressType int32

const (
	CompressType_Gzip CompressType = 0
)

// Enum value maps for CompressType.
var (
	CompressType_name = map[int32]string{
		0: "Gzip",
	}
	CompressType_value = map[string]int32{
		"Gzip": 0,
	}
)

func (x CompressType) Enum() *CompressType {
	p := new(CompressType)
	*p = x
	return p
}

func (x CompressType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (CompressType) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_proto_enumTypes[1].Descriptor()
}

func (CompressType) Type() protoreflect.EnumType {
	return &file_proto_proto_enumTypes[1]
}

func (x CompressType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use CompressType.Descriptor instead.
func (CompressType) EnumDescriptor() ([]byte, []int) {
	return file_proto_proto_rawDescGZIP(), []int{1}
}

type CallType int32

const (
	CallType_CallEntity  CallType = 0
	CallType_CallService CallType = 1
)

// Enum value maps for CallType.
var (
	CallType_name = map[int32]string{
		0: "CallEntity",
		1: "CallService",
	}
	CallType_value = map[string]int32{
		"CallEntity":  0,
		"CallService": 1,
	}
)

func (x CallType) Enum() *CallType {
	p := new(CallType)
	*p = x
	return p
}

func (x CallType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (CallType) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_proto_enumTypes[2].Descriptor()
}

func (CallType) Type() protoreflect.EnumType {
	return &file_proto_proto_enumTypes[2]
}

func (x CallType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use CallType.Descriptor instead.
func (CallType) EnumDescriptor() ([]byte, []int) {
	return file_proto_proto_rawDescGZIP(), []int{2}
}

type Entity struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID       string `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	TypeName string `protobuf:"bytes,2,opt,name=TypeName,proto3" json:"TypeName,omitempty"`
}

func (x *Entity) Reset() {
	*x = Entity{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Entity) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Entity) ProtoMessage() {}

func (x *Entity) ProtoReflect() protoreflect.Message {
	mi := &file_proto_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Entity.ProtoReflect.Descriptor instead.
func (*Entity) Descriptor() ([]byte, []int) {
	return file_proto_proto_rawDescGZIP(), []int{0}
}

func (x *Entity) GetID() string {
	if x != nil {
		return x.ID
	}
	return ""
}

func (x *Entity) GetTypeName() string {
	if x != nil {
		return x.TypeName
	}
	return ""
}

// 请求协议头
type Header struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 协议版本
	// 具体值与Version对应
	Version uint32 `protobuf:"varint,1,opt,name=Version,proto3" json:"Version,omitempty"`
	// 具体值与CallType对应
	// 请求唯一id
	RequestId string `protobuf:"bytes,2,opt,name=RequestId,proto3" json:"RequestId,omitempty"`
	// 请求的超时时间，单位ms
	Timeout uint32 `protobuf:"varint,3,opt,name=Timeout,proto3" json:"Timeout,omitempty"`
	// call sverice or call entity
	CallType CallType `protobuf:"varint,4,opt,name=CallType,proto3,enum=proto.CallType" json:"CallType,omitempty"`
	// 调用的原EntityID
	SrcEntity *Entity `protobuf:"bytes,5,opt,name=SrcEntity,proto3" json:"SrcEntity,omitempty"`
	// 调用的目的EntityID
	DstEntity *Entity `protobuf:"bytes,6,opt,name=DstEntity,proto3" json:"DstEntity,omitempty"`
	// 调用服务的接口名
	// 规范格式: /package.Service名称/接口名
	Method string `protobuf:"bytes,7,opt,name=Method,proto3" json:"Method,omitempty"`
	// 请求数据的序列化类型
	// 比如: proto/jce/json, 默认proto
	// 具体值与ContentEncodeType对应
	ContentType ContentType `protobuf:"varint,8,opt,name=ContentType,proto3,enum=proto.ContentType" json:"ContentType,omitempty"`
	// 请求数据使用的压缩方式
	// 比如: gzip/snappy/..., 默认不使用
	// 具体值与CompressType对应
	CompressType CompressType `protobuf:"varint,9,opt,name=CompressType,proto3,enum=proto.CompressType" json:"CompressType,omitempty"`
	// 是否检查包是否正确
	CheckFlags uint32 `protobuf:"varint,10,opt,name=CheckFlags,proto3" json:"CheckFlags,omitempty"`
	// 附加信息
	TransInfo map[string][]byte `protobuf:"bytes,11,rep,name=TransInfo,proto3" json:"TransInfo,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *Header) Reset() {
	*x = Header{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Header) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Header) ProtoMessage() {}

func (x *Header) ProtoReflect() protoreflect.Message {
	mi := &file_proto_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Header.ProtoReflect.Descriptor instead.
func (*Header) Descriptor() ([]byte, []int) {
	return file_proto_proto_rawDescGZIP(), []int{1}
}

func (x *Header) GetVersion() uint32 {
	if x != nil {
		return x.Version
	}
	return 0
}

func (x *Header) GetRequestId() string {
	if x != nil {
		return x.RequestId
	}
	return ""
}

func (x *Header) GetTimeout() uint32 {
	if x != nil {
		return x.Timeout
	}
	return 0
}

func (x *Header) GetCallType() CallType {
	if x != nil {
		return x.CallType
	}
	return CallType_CallEntity
}

func (x *Header) GetSrcEntity() *Entity {
	if x != nil {
		return x.SrcEntity
	}
	return nil
}

func (x *Header) GetDstEntity() *Entity {
	if x != nil {
		return x.DstEntity
	}
	return nil
}

func (x *Header) GetMethod() string {
	if x != nil {
		return x.Method
	}
	return ""
}

func (x *Header) GetContentType() ContentType {
	if x != nil {
		return x.ContentType
	}
	return ContentType_Proto
}

func (x *Header) GetCompressType() CompressType {
	if x != nil {
		return x.CompressType
	}
	return CompressType_Gzip
}

func (x *Header) GetCheckFlags() uint32 {
	if x != nil {
		return x.CheckFlags
	}
	return 0
}

func (x *Header) GetTransInfo() map[string][]byte {
	if x != nil {
		return x.TransInfo
	}
	return nil
}

var file_proto_proto_extTypes = []protoimpl.ExtensionInfo{
	{
		ExtendedType:  (*descriptor.ServiceOptions)(nil),
		ExtensionType: (*CallType)(nil),
		Field:         50001,
		Name:          "proto.call_type",
		Tag:           "varint,50001,opt,name=call_type,enum=proto.CallType",
		Filename:      "proto.proto",
	},
}

// Extension fields to descriptor.ServiceOptions.
var (
	// optional proto.CallType call_type = 50001;
	E_CallType = &file_proto_proto_extTypes[0]
)

var File_proto_proto protoreflect.FileDescriptor

var file_proto_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x20, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x34, 0x0a, 0x06, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79,
	0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x49, 0x44,
	0x12, 0x1a, 0x0a, 0x08, 0x54, 0x79, 0x70, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x54, 0x79, 0x70, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x82, 0x04, 0x0a,
	0x06, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x12, 0x18, 0x0a, 0x07, 0x56, 0x65, 0x72, 0x73, 0x69,
	0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x07, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x12, 0x1c, 0x0a, 0x09, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x49, 0x64, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x49, 0x64, 0x12,
	0x18, 0x0a, 0x07, 0x54, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x07, 0x54, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x12, 0x2b, 0x0a, 0x08, 0x43, 0x61, 0x6c,
	0x6c, 0x54, 0x79, 0x70, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0f, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x61, 0x6c, 0x6c, 0x54, 0x79, 0x70, 0x65, 0x52, 0x08, 0x43, 0x61,
	0x6c, 0x6c, 0x54, 0x79, 0x70, 0x65, 0x12, 0x2b, 0x0a, 0x09, 0x53, 0x72, 0x63, 0x45, 0x6e, 0x74,
	0x69, 0x74, 0x79, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x52, 0x09, 0x53, 0x72, 0x63, 0x45, 0x6e, 0x74,
	0x69, 0x74, 0x79, 0x12, 0x2b, 0x0a, 0x09, 0x44, 0x73, 0x74, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x45,
	0x6e, 0x74, 0x69, 0x74, 0x79, 0x52, 0x09, 0x44, 0x73, 0x74, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79,
	0x12, 0x16, 0x0a, 0x06, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x12, 0x34, 0x0a, 0x0b, 0x43, 0x6f, 0x6e, 0x74,
	0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x12, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70,
	0x65, 0x52, 0x0b, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x37,
	0x0a, 0x0c, 0x43, 0x6f, 0x6d, 0x70, 0x72, 0x65, 0x73, 0x73, 0x54, 0x79, 0x70, 0x65, 0x18, 0x09,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x13, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x6f, 0x6d,
	0x70, 0x72, 0x65, 0x73, 0x73, 0x54, 0x79, 0x70, 0x65, 0x52, 0x0c, 0x43, 0x6f, 0x6d, 0x70, 0x72,
	0x65, 0x73, 0x73, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x43, 0x68, 0x65, 0x63, 0x6b,
	0x46, 0x6c, 0x61, 0x67, 0x73, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0a, 0x43, 0x68, 0x65,
	0x63, 0x6b, 0x46, 0x6c, 0x61, 0x67, 0x73, 0x12, 0x3a, 0x0a, 0x09, 0x54, 0x72, 0x61, 0x6e, 0x73,
	0x49, 0x6e, 0x66, 0x6f, 0x18, 0x0b, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x49,
	0x6e, 0x66, 0x6f, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x09, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x49,
	0x6e, 0x66, 0x6f, 0x1a, 0x3c, 0x0a, 0x0e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x49, 0x6e, 0x66, 0x6f,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38,
	0x01, 0x2a, 0x18, 0x0a, 0x0b, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65,
	0x12, 0x09, 0x0a, 0x05, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x10, 0x00, 0x2a, 0x18, 0x0a, 0x0c, 0x43,
	0x6f, 0x6d, 0x70, 0x72, 0x65, 0x73, 0x73, 0x54, 0x79, 0x70, 0x65, 0x12, 0x08, 0x0a, 0x04, 0x47,
	0x7a, 0x69, 0x70, 0x10, 0x00, 0x2a, 0x2b, 0x0a, 0x08, 0x43, 0x61, 0x6c, 0x6c, 0x54, 0x79, 0x70,
	0x65, 0x12, 0x0e, 0x0a, 0x0a, 0x43, 0x61, 0x6c, 0x6c, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x10,
	0x00, 0x12, 0x0f, 0x0a, 0x0b, 0x43, 0x61, 0x6c, 0x6c, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x10, 0x01, 0x3a, 0x4f, 0x0a, 0x09, 0x63, 0x61, 0x6c, 0x6c, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x12,
	0x1f, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x18, 0xd1, 0x86, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x43, 0x61, 0x6c, 0x6c, 0x54, 0x79, 0x70, 0x65, 0x52, 0x08, 0x63, 0x61, 0x6c, 0x6c, 0x54,
	0x79, 0x70, 0x65, 0x42, 0x1e, 0x5a, 0x1c, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x30, 0x78, 0x30, 0x30, 0x62, 0x2f, 0x67, 0x6f, 0x62, 0x62, 0x71, 0x3b, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_proto_rawDescOnce sync.Once
	file_proto_proto_rawDescData = file_proto_proto_rawDesc
)

func file_proto_proto_rawDescGZIP() []byte {
	file_proto_proto_rawDescOnce.Do(func() {
		file_proto_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_proto_rawDescData)
	})
	return file_proto_proto_rawDescData
}

var file_proto_proto_enumTypes = make([]protoimpl.EnumInfo, 3)
var file_proto_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_proto_proto_goTypes = []interface{}{
	(ContentType)(0),                  // 0: proto.ContentType
	(CompressType)(0),                 // 1: proto.CompressType
	(CallType)(0),                     // 2: proto.CallType
	(*Entity)(nil),                    // 3: proto.Entity
	(*Header)(nil),                    // 4: proto.Header
	nil,                               // 5: proto.Header.TransInfoEntry
	(*descriptor.ServiceOptions)(nil), // 6: google.protobuf.ServiceOptions
}
var file_proto_proto_depIdxs = []int32{
	2, // 0: proto.Header.CallType:type_name -> proto.CallType
	3, // 1: proto.Header.SrcEntity:type_name -> proto.Entity
	3, // 2: proto.Header.DstEntity:type_name -> proto.Entity
	0, // 3: proto.Header.ContentType:type_name -> proto.ContentType
	1, // 4: proto.Header.CompressType:type_name -> proto.CompressType
	5, // 5: proto.Header.TransInfo:type_name -> proto.Header.TransInfoEntry
	6, // 6: proto.call_type:extendee -> google.protobuf.ServiceOptions
	2, // 7: proto.call_type:type_name -> proto.CallType
	8, // [8:8] is the sub-list for method output_type
	8, // [8:8] is the sub-list for method input_type
	7, // [7:8] is the sub-list for extension type_name
	6, // [6:7] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_proto_proto_init() }
func file_proto_proto_init() {
	if File_proto_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Entity); i {
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
		file_proto_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Header); i {
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
			RawDescriptor: file_proto_proto_rawDesc,
			NumEnums:      3,
			NumMessages:   3,
			NumExtensions: 1,
			NumServices:   0,
		},
		GoTypes:           file_proto_proto_goTypes,
		DependencyIndexes: file_proto_proto_depIdxs,
		EnumInfos:         file_proto_proto_enumTypes,
		MessageInfos:      file_proto_proto_msgTypes,
		ExtensionInfos:    file_proto_proto_extTypes,
	}.Build()
	File_proto_proto = out.File
	file_proto_proto_rawDesc = nil
	file_proto_proto_goTypes = nil
	file_proto_proto_depIdxs = nil
}