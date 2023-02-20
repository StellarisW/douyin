// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.21.1
// source: chat.sys.proto

package pb

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

type Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id         int64  `protobuf:"varint,1,opt,name=id,proto3" json:"id"`                                     // 消息id
	ToUserId   int64  `protobuf:"varint,2,opt,name=to_user_id,json=toUserId,proto3" json:"to_user_id"`       // 消息接收者id
	FromUserId int64  `protobuf:"varint,3,opt,name=from_user_id,json=fromUserId,proto3" json:"from_user_id"` // 消息发送者id
	Content    string `protobuf:"bytes,4,opt,name=content,proto3" json:"content"`                            // 消息内容
	CreateTime int64  `protobuf:"varint,5,opt,name=create_time,json=createTime,proto3" json:"create_time"`   // 消息创建时间
}

func (x *Message) Reset() {
	*x = Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chat_sys_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_chat_sys_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message.ProtoReflect.Descriptor instead.
func (*Message) Descriptor() ([]byte, []int) {
	return file_chat_sys_proto_rawDescGZIP(), []int{0}
}

func (x *Message) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Message) GetToUserId() int64 {
	if x != nil {
		return x.ToUserId
	}
	return 0
}

func (x *Message) GetFromUserId() int64 {
	if x != nil {
		return x.FromUserId
	}
	return 0
}

func (x *Message) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

func (x *Message) GetCreateTime() int64 {
	if x != nil {
		return x.CreateTime
	}
	return 0
}

type SendMessageReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SrcUserId  int64  `protobuf:"varint,1,opt,name=src_user_id,json=srcUserId,proto3" json:"src_user_id"`
	DstUserId  int64  `protobuf:"varint,2,opt,name=dst_user_id,json=dstUserId,proto3" json:"dst_user_id"`
	ActionType uint32 `protobuf:"varint,3,opt,name=action_type,json=actionType,proto3" json:"action_type"`
	Content    string `protobuf:"bytes,4,opt,name=content,proto3" json:"content"`
}

func (x *SendMessageReq) Reset() {
	*x = SendMessageReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chat_sys_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendMessageReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendMessageReq) ProtoMessage() {}

func (x *SendMessageReq) ProtoReflect() protoreflect.Message {
	mi := &file_chat_sys_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendMessageReq.ProtoReflect.Descriptor instead.
func (*SendMessageReq) Descriptor() ([]byte, []int) {
	return file_chat_sys_proto_rawDescGZIP(), []int{1}
}

func (x *SendMessageReq) GetSrcUserId() int64 {
	if x != nil {
		return x.SrcUserId
	}
	return 0
}

func (x *SendMessageReq) GetDstUserId() int64 {
	if x != nil {
		return x.DstUserId
	}
	return 0
}

func (x *SendMessageReq) GetActionType() uint32 {
	if x != nil {
		return x.ActionType
	}
	return 0
}

func (x *SendMessageReq) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

type SendMessageRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StatusCode uint32 `protobuf:"varint,1,opt,name=status_code,json=statusCode,proto3" json:"status_code"`
	StatusMsg  string `protobuf:"bytes,2,opt,name=status_msg,json=statusMsg,proto3" json:"status_msg"`
}

func (x *SendMessageRes) Reset() {
	*x = SendMessageRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chat_sys_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendMessageRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendMessageRes) ProtoMessage() {}

func (x *SendMessageRes) ProtoReflect() protoreflect.Message {
	mi := &file_chat_sys_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendMessageRes.ProtoReflect.Descriptor instead.
func (*SendMessageRes) Descriptor() ([]byte, []int) {
	return file_chat_sys_proto_rawDescGZIP(), []int{2}
}

func (x *SendMessageRes) GetStatusCode() uint32 {
	if x != nil {
		return x.StatusCode
	}
	return 0
}

func (x *SendMessageRes) GetStatusMsg() string {
	if x != nil {
		return x.StatusMsg
	}
	return ""
}

type GetMessageReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SrcUserId int64 `protobuf:"varint,1,opt,name=src_user_id,json=srcUserId,proto3" json:"src_user_id"`
	DstUserId int64 `protobuf:"varint,2,opt,name=dst_user_id,json=dstUserId,proto3" json:"dst_user_id"`
}

func (x *GetMessageReq) Reset() {
	*x = GetMessageReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chat_sys_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetMessageReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetMessageReq) ProtoMessage() {}

func (x *GetMessageReq) ProtoReflect() protoreflect.Message {
	mi := &file_chat_sys_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetMessageReq.ProtoReflect.Descriptor instead.
func (*GetMessageReq) Descriptor() ([]byte, []int) {
	return file_chat_sys_proto_rawDescGZIP(), []int{3}
}

func (x *GetMessageReq) GetSrcUserId() int64 {
	if x != nil {
		return x.SrcUserId
	}
	return 0
}

func (x *GetMessageReq) GetDstUserId() int64 {
	if x != nil {
		return x.DstUserId
	}
	return 0
}

type GetMessageRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StatusCode uint32     `protobuf:"varint,1,opt,name=status_code,json=statusCode,proto3" json:"status_code"`
	StatusMsg  string     `protobuf:"bytes,2,opt,name=status_msg,json=statusMsg,proto3" json:"status_msg"`
	Messages   []*Message `protobuf:"bytes,3,rep,name=messages,proto3" json:"messages"`
}

func (x *GetMessageRes) Reset() {
	*x = GetMessageRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chat_sys_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetMessageRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetMessageRes) ProtoMessage() {}

func (x *GetMessageRes) ProtoReflect() protoreflect.Message {
	mi := &file_chat_sys_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetMessageRes.ProtoReflect.Descriptor instead.
func (*GetMessageRes) Descriptor() ([]byte, []int) {
	return file_chat_sys_proto_rawDescGZIP(), []int{4}
}

func (x *GetMessageRes) GetStatusCode() uint32 {
	if x != nil {
		return x.StatusCode
	}
	return 0
}

func (x *GetMessageRes) GetStatusMsg() string {
	if x != nil {
		return x.StatusMsg
	}
	return ""
}

func (x *GetMessageRes) GetMessages() []*Message {
	if x != nil {
		return x.Messages
	}
	return nil
}

type StoreMessageReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SrcUserId int64  `protobuf:"varint,1,opt,name=src_user_id,json=srcUserId,proto3" json:"src_user_id"`
	DstUserId int64  `protobuf:"varint,2,opt,name=dst_user_id,json=dstUserId,proto3" json:"dst_user_id"`
	Content   string `protobuf:"bytes,3,opt,name=content,proto3" json:"content"`
}

func (x *StoreMessageReq) Reset() {
	*x = StoreMessageReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chat_sys_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StoreMessageReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StoreMessageReq) ProtoMessage() {}

func (x *StoreMessageReq) ProtoReflect() protoreflect.Message {
	mi := &file_chat_sys_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StoreMessageReq.ProtoReflect.Descriptor instead.
func (*StoreMessageReq) Descriptor() ([]byte, []int) {
	return file_chat_sys_proto_rawDescGZIP(), []int{5}
}

func (x *StoreMessageReq) GetSrcUserId() int64 {
	if x != nil {
		return x.SrcUserId
	}
	return 0
}

func (x *StoreMessageReq) GetDstUserId() int64 {
	if x != nil {
		return x.DstUserId
	}
	return 0
}

func (x *StoreMessageReq) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

type StoreMessageRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StatusCode uint32 `protobuf:"varint,1,opt,name=status_code,json=statusCode,proto3" json:"status_code"`
	StatusMsg  string `protobuf:"bytes,2,opt,name=status_msg,json=statusMsg,proto3" json:"status_msg"`
}

func (x *StoreMessageRes) Reset() {
	*x = StoreMessageRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chat_sys_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StoreMessageRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StoreMessageRes) ProtoMessage() {}

func (x *StoreMessageRes) ProtoReflect() protoreflect.Message {
	mi := &file_chat_sys_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StoreMessageRes.ProtoReflect.Descriptor instead.
func (*StoreMessageRes) Descriptor() ([]byte, []int) {
	return file_chat_sys_proto_rawDescGZIP(), []int{6}
}

func (x *StoreMessageRes) GetStatusCode() uint32 {
	if x != nil {
		return x.StatusCode
	}
	return 0
}

func (x *StoreMessageRes) GetStatusMsg() string {
	if x != nil {
		return x.StatusMsg
	}
	return ""
}

var File_chat_sys_proto protoreflect.FileDescriptor

var file_chat_sys_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x73, 0x79, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x08, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x73, 0x79, 0x73, 0x22, 0x94, 0x01, 0x0a, 0x07, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1c, 0x0a, 0x0a, 0x74, 0x6f, 0x5f, 0x75, 0x73, 0x65,
	0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x74, 0x6f, 0x55, 0x73,
	0x65, 0x72, 0x49, 0x64, 0x12, 0x20, 0x0a, 0x0c, 0x66, 0x72, 0x6f, 0x6d, 0x5f, 0x75, 0x73, 0x65,
	0x72, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x66, 0x72, 0x6f, 0x6d,
	0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e,
	0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74,
	0x12, 0x1f, 0x0a, 0x0b, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x69, 0x6d,
	0x65, 0x22, 0x8b, 0x01, 0x0a, 0x0e, 0x53, 0x65, 0x6e, 0x64, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x52, 0x65, 0x71, 0x12, 0x1e, 0x0a, 0x0b, 0x73, 0x72, 0x63, 0x5f, 0x75, 0x73, 0x65, 0x72,
	0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x73, 0x72, 0x63, 0x55, 0x73,
	0x65, 0x72, 0x49, 0x64, 0x12, 0x1e, 0x0a, 0x0b, 0x64, 0x73, 0x74, 0x5f, 0x75, 0x73, 0x65, 0x72,
	0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x64, 0x73, 0x74, 0x55, 0x73,
	0x65, 0x72, 0x49, 0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74,
	0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0a, 0x61, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x54, 0x79, 0x70, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x22,
	0x50, 0x0a, 0x0e, 0x53, 0x65, 0x6e, 0x64, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65,
	0x73, 0x12, 0x1f, 0x0a, 0x0b, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x5f, 0x63, 0x6f, 0x64, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0a, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x43, 0x6f,
	0x64, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x5f, 0x6d, 0x73, 0x67,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x4d, 0x73,
	0x67, 0x22, 0x4f, 0x0a, 0x0d, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52,
	0x65, 0x71, 0x12, 0x1e, 0x0a, 0x0b, 0x73, 0x72, 0x63, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x73, 0x72, 0x63, 0x55, 0x73, 0x65, 0x72,
	0x49, 0x64, 0x12, 0x1e, 0x0a, 0x0b, 0x64, 0x73, 0x74, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x64, 0x73, 0x74, 0x55, 0x73, 0x65, 0x72,
	0x49, 0x64, 0x22, 0x7e, 0x0a, 0x0d, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x52, 0x65, 0x73, 0x12, 0x1f, 0x0a, 0x0b, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x5f, 0x63, 0x6f,
	0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0a, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x43, 0x6f, 0x64, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x5f, 0x6d,
	0x73, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x4d, 0x73, 0x67, 0x12, 0x2d, 0x0a, 0x08, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x18,
	0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x73, 0x79, 0x73,
	0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x08, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x73, 0x22, 0x6b, 0x0a, 0x0f, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x52, 0x65, 0x71, 0x12, 0x1e, 0x0a, 0x0b, 0x73, 0x72, 0x63, 0x5f, 0x75, 0x73, 0x65,
	0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x73, 0x72, 0x63, 0x55,
	0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x1e, 0x0a, 0x0b, 0x64, 0x73, 0x74, 0x5f, 0x75, 0x73, 0x65,
	0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x64, 0x73, 0x74, 0x55,
	0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x22,
	0x51, 0x0a, 0x0f, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52,
	0x65, 0x73, 0x12, 0x1f, 0x0a, 0x0b, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x5f, 0x63, 0x6f, 0x64,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0a, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x43,
	0x6f, 0x64, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x5f, 0x6d, 0x73,
	0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x4d,
	0x73, 0x67, 0x32, 0xce, 0x01, 0x0a, 0x03, 0x53, 0x79, 0x73, 0x12, 0x41, 0x0a, 0x0b, 0x53, 0x65,
	0x6e, 0x64, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x18, 0x2e, 0x63, 0x68, 0x61, 0x74,
	0x2e, 0x73, 0x79, 0x73, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x52, 0x65, 0x71, 0x1a, 0x18, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x73, 0x79, 0x73, 0x2e, 0x53,
	0x65, 0x6e, 0x64, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x73, 0x12, 0x3e, 0x0a,
	0x0a, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x17, 0x2e, 0x63, 0x68,
	0x61, 0x74, 0x2e, 0x73, 0x79, 0x73, 0x2e, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x52, 0x65, 0x71, 0x1a, 0x17, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x73, 0x79, 0x73, 0x2e,
	0x47, 0x65, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x73, 0x12, 0x44, 0x0a,
	0x0c, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x19, 0x2e,
	0x63, 0x68, 0x61, 0x74, 0x2e, 0x73, 0x79, 0x73, 0x2e, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x71, 0x1a, 0x19, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x2e,
	0x73, 0x79, 0x73, 0x2e, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x52, 0x65, 0x73, 0x42, 0x06, 0x5a, 0x04, 0x2e, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_chat_sys_proto_rawDescOnce sync.Once
	file_chat_sys_proto_rawDescData = file_chat_sys_proto_rawDesc
)

func file_chat_sys_proto_rawDescGZIP() []byte {
	file_chat_sys_proto_rawDescOnce.Do(func() {
		file_chat_sys_proto_rawDescData = protoimpl.X.CompressGZIP(file_chat_sys_proto_rawDescData)
	})
	return file_chat_sys_proto_rawDescData
}

var file_chat_sys_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_chat_sys_proto_goTypes = []interface{}{
	(*Message)(nil),         // 0: chat.sys.Message
	(*SendMessageReq)(nil),  // 1: chat.sys.SendMessageReq
	(*SendMessageRes)(nil),  // 2: chat.sys.SendMessageRes
	(*GetMessageReq)(nil),   // 3: chat.sys.GetMessageReq
	(*GetMessageRes)(nil),   // 4: chat.sys.GetMessageRes
	(*StoreMessageReq)(nil), // 5: chat.sys.StoreMessageReq
	(*StoreMessageRes)(nil), // 6: chat.sys.StoreMessageRes
}
var file_chat_sys_proto_depIdxs = []int32{
	0, // 0: chat.sys.GetMessageRes.messages:type_name -> chat.sys.Message
	1, // 1: chat.sys.Sys.SendMessage:input_type -> chat.sys.SendMessageReq
	3, // 2: chat.sys.Sys.GetMessage:input_type -> chat.sys.GetMessageReq
	5, // 3: chat.sys.Sys.StoreMessage:input_type -> chat.sys.StoreMessageReq
	2, // 4: chat.sys.Sys.SendMessage:output_type -> chat.sys.SendMessageRes
	4, // 5: chat.sys.Sys.GetMessage:output_type -> chat.sys.GetMessageRes
	6, // 6: chat.sys.Sys.StoreMessage:output_type -> chat.sys.StoreMessageRes
	4, // [4:7] is the sub-list for method output_type
	1, // [1:4] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_chat_sys_proto_init() }
func file_chat_sys_proto_init() {
	if File_chat_sys_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_chat_sys_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Message); i {
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
		file_chat_sys_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendMessageReq); i {
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
		file_chat_sys_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendMessageRes); i {
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
		file_chat_sys_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetMessageReq); i {
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
		file_chat_sys_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetMessageRes); i {
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
		file_chat_sys_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StoreMessageReq); i {
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
		file_chat_sys_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StoreMessageRes); i {
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
			RawDescriptor: file_chat_sys_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_chat_sys_proto_goTypes,
		DependencyIndexes: file_chat_sys_proto_depIdxs,
		MessageInfos:      file_chat_sys_proto_msgTypes,
	}.Build()
	File_chat_sys_proto = out.File
	file_chat_sys_proto_rawDesc = nil
	file_chat_sys_proto_goTypes = nil
	file_chat_sys_proto_depIdxs = nil
}
