// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.6.1
// source: bbq.proto

package bbq

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
	return file_bbq_proto_enumTypes[0].Descriptor()
}

func (ContentType) Type() protoreflect.EnumType {
	return &file_bbq_proto_enumTypes[0]
}

func (x ContentType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ContentType.Descriptor instead.
func (ContentType) EnumDescriptor() ([]byte, []int) {
	return file_bbq_proto_rawDescGZIP(), []int{0}
}

type CompressType int32

const (
	CompressType_None CompressType = 0
	CompressType_Gzip CompressType = 1
)

// Enum value maps for CompressType.
var (
	CompressType_name = map[int32]string{
		0: "None",
		1: "Gzip",
	}
	CompressType_value = map[string]int32{
		"None": 0,
		"Gzip": 1,
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
	return file_bbq_proto_enumTypes[1].Descriptor()
}

func (CompressType) Type() protoreflect.EnumType {
	return &file_bbq_proto_enumTypes[1]
}

func (x CompressType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use CompressType.Descriptor instead.
func (CompressType) EnumDescriptor() ([]byte, []int) {
	return file_bbq_proto_rawDescGZIP(), []int{1}
}

type RequestType int32

const (
	RequestType_RequestRequest RequestType = 0
	RequestType_RequestRespone RequestType = 1
)

// Enum value maps for RequestType.
var (
	RequestType_name = map[int32]string{
		0: "RequestRequest",
		1: "RequestRespone",
	}
	RequestType_value = map[string]int32{
		"RequestRequest": 0,
		"RequestRespone": 1,
	}
)

func (x RequestType) Enum() *RequestType {
	p := new(RequestType)
	*p = x
	return p
}

func (x RequestType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (RequestType) Descriptor() protoreflect.EnumDescriptor {
	return file_bbq_proto_enumTypes[2].Descriptor()
}

func (RequestType) Type() protoreflect.EnumType {
	return &file_bbq_proto_enumTypes[2]
}

func (x RequestType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use RequestType.Descriptor instead.
func (RequestType) EnumDescriptor() ([]byte, []int) {
	return file_bbq_proto_rawDescGZIP(), []int{2}
}

type ServiceType int32

const (
	// 请求entity，需要提供entity id， entity是有ID的service, entity可以创建很多
	ServiceType_Entity ServiceType = 0
	// 请求service，只需要提供完整接口名，service是单例entity，只能有一个
	ServiceType_Service ServiceType = 1
	// 系统接口，proxyid+instid
	ServiceType_System ServiceType = 2
)

// Enum value maps for ServiceType.
var (
	ServiceType_name = map[int32]string{
		0: "Entity",
		1: "Service",
		2: "System",
	}
	ServiceType_value = map[string]int32{
		"Entity":  0,
		"Service": 1,
		"System":  2,
	}
)

func (x ServiceType) Enum() *ServiceType {
	p := new(ServiceType)
	*p = x
	return p
}

func (x ServiceType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ServiceType) Descriptor() protoreflect.EnumDescriptor {
	return file_bbq_proto_enumTypes[3].Descriptor()
}

func (ServiceType) Type() protoreflect.EnumType {
	return &file_bbq_proto_enumTypes[3]
}

func (x ServiceType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ServiceType.Descriptor instead.
func (ServiceType) EnumDescriptor() ([]byte, []int) {
	return file_bbq_proto_rawDescGZIP(), []int{3}
}

// 请求协议头
type Header struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 协议版本
	Version uint32 `protobuf:"varint,1,opt,name=Version,proto3" json:"Version,omitempty"`
	// 请求唯一id
	RequestId string `protobuf:"bytes,2,opt,name=RequestId,proto3" json:"RequestId,omitempty"`
	// 请求的超时时间，单位ms
	Timeout uint32 `protobuf:"varint,3,opt,name=Timeout,proto3" json:"Timeout,omitempty"`
	// 是请求包，还是返回包
	RequestType RequestType `protobuf:"varint,4,opt,name=RequestType,proto3,enum=bbq.RequestType" json:"RequestType,omitempty"`
	// sverice or entity
	ServiceType ServiceType `protobuf:"varint,5,opt,name=ServiceType,proto3,enum=bbq.ServiceType" json:"ServiceType,omitempty"`
	// 调用的原EntityID
	SrcEntity uint64 `protobuf:"varint,6,opt,name=SrcEntity,proto3" json:"SrcEntity,omitempty"`
	// 调用的目的EntityID
	DstEntity uint64 `protobuf:"varint,7,opt,name=DstEntity,proto3" json:"DstEntity,omitempty"`
	// 规范格式: 类名，服务名
	Type string `protobuf:"bytes,8,opt,name=Type,proto3" json:"Type,omitempty"`
	// 规范格式: 接口名
	Method string `protobuf:"bytes,9,opt,name=Method,proto3" json:"Method,omitempty"`
	// 请求数据的序列化类型
	// 比如: proto/jce/json, 默认proto
	// 具体值与ContentEncodeType对应
	ContentType ContentType `protobuf:"varint,10,opt,name=ContentType,proto3,enum=bbq.ContentType" json:"ContentType,omitempty"`
	// 请求数据使用的压缩方式
	// 比如: gzip/snappy/..., 默认不使用
	// 具体值与CompressType对应
	CompressType CompressType `protobuf:"varint,11,opt,name=CompressType,proto3,enum=bbq.CompressType" json:"CompressType,omitempty"`
	// 是否检查包是否正确
	CheckFlags uint32 `protobuf:"varint,12,opt,name=CheckFlags,proto3" json:"CheckFlags,omitempty"`
	// 附加信息
	TransInfo map[string][]byte `protobuf:"bytes,13,rep,name=TransInfo,proto3" json:"TransInfo,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// 返回值
	ErrCode int32  `protobuf:"varint,14,opt,name=ErrCode,proto3" json:"ErrCode,omitempty"`
	ErrMsg  string `protobuf:"bytes,15,opt,name=ErrMsg,proto3" json:"ErrMsg,omitempty"`
}

func (x *Header) Reset() {
	*x = Header{}
	if protoimpl.UnsafeEnabled {
		mi := &file_bbq_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Header) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Header) ProtoMessage() {}

func (x *Header) ProtoReflect() protoreflect.Message {
	mi := &file_bbq_proto_msgTypes[0]
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
	return file_bbq_proto_rawDescGZIP(), []int{0}
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

func (x *Header) GetRequestType() RequestType {
	if x != nil {
		return x.RequestType
	}
	return RequestType_RequestRequest
}

func (x *Header) GetServiceType() ServiceType {
	if x != nil {
		return x.ServiceType
	}
	return ServiceType_Entity
}

func (x *Header) GetSrcEntity() uint64 {
	if x != nil {
		return x.SrcEntity
	}
	return 0
}

func (x *Header) GetDstEntity() uint64 {
	if x != nil {
		return x.DstEntity
	}
	return 0
}

func (x *Header) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
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
	return CompressType_None
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

func (x *Header) GetErrCode() int32 {
	if x != nil {
		return x.ErrCode
	}
	return 0
}

func (x *Header) GetErrMsg() string {
	if x != nil {
		return x.ErrMsg
	}
	return ""
}

var file_bbq_proto_extTypes = []protoimpl.ExtensionInfo{
	{
		ExtendedType:  (*descriptor.ServiceOptions)(nil),
		ExtensionType: (*ServiceType)(nil),
		Field:         50100,
		Name:          "bbq.service_type",
		Tag:           "varint,50100,opt,name=service_type,enum=bbq.ServiceType",
		Filename:      "bbq.proto",
	},
}

// Extension fields to descriptor.ServiceOptions.
var (
	// optional bbq.ServiceType service_type = 50100;
	E_ServiceType = &file_bbq_proto_extTypes[0]
)

var File_bbq_proto protoreflect.FileDescriptor

var file_bbq_proto_rawDesc = []byte{
	0x0a, 0x09, 0x62, 0x62, 0x71, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x03, 0x62, 0x62, 0x71,
	0x1a, 0x20, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0xdf, 0x04, 0x0a, 0x06, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x12, 0x18, 0x0a,
	0x07, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x07,
	0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x1c, 0x0a, 0x09, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x54, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x07, 0x54, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x12,
	0x32, 0x0a, 0x0b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x54, 0x79, 0x70, 0x65, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x10, 0x2e, 0x62, 0x62, 0x71, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x54, 0x79, 0x70, 0x65, 0x52, 0x0b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x54,
	0x79, 0x70, 0x65, 0x12, 0x32, 0x0a, 0x0b, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x54, 0x79,
	0x70, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x10, 0x2e, 0x62, 0x62, 0x71, 0x2e, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x54, 0x79, 0x70, 0x65, 0x52, 0x0b, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x53, 0x72, 0x63, 0x45, 0x6e,
	0x74, 0x69, 0x74, 0x79, 0x18, 0x06, 0x20, 0x01, 0x28, 0x04, 0x52, 0x09, 0x53, 0x72, 0x63, 0x45,
	0x6e, 0x74, 0x69, 0x74, 0x79, 0x12, 0x1c, 0x0a, 0x09, 0x44, 0x73, 0x74, 0x45, 0x6e, 0x74, 0x69,
	0x74, 0x79, 0x18, 0x07, 0x20, 0x01, 0x28, 0x04, 0x52, 0x09, 0x44, 0x73, 0x74, 0x45, 0x6e, 0x74,
	0x69, 0x74, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x54, 0x79, 0x70, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x54, 0x79, 0x70, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x4d, 0x65, 0x74, 0x68, 0x6f,
	0x64, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x12,
	0x32, 0x0a, 0x0b, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x18, 0x0a,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x10, 0x2e, 0x62, 0x62, 0x71, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x65,
	0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x52, 0x0b, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x54,
	0x79, 0x70, 0x65, 0x12, 0x35, 0x0a, 0x0c, 0x43, 0x6f, 0x6d, 0x70, 0x72, 0x65, 0x73, 0x73, 0x54,
	0x79, 0x70, 0x65, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x11, 0x2e, 0x62, 0x62, 0x71, 0x2e,
	0x43, 0x6f, 0x6d, 0x70, 0x72, 0x65, 0x73, 0x73, 0x54, 0x79, 0x70, 0x65, 0x52, 0x0c, 0x43, 0x6f,
	0x6d, 0x70, 0x72, 0x65, 0x73, 0x73, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x43, 0x68,
	0x65, 0x63, 0x6b, 0x46, 0x6c, 0x61, 0x67, 0x73, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0a,
	0x43, 0x68, 0x65, 0x63, 0x6b, 0x46, 0x6c, 0x61, 0x67, 0x73, 0x12, 0x38, 0x0a, 0x09, 0x54, 0x72,
	0x61, 0x6e, 0x73, 0x49, 0x6e, 0x66, 0x6f, 0x18, 0x0d, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x62, 0x62, 0x71, 0x2e, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73,
	0x49, 0x6e, 0x66, 0x6f, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x09, 0x54, 0x72, 0x61, 0x6e, 0x73,
	0x49, 0x6e, 0x66, 0x6f, 0x12, 0x18, 0x0a, 0x07, 0x45, 0x72, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x18,
	0x0e, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x45, 0x72, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x16,
	0x0a, 0x06, 0x45, 0x72, 0x72, 0x4d, 0x73, 0x67, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x45, 0x72, 0x72, 0x4d, 0x73, 0x67, 0x1a, 0x3c, 0x0a, 0x0e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x49,
	0x6e, 0x66, 0x6f, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x3a, 0x02, 0x38, 0x01, 0x2a, 0x18, 0x0a, 0x0b, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x54,
	0x79, 0x70, 0x65, 0x12, 0x09, 0x0a, 0x05, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x10, 0x00, 0x2a, 0x22,
	0x0a, 0x0c, 0x43, 0x6f, 0x6d, 0x70, 0x72, 0x65, 0x73, 0x73, 0x54, 0x79, 0x70, 0x65, 0x12, 0x08,
	0x0a, 0x04, 0x4e, 0x6f, 0x6e, 0x65, 0x10, 0x00, 0x12, 0x08, 0x0a, 0x04, 0x47, 0x7a, 0x69, 0x70,
	0x10, 0x01, 0x2a, 0x35, 0x0a, 0x0b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x54, 0x79, 0x70,
	0x65, 0x12, 0x12, 0x0a, 0x0e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x10, 0x00, 0x12, 0x12, 0x0a, 0x0e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x65, 0x10, 0x01, 0x2a, 0x32, 0x0a, 0x0b, 0x53, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0a, 0x0a, 0x06, 0x45, 0x6e, 0x74, 0x69,
	0x74, 0x79, 0x10, 0x00, 0x12, 0x0b, 0x0a, 0x07, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x10,
	0x01, 0x12, 0x0a, 0x0a, 0x06, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x10, 0x02, 0x3a, 0x56, 0x0a,
	0x0c, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x12, 0x1f, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xb4,
	0x87, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x10, 0x2e, 0x62, 0x62, 0x71, 0x2e, 0x53, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x54, 0x79, 0x70, 0x65, 0x52, 0x0b, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x54, 0x79, 0x70, 0x65, 0x42, 0x26, 0x5a, 0x24, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x30, 0x78, 0x30, 0x30, 0x62, 0x2f, 0x67, 0x6f, 0x62, 0x62, 0x71, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x62, 0x62, 0x71, 0x3b, 0x62, 0x62, 0x71, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_bbq_proto_rawDescOnce sync.Once
	file_bbq_proto_rawDescData = file_bbq_proto_rawDesc
)

func file_bbq_proto_rawDescGZIP() []byte {
	file_bbq_proto_rawDescOnce.Do(func() {
		file_bbq_proto_rawDescData = protoimpl.X.CompressGZIP(file_bbq_proto_rawDescData)
	})
	return file_bbq_proto_rawDescData
}

var file_bbq_proto_enumTypes = make([]protoimpl.EnumInfo, 4)
var file_bbq_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_bbq_proto_goTypes = []interface{}{
	(ContentType)(0),                  // 0: bbq.ContentType
	(CompressType)(0),                 // 1: bbq.CompressType
	(RequestType)(0),                  // 2: bbq.RequestType
	(ServiceType)(0),                  // 3: bbq.ServiceType
	(*Header)(nil),                    // 4: bbq.Header
	nil,                               // 5: bbq.Header.TransInfoEntry
	(*descriptor.ServiceOptions)(nil), // 6: google.protobuf.ServiceOptions
}
var file_bbq_proto_depIdxs = []int32{
	2, // 0: bbq.Header.RequestType:type_name -> bbq.RequestType
	3, // 1: bbq.Header.ServiceType:type_name -> bbq.ServiceType
	0, // 2: bbq.Header.ContentType:type_name -> bbq.ContentType
	1, // 3: bbq.Header.CompressType:type_name -> bbq.CompressType
	5, // 4: bbq.Header.TransInfo:type_name -> bbq.Header.TransInfoEntry
	6, // 5: bbq.service_type:extendee -> google.protobuf.ServiceOptions
	3, // 6: bbq.service_type:type_name -> bbq.ServiceType
	7, // [7:7] is the sub-list for method output_type
	7, // [7:7] is the sub-list for method input_type
	6, // [6:7] is the sub-list for extension type_name
	5, // [5:6] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_bbq_proto_init() }
func file_bbq_proto_init() {
	if File_bbq_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_bbq_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
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
			RawDescriptor: file_bbq_proto_rawDesc,
			NumEnums:      4,
			NumMessages:   2,
			NumExtensions: 1,
			NumServices:   0,
		},
		GoTypes:           file_bbq_proto_goTypes,
		DependencyIndexes: file_bbq_proto_depIdxs,
		EnumInfos:         file_bbq_proto_enumTypes,
		MessageInfos:      file_bbq_proto_msgTypes,
		ExtensionInfos:    file_bbq_proto_extTypes,
	}.Build()
	File_bbq_proto = out.File
	file_bbq_proto_rawDesc = nil
	file_bbq_proto_goTypes = nil
	file_bbq_proto_depIdxs = nil
}
