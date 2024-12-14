// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        v3.21.12
// source: streamexample.proto

package protos

import (
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

type ServiceBusMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TargetLeaf      string `protobuf:"bytes,1,opt,name=TargetLeaf,proto3" json:"TargetLeaf,omitempty"`
	TargetProcessor string `protobuf:"bytes,2,opt,name=TargetProcessor,proto3" json:"TargetProcessor,omitempty"`
	Message         string `protobuf:"bytes,3,opt,name=Message,proto3" json:"Message,omitempty"`
	Type            string `protobuf:"bytes,4,opt,name=Type,proto3" json:"Type,omitempty"`
	MessageId       string `protobuf:"bytes,5,opt,name=MessageId,proto3" json:"MessageId,omitempty"`
}

func (x *ServiceBusMessage) Reset() {
	*x = ServiceBusMessage{}
	mi := &file_streamexample_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ServiceBusMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServiceBusMessage) ProtoMessage() {}

func (x *ServiceBusMessage) ProtoReflect() protoreflect.Message {
	mi := &file_streamexample_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServiceBusMessage.ProtoReflect.Descriptor instead.
func (*ServiceBusMessage) Descriptor() ([]byte, []int) {
	return file_streamexample_proto_rawDescGZIP(), []int{0}
}

func (x *ServiceBusMessage) GetTargetLeaf() string {
	if x != nil {
		return x.TargetLeaf
	}
	return ""
}

func (x *ServiceBusMessage) GetTargetProcessor() string {
	if x != nil {
		return x.TargetProcessor
	}
	return ""
}

func (x *ServiceBusMessage) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *ServiceBusMessage) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *ServiceBusMessage) GetMessageId() string {
	if x != nil {
		return x.MessageId
	}
	return ""
}

type ClientRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ClientName     string `protobuf:"bytes,1,opt,name=ClientName,proto3" json:"ClientName,omitempty"`
	ClientType     string `protobuf:"bytes,2,opt,name=ClientType,proto3" json:"ClientType,omitempty"`
	IsRegistration bool   `protobuf:"varint,3,opt,name=IsRegistration,proto3" json:"IsRegistration,omitempty"`
	MessageId      string `protobuf:"bytes,4,opt,name=MessageId,proto3" json:"MessageId,omitempty"`
	Error          *Error `protobuf:"bytes,5,opt,name=Error,proto3" json:"Error,omitempty"`
}

func (x *ClientRequest) Reset() {
	*x = ClientRequest{}
	mi := &file_streamexample_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ClientRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClientRequest) ProtoMessage() {}

func (x *ClientRequest) ProtoReflect() protoreflect.Message {
	mi := &file_streamexample_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClientRequest.ProtoReflect.Descriptor instead.
func (*ClientRequest) Descriptor() ([]byte, []int) {
	return file_streamexample_proto_rawDescGZIP(), []int{1}
}

func (x *ClientRequest) GetClientName() string {
	if x != nil {
		return x.ClientName
	}
	return ""
}

func (x *ClientRequest) GetClientType() string {
	if x != nil {
		return x.ClientType
	}
	return ""
}

func (x *ClientRequest) GetIsRegistration() bool {
	if x != nil {
		return x.IsRegistration
	}
	return false
}

func (x *ClientRequest) GetMessageId() string {
	if x != nil {
		return x.MessageId
	}
	return ""
}

func (x *ClientRequest) GetError() *Error {
	if x != nil {
		return x.Error
	}
	return nil
}

type Error struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message     string `protobuf:"bytes,1,opt,name=Message,proto3" json:"Message,omitempty"`
	IsRetryable bool   `protobuf:"varint,2,opt,name=IsRetryable,proto3" json:"IsRetryable,omitempty"`
}

func (x *Error) Reset() {
	*x = Error{}
	mi := &file_streamexample_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Error) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Error) ProtoMessage() {}

func (x *Error) ProtoReflect() protoreflect.Message {
	mi := &file_streamexample_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Error.ProtoReflect.Descriptor instead.
func (*Error) Descriptor() ([]byte, []int) {
	return file_streamexample_proto_rawDescGZIP(), []int{2}
}

func (x *Error) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *Error) GetIsRetryable() bool {
	if x != nil {
		return x.IsRetryable
	}
	return false
}

var File_streamexample_proto protoreflect.FileDescriptor

var file_streamexample_proto_rawDesc = []byte{
	0x0a, 0x13, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x22, 0xa9, 0x01,
	0x0a, 0x11, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x42, 0x75, 0x73, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x54, 0x61, 0x72, 0x67, 0x65, 0x74, 0x4c, 0x65, 0x61,
	0x66, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x54, 0x61, 0x72, 0x67, 0x65, 0x74, 0x4c,
	0x65, 0x61, 0x66, 0x12, 0x28, 0x0a, 0x0f, 0x54, 0x61, 0x72, 0x67, 0x65, 0x74, 0x50, 0x72, 0x6f,
	0x63, 0x65, 0x73, 0x73, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x54, 0x61,
	0x72, 0x67, 0x65, 0x74, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x6f, 0x72, 0x12, 0x18, 0x0a,
	0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x54, 0x79, 0x70, 0x65, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x49, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x49, 0x64, 0x22, 0xba, 0x01, 0x0a, 0x0d, 0x43, 0x6c,
	0x69, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1e, 0x0a, 0x0a, 0x43,
	0x6c, 0x69, 0x65, 0x6e, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0a, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x43,
	0x6c, 0x69, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0a, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x26, 0x0a, 0x0e, 0x49,
	0x73, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x0e, 0x49, 0x73, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x12, 0x1c, 0x0a, 0x09, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x49, 0x64,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x49,
	0x64, 0x12, 0x23, 0x0a, 0x05, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x0d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x52,
	0x05, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x43, 0x0a, 0x05, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x12,
	0x18, 0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x49, 0x73, 0x52,
	0x65, 0x74, 0x72, 0x79, 0x61, 0x62, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b,
	0x49, 0x73, 0x52, 0x65, 0x74, 0x72, 0x79, 0x61, 0x62, 0x6c, 0x65, 0x32, 0x56, 0x0a, 0x0d, 0x53,
	0x74, 0x72, 0x65, 0x61, 0x6d, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x45, 0x0a, 0x0b,
	0x47, 0x65, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x12, 0x15, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x73, 0x2e, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x19, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2e, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x42, 0x75, 0x73, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x00, 0x28,
	0x01, 0x30, 0x01, 0x42, 0x2f, 0x5a, 0x2d, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x69, 0x6d, 0x69, 0x6c, 0x6c, 0x65, 0x72, 0x33, 0x31, 0x2f, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x62, 0x75, 0x73, 0x2d, 0x66, 0x61, 0x6e, 0x6f, 0x75, 0x74, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_streamexample_proto_rawDescOnce sync.Once
	file_streamexample_proto_rawDescData = file_streamexample_proto_rawDesc
)

func file_streamexample_proto_rawDescGZIP() []byte {
	file_streamexample_proto_rawDescOnce.Do(func() {
		file_streamexample_proto_rawDescData = protoimpl.X.CompressGZIP(file_streamexample_proto_rawDescData)
	})
	return file_streamexample_proto_rawDescData
}

var file_streamexample_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_streamexample_proto_goTypes = []any{
	(*ServiceBusMessage)(nil), // 0: protos.ServiceBusMessage
	(*ClientRequest)(nil),     // 1: protos.ClientRequest
	(*Error)(nil),             // 2: protos.Error
}
var file_streamexample_proto_depIdxs = []int32{
	2, // 0: protos.ClientRequest.Error:type_name -> protos.Error
	1, // 1: protos.StreamService.GetMessages:input_type -> protos.ClientRequest
	0, // 2: protos.StreamService.GetMessages:output_type -> protos.ServiceBusMessage
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_streamexample_proto_init() }
func file_streamexample_proto_init() {
	if File_streamexample_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_streamexample_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_streamexample_proto_goTypes,
		DependencyIndexes: file_streamexample_proto_depIdxs,
		MessageInfos:      file_streamexample_proto_msgTypes,
	}.Build()
	File_streamexample_proto = out.File
	file_streamexample_proto_rawDesc = nil
	file_streamexample_proto_goTypes = nil
	file_streamexample_proto_depIdxs = nil
}
