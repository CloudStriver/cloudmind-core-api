// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.19.4
// source: cloudmind/core_api/system.proto

package core_api

import (
	_ "github.com/CloudStriver/cloudmind-core-api/biz/application/dto/basic"
	_ "github.com/CloudStriver/cloudmind-core-api/biz/application/dto/http"
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

type GetNotificationsReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OnlyType  *int64  `protobuf:"varint,1,opt,name=onlyType,proto3,oneof" json:"onlyType" form:"onlyType" query:"onlyType"`
	Limit     *int64  `protobuf:"varint,2,opt,name=limit,proto3,oneof" json:"limit" form:"limit" query:"limit"`
	LastToken *string `protobuf:"bytes,3,opt,name=lastToken,proto3,oneof" json:"lastToken" form:"lastToken" query:"lastToken"`
	Backward  *bool   `protobuf:"varint,4,opt,name=backward,proto3,oneof" json:"backward" form:"backward" query:"backward"`
	Offset    *int64  `protobuf:"varint,5,opt,name=offset,proto3,oneof" json:"offset" form:"offset" query:"offset"`
}

func (x *GetNotificationsReq) Reset() {
	*x = GetNotificationsReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cloudmind_core_api_system_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetNotificationsReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetNotificationsReq) ProtoMessage() {}

func (x *GetNotificationsReq) ProtoReflect() protoreflect.Message {
	mi := &file_cloudmind_core_api_system_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetNotificationsReq.ProtoReflect.Descriptor instead.
func (*GetNotificationsReq) Descriptor() ([]byte, []int) {
	return file_cloudmind_core_api_system_proto_rawDescGZIP(), []int{0}
}

func (x *GetNotificationsReq) GetOnlyType() int64 {
	if x != nil && x.OnlyType != nil {
		return *x.OnlyType
	}
	return 0
}

func (x *GetNotificationsReq) GetLimit() int64 {
	if x != nil && x.Limit != nil {
		return *x.Limit
	}
	return 0
}

func (x *GetNotificationsReq) GetLastToken() string {
	if x != nil && x.LastToken != nil {
		return *x.LastToken
	}
	return ""
}

func (x *GetNotificationsReq) GetBackward() bool {
	if x != nil && x.Backward != nil {
		return *x.Backward
	}
	return false
}

func (x *GetNotificationsReq) GetOffset() int64 {
	if x != nil && x.Offset != nil {
		return *x.Offset
	}
	return 0
}

type GetNotificationsResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Notifications []*Notification `protobuf:"bytes,1,rep,name=notifications,proto3" json:"notifications" form:"notifications" query:"notifications"`
	Token         string          `protobuf:"bytes,2,opt,name=token,proto3" json:"token" form:"token" query:"token"`
}

func (x *GetNotificationsResp) Reset() {
	*x = GetNotificationsResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cloudmind_core_api_system_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetNotificationsResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetNotificationsResp) ProtoMessage() {}

func (x *GetNotificationsResp) ProtoReflect() protoreflect.Message {
	mi := &file_cloudmind_core_api_system_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetNotificationsResp.ProtoReflect.Descriptor instead.
func (*GetNotificationsResp) Descriptor() ([]byte, []int) {
	return file_cloudmind_core_api_system_proto_rawDescGZIP(), []int{1}
}

func (x *GetNotificationsResp) GetNotifications() []*Notification {
	if x != nil {
		return x.Notifications
	}
	return nil
}

func (x *GetNotificationsResp) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

type GetNotificationCountReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OnlyType *int64 `protobuf:"varint,1,opt,name=onlyType,proto3,oneof" json:"onlyType" form:"onlyType" query:"onlyType"`
}

func (x *GetNotificationCountReq) Reset() {
	*x = GetNotificationCountReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cloudmind_core_api_system_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetNotificationCountReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetNotificationCountReq) ProtoMessage() {}

func (x *GetNotificationCountReq) ProtoReflect() protoreflect.Message {
	mi := &file_cloudmind_core_api_system_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetNotificationCountReq.ProtoReflect.Descriptor instead.
func (*GetNotificationCountReq) Descriptor() ([]byte, []int) {
	return file_cloudmind_core_api_system_proto_rawDescGZIP(), []int{2}
}

func (x *GetNotificationCountReq) GetOnlyType() int64 {
	if x != nil && x.OnlyType != nil {
		return *x.OnlyType
	}
	return 0
}

type GetNotificationCountResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Total int64 `protobuf:"varint,1,opt,name=total,proto3" json:"total" form:"total" query:"total"`
}

func (x *GetNotificationCountResp) Reset() {
	*x = GetNotificationCountResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cloudmind_core_api_system_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetNotificationCountResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetNotificationCountResp) ProtoMessage() {}

func (x *GetNotificationCountResp) ProtoReflect() protoreflect.Message {
	mi := &file_cloudmind_core_api_system_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetNotificationCountResp.ProtoReflect.Descriptor instead.
func (*GetNotificationCountResp) Descriptor() ([]byte, []int) {
	return file_cloudmind_core_api_system_proto_rawDescGZIP(), []int{3}
}

func (x *GetNotificationCountResp) GetTotal() int64 {
	if x != nil {
		return x.Total
	}
	return 0
}

type CreateSliderReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ImageUrl string `protobuf:"bytes,1,opt,name=imageUrl,proto3" json:"imageUrl" form:"imageUrl" query:"imageUrl"`  // 图片链接
	LinkUrl  string `protobuf:"bytes,2,opt,name=linkUrl,proto3" json:"linkUrl" form:"linkUrl" query:"linkUrl"`      // 跳转链接
	Type     int64  `protobuf:"varint,3,opt,name=type,proto3" json:"type" form:"type" query:"type"`                 // 类型
	IsPublic int64  `protobuf:"varint,4,opt,name=isPublic,proto3" json:"isPublic" form:"isPublic" query:"isPublic"` // 是否公开
}

func (x *CreateSliderReq) Reset() {
	*x = CreateSliderReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cloudmind_core_api_system_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateSliderReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateSliderReq) ProtoMessage() {}

func (x *CreateSliderReq) ProtoReflect() protoreflect.Message {
	mi := &file_cloudmind_core_api_system_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateSliderReq.ProtoReflect.Descriptor instead.
func (*CreateSliderReq) Descriptor() ([]byte, []int) {
	return file_cloudmind_core_api_system_proto_rawDescGZIP(), []int{4}
}

func (x *CreateSliderReq) GetImageUrl() string {
	if x != nil {
		return x.ImageUrl
	}
	return ""
}

func (x *CreateSliderReq) GetLinkUrl() string {
	if x != nil {
		return x.LinkUrl
	}
	return ""
}

func (x *CreateSliderReq) GetType() int64 {
	if x != nil {
		return x.Type
	}
	return 0
}

func (x *CreateSliderReq) GetIsPublic() int64 {
	if x != nil {
		return x.IsPublic
	}
	return 0
}

type CreateSliderResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *CreateSliderResp) Reset() {
	*x = CreateSliderResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cloudmind_core_api_system_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateSliderResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateSliderResp) ProtoMessage() {}

func (x *CreateSliderResp) ProtoReflect() protoreflect.Message {
	mi := &file_cloudmind_core_api_system_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateSliderResp.ProtoReflect.Descriptor instead.
func (*CreateSliderResp) Descriptor() ([]byte, []int) {
	return file_cloudmind_core_api_system_proto_rawDescGZIP(), []int{5}
}

type DeleteSliderReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SliderId string `protobuf:"bytes,1,opt,name=sliderId,proto3" json:"sliderId" form:"sliderId" query:"sliderId"`
}

func (x *DeleteSliderReq) Reset() {
	*x = DeleteSliderReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cloudmind_core_api_system_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteSliderReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteSliderReq) ProtoMessage() {}

func (x *DeleteSliderReq) ProtoReflect() protoreflect.Message {
	mi := &file_cloudmind_core_api_system_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteSliderReq.ProtoReflect.Descriptor instead.
func (*DeleteSliderReq) Descriptor() ([]byte, []int) {
	return file_cloudmind_core_api_system_proto_rawDescGZIP(), []int{6}
}

func (x *DeleteSliderReq) GetSliderId() string {
	if x != nil {
		return x.SliderId
	}
	return ""
}

type DeleteSliderResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DeleteSliderResp) Reset() {
	*x = DeleteSliderResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cloudmind_core_api_system_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteSliderResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteSliderResp) ProtoMessage() {}

func (x *DeleteSliderResp) ProtoReflect() protoreflect.Message {
	mi := &file_cloudmind_core_api_system_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteSliderResp.ProtoReflect.Descriptor instead.
func (*DeleteSliderResp) Descriptor() ([]byte, []int) {
	return file_cloudmind_core_api_system_proto_rawDescGZIP(), []int{7}
}

type UpdateSliderReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SliderId string `protobuf:"bytes,1,opt,name=sliderId,proto3" json:"sliderId" form:"sliderId" query:"sliderId"`
	ImageUrl string `protobuf:"bytes,2,opt,name=imageUrl,proto3" json:"imageUrl" form:"imageUrl" query:"imageUrl"`
	LinkUrl  string `protobuf:"bytes,3,opt,name=linkUrl,proto3" json:"linkUrl" form:"linkUrl" query:"linkUrl"`
	Type     int64  `protobuf:"varint,4,opt,name=type,proto3" json:"type" form:"type" query:"type"`
	IsPublic int64  `protobuf:"varint,5,opt,name=isPublic,proto3" json:"isPublic" form:"isPublic" query:"isPublic"`
}

func (x *UpdateSliderReq) Reset() {
	*x = UpdateSliderReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cloudmind_core_api_system_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateSliderReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateSliderReq) ProtoMessage() {}

func (x *UpdateSliderReq) ProtoReflect() protoreflect.Message {
	mi := &file_cloudmind_core_api_system_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateSliderReq.ProtoReflect.Descriptor instead.
func (*UpdateSliderReq) Descriptor() ([]byte, []int) {
	return file_cloudmind_core_api_system_proto_rawDescGZIP(), []int{8}
}

func (x *UpdateSliderReq) GetSliderId() string {
	if x != nil {
		return x.SliderId
	}
	return ""
}

func (x *UpdateSliderReq) GetImageUrl() string {
	if x != nil {
		return x.ImageUrl
	}
	return ""
}

func (x *UpdateSliderReq) GetLinkUrl() string {
	if x != nil {
		return x.LinkUrl
	}
	return ""
}

func (x *UpdateSliderReq) GetType() int64 {
	if x != nil {
		return x.Type
	}
	return 0
}

func (x *UpdateSliderReq) GetIsPublic() int64 {
	if x != nil {
		return x.IsPublic
	}
	return 0
}

type UpdateSliderResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *UpdateSliderResp) Reset() {
	*x = UpdateSliderResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cloudmind_core_api_system_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateSliderResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateSliderResp) ProtoMessage() {}

func (x *UpdateSliderResp) ProtoReflect() protoreflect.Message {
	mi := &file_cloudmind_core_api_system_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateSliderResp.ProtoReflect.Descriptor instead.
func (*UpdateSliderResp) Descriptor() ([]byte, []int) {
	return file_cloudmind_core_api_system_proto_rawDescGZIP(), []int{9}
}

type GetSlidersReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Limit     *int64  `protobuf:"varint,1,opt,name=limit,proto3,oneof" json:"limit" form:"limit" query:"limit"`
	LastToken *string `protobuf:"bytes,2,opt,name=lastToken,proto3,oneof" json:"lastToken" form:"lastToken" query:"lastToken"`
	Backward  *bool   `protobuf:"varint,3,opt,name=backward,proto3,oneof" json:"backward" form:"backward" query:"backward"`
	Offset    *int64  `protobuf:"varint,4,opt,name=offset,proto3,oneof" json:"offset" form:"offset" query:"offset"`
}

func (x *GetSlidersReq) Reset() {
	*x = GetSlidersReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cloudmind_core_api_system_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetSlidersReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetSlidersReq) ProtoMessage() {}

func (x *GetSlidersReq) ProtoReflect() protoreflect.Message {
	mi := &file_cloudmind_core_api_system_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetSlidersReq.ProtoReflect.Descriptor instead.
func (*GetSlidersReq) Descriptor() ([]byte, []int) {
	return file_cloudmind_core_api_system_proto_rawDescGZIP(), []int{10}
}

func (x *GetSlidersReq) GetLimit() int64 {
	if x != nil && x.Limit != nil {
		return *x.Limit
	}
	return 0
}

func (x *GetSlidersReq) GetLastToken() string {
	if x != nil && x.LastToken != nil {
		return *x.LastToken
	}
	return ""
}

func (x *GetSlidersReq) GetBackward() bool {
	if x != nil && x.Backward != nil {
		return *x.Backward
	}
	return false
}

func (x *GetSlidersReq) GetOffset() int64 {
	if x != nil && x.Offset != nil {
		return *x.Offset
	}
	return 0
}

type GetSlidersResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sliders []*Slider `protobuf:"bytes,1,rep,name=sliders,proto3" json:"sliders" form:"sliders" query:"sliders"`
	Total   int64     `protobuf:"varint,2,opt,name=total,proto3" json:"total" form:"total" query:"total"`
	Token   string    `protobuf:"bytes,3,opt,name=token,proto3" json:"token" form:"token" query:"token"`
}

func (x *GetSlidersResp) Reset() {
	*x = GetSlidersResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cloudmind_core_api_system_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetSlidersResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetSlidersResp) ProtoMessage() {}

func (x *GetSlidersResp) ProtoReflect() protoreflect.Message {
	mi := &file_cloudmind_core_api_system_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetSlidersResp.ProtoReflect.Descriptor instead.
func (*GetSlidersResp) Descriptor() ([]byte, []int) {
	return file_cloudmind_core_api_system_proto_rawDescGZIP(), []int{11}
}

func (x *GetSlidersResp) GetSliders() []*Slider {
	if x != nil {
		return x.Sliders
	}
	return nil
}

func (x *GetSlidersResp) GetTotal() int64 {
	if x != nil {
		return x.Total
	}
	return 0
}

func (x *GetSlidersResp) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

var File_cloudmind_core_api_system_proto protoreflect.FileDescriptor

var file_cloudmind_core_api_system_proto_rawDesc = []byte{
	0x0a, 0x1f, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x6d, 0x69, 0x6e, 0x64, 0x2f, 0x63, 0x6f, 0x72, 0x65,
	0x5f, 0x61, 0x70, 0x69, 0x2f, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x12, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x6d, 0x69, 0x6e, 0x64, 0x2e, 0x63, 0x6f, 0x72,
	0x65, 0x5f, 0x61, 0x70, 0x69, 0x1a, 0x1f, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x6d, 0x69, 0x6e, 0x64,
	0x2f, 0x63, 0x6f, 0x72, 0x65, 0x5f, 0x61, 0x70, 0x69, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x16, 0x62, 0x61, 0x73, 0x69, 0x63, 0x2f, 0x70, 0x61,
	0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0f,
	0x68, 0x74, 0x74, 0x70, 0x2f, 0x68, 0x74, 0x74, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0xef, 0x01, 0x0a, 0x13, 0x47, 0x65, 0x74, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x71, 0x12, 0x1f, 0x0a, 0x08, 0x6f, 0x6e, 0x6c, 0x79, 0x54,
	0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x48, 0x00, 0x52, 0x08, 0x6f, 0x6e, 0x6c,
	0x79, 0x54, 0x79, 0x70, 0x65, 0x88, 0x01, 0x01, 0x12, 0x19, 0x0a, 0x05, 0x6c, 0x69, 0x6d, 0x69,
	0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x48, 0x01, 0x52, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74,
	0x88, 0x01, 0x01, 0x12, 0x21, 0x0a, 0x09, 0x6c, 0x61, 0x73, 0x74, 0x54, 0x6f, 0x6b, 0x65, 0x6e,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x48, 0x02, 0x52, 0x09, 0x6c, 0x61, 0x73, 0x74, 0x54, 0x6f,
	0x6b, 0x65, 0x6e, 0x88, 0x01, 0x01, 0x12, 0x1f, 0x0a, 0x08, 0x62, 0x61, 0x63, 0x6b, 0x77, 0x61,
	0x72, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x48, 0x03, 0x52, 0x08, 0x62, 0x61, 0x63, 0x6b,
	0x77, 0x61, 0x72, 0x64, 0x88, 0x01, 0x01, 0x12, 0x1b, 0x0a, 0x06, 0x6f, 0x66, 0x66, 0x73, 0x65,
	0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x48, 0x04, 0x52, 0x06, 0x6f, 0x66, 0x66, 0x73, 0x65,
	0x74, 0x88, 0x01, 0x01, 0x42, 0x0b, 0x0a, 0x09, 0x5f, 0x6f, 0x6e, 0x6c, 0x79, 0x54, 0x79, 0x70,
	0x65, 0x42, 0x08, 0x0a, 0x06, 0x5f, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x42, 0x0c, 0x0a, 0x0a, 0x5f,
	0x6c, 0x61, 0x73, 0x74, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x42, 0x0b, 0x0a, 0x09, 0x5f, 0x62, 0x61,
	0x63, 0x6b, 0x77, 0x61, 0x72, 0x64, 0x42, 0x09, 0x0a, 0x07, 0x5f, 0x6f, 0x66, 0x66, 0x73, 0x65,
	0x74, 0x22, 0x74, 0x0a, 0x14, 0x47, 0x65, 0x74, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x73, 0x70, 0x12, 0x46, 0x0a, 0x0d, 0x6e, 0x6f, 0x74,
	0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x20, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x6d, 0x69, 0x6e, 0x64, 0x2e, 0x63, 0x6f, 0x72,
	0x65, 0x5f, 0x61, 0x70, 0x69, 0x2e, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x52, 0x0d, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x47, 0x0a, 0x17, 0x47, 0x65, 0x74, 0x4e, 0x6f,
	0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x52,
	0x65, 0x71, 0x12, 0x1f, 0x0a, 0x08, 0x6f, 0x6e, 0x6c, 0x79, 0x54, 0x79, 0x70, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x03, 0x48, 0x00, 0x52, 0x08, 0x6f, 0x6e, 0x6c, 0x79, 0x54, 0x79, 0x70, 0x65,
	0x88, 0x01, 0x01, 0x42, 0x0b, 0x0a, 0x09, 0x5f, 0x6f, 0x6e, 0x6c, 0x79, 0x54, 0x79, 0x70, 0x65,
	0x22, 0x30, 0x0a, 0x18, 0x47, 0x65, 0x74, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x12, 0x14, 0x0a, 0x05,
	0x74, 0x6f, 0x74, 0x61, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x74, 0x6f, 0x74,
	0x61, 0x6c, 0x22, 0x77, 0x0a, 0x0f, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x6c, 0x69, 0x64,
	0x65, 0x72, 0x52, 0x65, 0x71, 0x12, 0x1a, 0x0a, 0x08, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x55, 0x72,
	0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x55, 0x72,
	0x6c, 0x12, 0x18, 0x0a, 0x07, 0x6c, 0x69, 0x6e, 0x6b, 0x55, 0x72, 0x6c, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x6c, 0x69, 0x6e, 0x6b, 0x55, 0x72, 0x6c, 0x12, 0x12, 0x0a, 0x04, 0x74,
	0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12,
	0x1a, 0x0a, 0x08, 0x69, 0x73, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x08, 0x69, 0x73, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x22, 0x12, 0x0a, 0x10, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x6c, 0x69, 0x64, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x22,
	0x2d, 0x0a, 0x0f, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x53, 0x6c, 0x69, 0x64, 0x65, 0x72, 0x52,
	0x65, 0x71, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x6c, 0x69, 0x64, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x6c, 0x69, 0x64, 0x65, 0x72, 0x49, 0x64, 0x22, 0x12,
	0x0a, 0x10, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x53, 0x6c, 0x69, 0x64, 0x65, 0x72, 0x52, 0x65,
	0x73, 0x70, 0x22, 0x93, 0x01, 0x0a, 0x0f, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x53, 0x6c, 0x69,
	0x64, 0x65, 0x72, 0x52, 0x65, 0x71, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x6c, 0x69, 0x64, 0x65, 0x72,
	0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x6c, 0x69, 0x64, 0x65, 0x72,
	0x49, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x55, 0x72, 0x6c, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x55, 0x72, 0x6c, 0x12, 0x18,
	0x0a, 0x07, 0x6c, 0x69, 0x6e, 0x6b, 0x55, 0x72, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x6c, 0x69, 0x6e, 0x6b, 0x55, 0x72, 0x6c, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x1a, 0x0a, 0x08,
	0x69, 0x73, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08,
	0x69, 0x73, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x22, 0x12, 0x0a, 0x10, 0x55, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x53, 0x6c, 0x69, 0x64, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x22, 0xbb, 0x01, 0x0a,
	0x0d, 0x47, 0x65, 0x74, 0x53, 0x6c, 0x69, 0x64, 0x65, 0x72, 0x73, 0x52, 0x65, 0x71, 0x12, 0x19,
	0x0a, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x48, 0x00, 0x52,
	0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x88, 0x01, 0x01, 0x12, 0x21, 0x0a, 0x09, 0x6c, 0x61, 0x73,
	0x74, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x48, 0x01, 0x52, 0x09,
	0x6c, 0x61, 0x73, 0x74, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x88, 0x01, 0x01, 0x12, 0x1f, 0x0a, 0x08,
	0x62, 0x61, 0x63, 0x6b, 0x77, 0x61, 0x72, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x48, 0x02,
	0x52, 0x08, 0x62, 0x61, 0x63, 0x6b, 0x77, 0x61, 0x72, 0x64, 0x88, 0x01, 0x01, 0x12, 0x1b, 0x0a,
	0x06, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x48, 0x03, 0x52,
	0x06, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x88, 0x01, 0x01, 0x42, 0x08, 0x0a, 0x06, 0x5f, 0x6c,
	0x69, 0x6d, 0x69, 0x74, 0x42, 0x0c, 0x0a, 0x0a, 0x5f, 0x6c, 0x61, 0x73, 0x74, 0x54, 0x6f, 0x6b,
	0x65, 0x6e, 0x42, 0x0b, 0x0a, 0x09, 0x5f, 0x62, 0x61, 0x63, 0x6b, 0x77, 0x61, 0x72, 0x64, 0x42,
	0x09, 0x0a, 0x07, 0x5f, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x22, 0x72, 0x0a, 0x0e, 0x47, 0x65,
	0x74, 0x53, 0x6c, 0x69, 0x64, 0x65, 0x72, 0x73, 0x52, 0x65, 0x73, 0x70, 0x12, 0x34, 0x0a, 0x07,
	0x73, 0x6c, 0x69, 0x64, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x63, 0x6c, 0x6f, 0x75, 0x64, 0x6d, 0x69, 0x6e, 0x64, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x5f, 0x61,
	0x70, 0x69, 0x2e, 0x53, 0x6c, 0x69, 0x64, 0x65, 0x72, 0x52, 0x07, 0x73, 0x6c, 0x69, 0x64, 0x65,
	0x72, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x05, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65,
	0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x42, 0x53,
	0x5a, 0x51, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x43, 0x6c, 0x6f,
	0x75, 0x64, 0x53, 0x74, 0x72, 0x69, 0x76, 0x65, 0x72, 0x2f, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x6d,
	0x69, 0x6e, 0x64, 0x2d, 0x63, 0x6f, 0x72, 0x65, 0x2d, 0x61, 0x70, 0x69, 0x2f, 0x62, 0x69, 0x7a,
	0x2f, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x64, 0x74, 0x6f,
	0x2f, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x6d, 0x69, 0x6e, 0x64, 0x2f, 0x63, 0x6f, 0x72, 0x65, 0x5f,
	0x61, 0x70, 0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_cloudmind_core_api_system_proto_rawDescOnce sync.Once
	file_cloudmind_core_api_system_proto_rawDescData = file_cloudmind_core_api_system_proto_rawDesc
)

func file_cloudmind_core_api_system_proto_rawDescGZIP() []byte {
	file_cloudmind_core_api_system_proto_rawDescOnce.Do(func() {
		file_cloudmind_core_api_system_proto_rawDescData = protoimpl.X.CompressGZIP(file_cloudmind_core_api_system_proto_rawDescData)
	})
	return file_cloudmind_core_api_system_proto_rawDescData
}

var file_cloudmind_core_api_system_proto_msgTypes = make([]protoimpl.MessageInfo, 12)
var file_cloudmind_core_api_system_proto_goTypes = []interface{}{
	(*GetNotificationsReq)(nil),      // 0: cloudmind.core_api.GetNotificationsReq
	(*GetNotificationsResp)(nil),     // 1: cloudmind.core_api.GetNotificationsResp
	(*GetNotificationCountReq)(nil),  // 2: cloudmind.core_api.GetNotificationCountReq
	(*GetNotificationCountResp)(nil), // 3: cloudmind.core_api.GetNotificationCountResp
	(*CreateSliderReq)(nil),          // 4: cloudmind.core_api.CreateSliderReq
	(*CreateSliderResp)(nil),         // 5: cloudmind.core_api.CreateSliderResp
	(*DeleteSliderReq)(nil),          // 6: cloudmind.core_api.DeleteSliderReq
	(*DeleteSliderResp)(nil),         // 7: cloudmind.core_api.DeleteSliderResp
	(*UpdateSliderReq)(nil),          // 8: cloudmind.core_api.UpdateSliderReq
	(*UpdateSliderResp)(nil),         // 9: cloudmind.core_api.UpdateSliderResp
	(*GetSlidersReq)(nil),            // 10: cloudmind.core_api.GetSlidersReq
	(*GetSlidersResp)(nil),           // 11: cloudmind.core_api.GetSlidersResp
	(*Notification)(nil),             // 12: cloudmind.core_api.Notification
	(*Slider)(nil),                   // 13: cloudmind.core_api.Slider
}
var file_cloudmind_core_api_system_proto_depIdxs = []int32{
	12, // 0: cloudmind.core_api.GetNotificationsResp.notifications:type_name -> cloudmind.core_api.Notification
	13, // 1: cloudmind.core_api.GetSlidersResp.sliders:type_name -> cloudmind.core_api.Slider
	2,  // [2:2] is the sub-list for method output_type
	2,  // [2:2] is the sub-list for method input_type
	2,  // [2:2] is the sub-list for extension type_name
	2,  // [2:2] is the sub-list for extension extendee
	0,  // [0:2] is the sub-list for field type_name
}


func file_cloudmind_core_api_system_proto_init() {
	if File_cloudmind_core_api_system_proto != nil {
		return
	}
	file_cloudmind_core_api_common_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_cloudmind_core_api_system_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetNotificationsReq); i {
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
		file_cloudmind_core_api_system_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetNotificationsResp); i {
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
		file_cloudmind_core_api_system_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetNotificationCountReq); i {
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
		file_cloudmind_core_api_system_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetNotificationCountResp); i {
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
		file_cloudmind_core_api_system_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateSliderReq); i {
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
		file_cloudmind_core_api_system_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateSliderResp); i {
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
		file_cloudmind_core_api_system_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteSliderReq); i {
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
		file_cloudmind_core_api_system_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteSliderResp); i {
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
		file_cloudmind_core_api_system_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateSliderReq); i {
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
		file_cloudmind_core_api_system_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateSliderResp); i {
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
		file_cloudmind_core_api_system_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetSlidersReq); i {
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
		file_cloudmind_core_api_system_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetSlidersResp); i {
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
	file_cloudmind_core_api_system_proto_msgTypes[0].OneofWrappers = []interface{}{}
	file_cloudmind_core_api_system_proto_msgTypes[2].OneofWrappers = []interface{}{}
	file_cloudmind_core_api_system_proto_msgTypes[10].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_cloudmind_core_api_system_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   12,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_cloudmind_core_api_system_proto_goTypes,
		DependencyIndexes: file_cloudmind_core_api_system_proto_depIdxs,
		MessageInfos:      file_cloudmind_core_api_system_proto_msgTypes,
	}.Build()
	File_cloudmind_core_api_system_proto = out.File
	file_cloudmind_core_api_system_proto_rawDesc = nil
	file_cloudmind_core_api_system_proto_goTypes = nil
	file_cloudmind_core_api_system_proto_depIdxs = nil
}
