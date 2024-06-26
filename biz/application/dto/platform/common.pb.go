// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.19.4
// source: platform/common.proto

package platform

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

type Attrs int32

const (
	Attrs_UnknownAttrs         Attrs = 0
	Attrs_None                 Attrs = 1 // 无
	Attrs_Pinned               Attrs = 2 // 置顶
	Attrs_Highlighted          Attrs = 3 // 精华
	Attrs_PinnedAndHighlighted Attrs = 4 // 置顶+精华
)

// Enum value maps for Attrs.
var (
	Attrs_name = map[int32]string{
		0: "UnknownAttrs",
		1: "None",
		2: "Pinned",
		3: "Highlighted",
		4: "PinnedAndHighlighted",
	}
	Attrs_value = map[string]int32{
		"UnknownAttrs":         0,
		"None":                 1,
		"Pinned":               2,
		"Highlighted":          3,
		"PinnedAndHighlighted": 4,
	}
)

func (x Attrs) Enum() *Attrs {
	p := new(Attrs)
	*p = x
	return p
}

func (x Attrs) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Attrs) Descriptor() protoreflect.EnumDescriptor {
	return file_platform_common_proto_enumTypes[0].Descriptor()
}

func (Attrs) Type() protoreflect.EnumType {
	return &file_platform_common_proto_enumTypes[0]
}

func (x Attrs) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Attrs.Descriptor instead.
func (Attrs) EnumDescriptor() ([]byte, []int) {
	return file_platform_common_proto_rawDescGZIP(), []int{0}
}

type State int32

const (
	State_UnknownState State = 0
	State_Normal       State = 1 // 正常状态
	State_Hidden       State = 2 // 隐藏状态
)

// Enum value maps for State.
var (
	State_name = map[int32]string{
		0: "UnknownState",
		1: "Normal",
		2: "Hidden",
	}
	State_value = map[string]int32{
		"UnknownState": 0,
		"Normal":       1,
		"Hidden":       2,
	}
)

func (x State) Enum() *State {
	p := new(State)
	*p = x
	return p
}

func (x State) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (State) Descriptor() protoreflect.EnumDescriptor {
	return file_platform_common_proto_enumTypes[1].Descriptor()
}

func (State) Type() protoreflect.EnumType {
	return &file_platform_common_proto_enumTypes[1]
}

func (x State) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use State.Descriptor instead.
func (State) EnumDescriptor() ([]byte, []int) {
	return file_platform_common_proto_rawDescGZIP(), []int{1}
}

type Comment struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CommentId  string   `protobuf:"bytes,1,opt,name=commentId,proto3" json:"commentId" form:"commentId" query:"commentId"`
	SubjectId  string   `protobuf:"bytes,2,opt,name=subjectId,proto3" json:"subjectId" form:"subjectId" query:"subjectId"`
	RootId     string   `protobuf:"bytes,3,opt,name=rootId,proto3" json:"rootId" form:"rootId" query:"rootId"`
	FatherId   string   `protobuf:"bytes,4,opt,name=fatherId,proto3" json:"fatherId" form:"fatherId" query:"fatherId"`
	Count      int64    `protobuf:"varint,5,opt,name=count,proto3" json:"count" form:"count" query:"count"`             // 回复数
	State      int64    `protobuf:"varint,6,opt,name=state,proto3" json:"state" form:"state" query:"state"`             // 1: 正常, 2: 删除
	Attrs      int64    `protobuf:"varint,7,opt,name=attrs,proto3" json:"attrs" form:"attrs" query:"attrs"`             // 1: 无, 2: 置顶, 3: 精华, 4: 置顶+精华
	Labels     []string `protobuf:"bytes,8,rep,name=labels,proto3" json:"labels" form:"labels" query:"labels"`          // 标签：作者点赞，作者回复等
	UserId     string   `protobuf:"bytes,9,opt,name=userId,proto3" json:"userId" form:"userId" query:"userId"`          // 评论者
	AtUserId   string   `protobuf:"bytes,10,opt,name=atUserId,proto3" json:"atUserId" form:"atUserId" query:"atUserId"` // @谁
	Content    string   `protobuf:"bytes,11,opt,name=content,proto3" json:"content" form:"content" query:"content"`     // 内容
	Meta       string   `protobuf:"bytes,12,opt,name=meta,proto3" json:"meta" form:"meta" query:"meta"`                 // 皮肤，字体等
	CreateTime int64    `protobuf:"varint,13,opt,name=createTime,proto3" json:"createTime" form:"createTime" query:"createTime"`
	Type       int64    `protobuf:"varint,14,opt,name=type,proto3" json:"type" form:"type" query:"type"`
}

func (x *Comment) Reset() {
	*x = Comment{}
	if protoimpl.UnsafeEnabled {
		mi := &file_platform_common_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Comment) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Comment) ProtoMessage() {}

func (x *Comment) ProtoReflect() protoreflect.Message {
	mi := &file_platform_common_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Comment.ProtoReflect.Descriptor instead.
func (*Comment) Descriptor() ([]byte, []int) {
	return file_platform_common_proto_rawDescGZIP(), []int{0}
}

func (x *Comment) GetCommentId() string {
	if x != nil {
		return x.CommentId
	}
	return ""
}

func (x *Comment) GetSubjectId() string {
	if x != nil {
		return x.SubjectId
	}
	return ""
}

func (x *Comment) GetRootId() string {
	if x != nil {
		return x.RootId
	}
	return ""
}

func (x *Comment) GetFatherId() string {
	if x != nil {
		return x.FatherId
	}
	return ""
}

func (x *Comment) GetCount() int64 {
	if x != nil {
		return x.Count
	}
	return 0
}

func (x *Comment) GetState() int64 {
	if x != nil {
		return x.State
	}
	return 0
}

func (x *Comment) GetAttrs() int64 {
	if x != nil {
		return x.Attrs
	}
	return 0
}

func (x *Comment) GetLabels() []string {
	if x != nil {
		return x.Labels
	}
	return nil
}

func (x *Comment) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *Comment) GetAtUserId() string {
	if x != nil {
		return x.AtUserId
	}
	return ""
}

func (x *Comment) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

func (x *Comment) GetMeta() string {
	if x != nil {
		return x.Meta
	}
	return ""
}

func (x *Comment) GetCreateTime() int64 {
	if x != nil {
		return x.CreateTime
	}
	return 0
}

func (x *Comment) GetType() int64 {
	if x != nil {
		return x.Type
	}
	return 0
}

type CommentBlock struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RootComment *Comment   `protobuf:"bytes,1,opt,name=rootComment,proto3" json:"rootComment" form:"rootComment" query:"rootComment"`
	ReplyList   *ReplyList `protobuf:"bytes,2,opt,name=replyList,proto3" json:"replyList" form:"replyList" query:"replyList"`
}

func (x *CommentBlock) Reset() {
	*x = CommentBlock{}
	if protoimpl.UnsafeEnabled {
		mi := &file_platform_common_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CommentBlock) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CommentBlock) ProtoMessage() {}

func (x *CommentBlock) ProtoReflect() protoreflect.Message {
	mi := &file_platform_common_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CommentBlock.ProtoReflect.Descriptor instead.
func (*CommentBlock) Descriptor() ([]byte, []int) {
	return file_platform_common_proto_rawDescGZIP(), []int{1}
}

func (x *CommentBlock) GetRootComment() *Comment {
	if x != nil {
		return x.RootComment
	}
	return nil
}

func (x *CommentBlock) GetReplyList() *ReplyList {
	if x != nil {
		return x.ReplyList
	}
	return nil
}

type ReplyList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Comments []*Comment `protobuf:"bytes,1,rep,name=comments,proto3" json:"comments" form:"comments" query:"comments"`
	Total    int64      `protobuf:"varint,2,opt,name=total,proto3" json:"total" form:"total" query:"total"`
	Token    string     `protobuf:"bytes,3,opt,name=token,proto3" json:"token" form:"token" query:"token"`
}

func (x *ReplyList) Reset() {
	*x = ReplyList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_platform_common_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReplyList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReplyList) ProtoMessage() {}

func (x *ReplyList) ProtoReflect() protoreflect.Message {
	mi := &file_platform_common_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReplyList.ProtoReflect.Descriptor instead.
func (*ReplyList) Descriptor() ([]byte, []int) {
	return file_platform_common_proto_rawDescGZIP(), []int{2}
}

func (x *ReplyList) GetComments() []*Comment {
	if x != nil {
		return x.Comments
	}
	return nil
}

func (x *ReplyList) GetTotal() int64 {
	if x != nil {
		return x.Total
	}
	return 0
}

func (x *ReplyList) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

type CommentFilterOptions struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OnlyUserId     *string  `protobuf:"bytes,1,opt,name=onlyUserId,proto3,oneof" json:"onlyUserId" form:"onlyUserId" query:"onlyUserId"`
	OnlyAtUserId   *string  `protobuf:"bytes,2,opt,name=onlyAtUserId,proto3,oneof" json:"onlyAtUserId" form:"onlyAtUserId" query:"onlyAtUserId"`
	OnlyState      *int64   `protobuf:"varint,3,opt,name=onlyState,proto3,oneof" json:"onlyState" form:"onlyState" query:"onlyState"`
	OnlyAttrs      *int64   `protobuf:"varint,4,opt,name=onlyAttrs,proto3,oneof" json:"onlyAttrs" form:"onlyAttrs" query:"onlyAttrs"`
	OnlyCommentIds []string `protobuf:"bytes,5,rep,name=onlyCommentIds,proto3" json:"onlyCommentIds" form:"onlyCommentIds" query:"onlyCommentIds"`
}

func (x *CommentFilterOptions) Reset() {
	*x = CommentFilterOptions{}
	if protoimpl.UnsafeEnabled {
		mi := &file_platform_common_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CommentFilterOptions) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CommentFilterOptions) ProtoMessage() {}

func (x *CommentFilterOptions) ProtoReflect() protoreflect.Message {
	mi := &file_platform_common_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CommentFilterOptions.ProtoReflect.Descriptor instead.
func (*CommentFilterOptions) Descriptor() ([]byte, []int) {
	return file_platform_common_proto_rawDescGZIP(), []int{3}
}

func (x *CommentFilterOptions) GetOnlyUserId() string {
	if x != nil && x.OnlyUserId != nil {
		return *x.OnlyUserId
	}
	return ""
}

func (x *CommentFilterOptions) GetOnlyAtUserId() string {
	if x != nil && x.OnlyAtUserId != nil {
		return *x.OnlyAtUserId
	}
	return ""
}

func (x *CommentFilterOptions) GetOnlyState() int64 {
	if x != nil && x.OnlyState != nil {
		return *x.OnlyState
	}
	return 0
}

func (x *CommentFilterOptions) GetOnlyAttrs() int64 {
	if x != nil && x.OnlyAttrs != nil {
		return *x.OnlyAttrs
	}
	return 0
}

func (x *CommentFilterOptions) GetOnlyCommentIds() []string {
	if x != nil {
		return x.OnlyCommentIds
	}
	return nil
}

type LabelFilterOptions struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *LabelFilterOptions) Reset() {
	*x = LabelFilterOptions{}
	if protoimpl.UnsafeEnabled {
		mi := &file_platform_common_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LabelFilterOptions) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LabelFilterOptions) ProtoMessage() {}

func (x *LabelFilterOptions) ProtoReflect() protoreflect.Message {
	mi := &file_platform_common_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LabelFilterOptions.ProtoReflect.Descriptor instead.
func (*LabelFilterOptions) Descriptor() ([]byte, []int) {
	return file_platform_common_proto_rawDescGZIP(), []int{4}
}

type Subject struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SubjectId    string `protobuf:"bytes,1,opt,name=subjectId,proto3" json:"subjectId" form:"subjectId" query:"subjectId"`
	UserId       string `protobuf:"bytes,2,opt,name=userId,proto3" json:"userId" form:"userId" query:"userId"`
	TopCommentId string `protobuf:"bytes,3,opt,name=TopCommentId,proto3" json:"TopCommentId" form:"TopCommentId" query:"TopCommentId"`
	RootCount    int64  `protobuf:"varint,4,opt,name=rootCount,proto3" json:"rootCount" form:"rootCount" query:"rootCount"`
	AllCount     int64  `protobuf:"varint,5,opt,name=allCount,proto3" json:"allCount" form:"allCount" query:"allCount"`
	State        int64  `protobuf:"varint,6,opt,name=state,proto3" json:"state" form:"state" query:"state"` // 1: 正常, 2: 删除
	Attrs        int64  `protobuf:"varint,7,opt,name=attrs,proto3" json:"attrs" form:"attrs" query:"attrs"` // 1: 无, 2: 置顶, 3: 精华, 4: 置顶+精华
	Type         int64  `protobuf:"varint,8,opt,name=type,proto3" json:"type" form:"type" query:"type"`
}

func (x *Subject) Reset() {
	*x = Subject{}
	if protoimpl.UnsafeEnabled {
		mi := &file_platform_common_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Subject) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Subject) ProtoMessage() {}

func (x *Subject) ProtoReflect() protoreflect.Message {
	mi := &file_platform_common_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Subject.ProtoReflect.Descriptor instead.
func (*Subject) Descriptor() ([]byte, []int) {
	return file_platform_common_proto_rawDescGZIP(), []int{5}
}

func (x *Subject) GetSubjectId() string {
	if x != nil {
		return x.SubjectId
	}
	return ""
}

func (x *Subject) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *Subject) GetTopCommentId() string {
	if x != nil {
		return x.TopCommentId
	}
	return ""
}

func (x *Subject) GetRootCount() int64 {
	if x != nil {
		return x.RootCount
	}
	return 0
}

func (x *Subject) GetAllCount() int64 {
	if x != nil {
		return x.AllCount
	}
	return 0
}

func (x *Subject) GetState() int64 {
	if x != nil {
		return x.State
	}
	return 0
}

func (x *Subject) GetAttrs() int64 {
	if x != nil {
		return x.Attrs
	}
	return 0
}

func (x *Subject) GetType() int64 {
	if x != nil {
		return x.Type
	}
	return 0
}

type Label struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LabelId  string `protobuf:"bytes,1,opt,name=labelId,proto3" json:"labelId" form:"labelId" query:"labelId"`
	FatherId string `protobuf:"bytes,2,opt,name=fatherId,proto3" json:"fatherId" form:"fatherId" query:"fatherId"`
	Value    string `protobuf:"bytes,3,opt,name=value,proto3" json:"value" form:"value" query:"value"`
}

func (x *Label) Reset() {
	*x = Label{}
	if protoimpl.UnsafeEnabled {
		mi := &file_platform_common_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Label) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Label) ProtoMessage() {}

func (x *Label) ProtoReflect() protoreflect.Message {
	mi := &file_platform_common_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Label.ProtoReflect.Descriptor instead.
func (*Label) Descriptor() ([]byte, []int) {
	return file_platform_common_proto_rawDescGZIP(), []int{6}
}

func (x *Label) GetLabelId() string {
	if x != nil {
		return x.LabelId
	}
	return ""
}

func (x *Label) GetFatherId() string {
	if x != nil {
		return x.FatherId
	}
	return ""
}

func (x *Label) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

type Relation struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FromType     int64  `protobuf:"varint,1,opt,name=fromType,proto3" json:"fromType" form:"fromType" query:"fromType"`
	FromId       string `protobuf:"bytes,2,opt,name=fromId,proto3" json:"fromId" form:"fromId" query:"fromId"`
	ToType       int64  `protobuf:"varint,3,opt,name=toType,proto3" json:"toType" form:"toType" query:"toType"`
	ToId         string `protobuf:"bytes,4,opt,name=toId,proto3" json:"toId" form:"toId" query:"toId"`
	RelationType int64  `protobuf:"varint,5,opt,name=relationType,proto3" json:"relationType" form:"relationType" query:"relationType"`
	CreateTime   int64  `protobuf:"varint,6,opt,name=createTime,proto3" json:"createTime" form:"createTime" query:"createTime"`
}

func (x *Relation) Reset() {
	*x = Relation{}
	if protoimpl.UnsafeEnabled {
		mi := &file_platform_common_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Relation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Relation) ProtoMessage() {}

func (x *Relation) ProtoReflect() protoreflect.Message {
	mi := &file_platform_common_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Relation.ProtoReflect.Descriptor instead.
func (*Relation) Descriptor() ([]byte, []int) {
	return file_platform_common_proto_rawDescGZIP(), []int{7}
}

func (x *Relation) GetFromType() int64 {
	if x != nil {
		return x.FromType
	}
	return 0
}

func (x *Relation) GetFromId() string {
	if x != nil {
		return x.FromId
	}
	return ""
}

func (x *Relation) GetToType() int64 {
	if x != nil {
		return x.ToType
	}
	return 0
}

func (x *Relation) GetToId() string {
	if x != nil {
		return x.ToId
	}
	return ""
}

func (x *Relation) GetRelationType() int64 {
	if x != nil {
		return x.RelationType
	}
	return 0
}

func (x *Relation) GetCreateTime() int64 {
	if x != nil {
		return x.CreateTime
	}
	return 0
}

var File_platform_common_proto protoreflect.FileDescriptor

var file_platform_common_proto_rawDesc = []byte{
	0x0a, 0x15, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72,
	0x6d, 0x22, 0xe9, 0x02, 0x0a, 0x07, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x1c, 0x0a,
	0x09, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x09, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x73,
	0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09,
	0x73, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x6f, 0x6f,
	0x74, 0x49, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x6f, 0x6f, 0x74, 0x49,
	0x64, 0x12, 0x1a, 0x0a, 0x08, 0x66, 0x61, 0x74, 0x68, 0x65, 0x72, 0x49, 0x64, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x61, 0x74, 0x68, 0x65, 0x72, 0x49, 0x64, 0x12, 0x14, 0x0a,
	0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x63, 0x6f,
	0x75, 0x6e, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x61, 0x74, 0x74,
	0x72, 0x73, 0x18, 0x07, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x61, 0x74, 0x74, 0x72, 0x73, 0x12,
	0x16, 0x0a, 0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x18, 0x08, 0x20, 0x03, 0x28, 0x09, 0x52,
	0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49,
	0x64, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12,
	0x1a, 0x0a, 0x08, 0x61, 0x74, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x0a, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x61, 0x74, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x63,
	0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f,
	0x6e, 0x74, 0x65, 0x6e, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6d, 0x65, 0x74, 0x61, 0x18, 0x0c, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x6d, 0x65, 0x74, 0x61, 0x12, 0x1e, 0x0a, 0x0a, 0x63, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x63,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70,
	0x65, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x22, 0x76, 0x0a,
	0x0c, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x12, 0x33, 0x0a,
	0x0b, 0x72, 0x6f, 0x6f, 0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x11, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x43, 0x6f,
	0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x0b, 0x72, 0x6f, 0x6f, 0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x65,
	0x6e, 0x74, 0x12, 0x31, 0x0a, 0x09, 0x72, 0x65, 0x70, 0x6c, 0x79, 0x4c, 0x69, 0x73, 0x74, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d,
	0x2e, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x09, 0x72, 0x65, 0x70, 0x6c,
	0x79, 0x4c, 0x69, 0x73, 0x74, 0x22, 0x66, 0x0a, 0x09, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x4c, 0x69,
	0x73, 0x74, 0x12, 0x2d, 0x0a, 0x08, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e,
	0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x08, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74,
	0x73, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x05, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x8e, 0x02,
	0x0a, 0x14, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x4f,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x23, 0x0a, 0x0a, 0x6f, 0x6e, 0x6c, 0x79, 0x55, 0x73,
	0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x0a, 0x6f, 0x6e,
	0x6c, 0x79, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x88, 0x01, 0x01, 0x12, 0x27, 0x0a, 0x0c, 0x6f,
	0x6e, 0x6c, 0x79, 0x41, 0x74, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x48, 0x01, 0x52, 0x0c, 0x6f, 0x6e, 0x6c, 0x79, 0x41, 0x74, 0x55, 0x73, 0x65, 0x72, 0x49,
	0x64, 0x88, 0x01, 0x01, 0x12, 0x21, 0x0a, 0x09, 0x6f, 0x6e, 0x6c, 0x79, 0x53, 0x74, 0x61, 0x74,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x48, 0x02, 0x52, 0x09, 0x6f, 0x6e, 0x6c, 0x79, 0x53,
	0x74, 0x61, 0x74, 0x65, 0x88, 0x01, 0x01, 0x12, 0x21, 0x0a, 0x09, 0x6f, 0x6e, 0x6c, 0x79, 0x41,
	0x74, 0x74, 0x72, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x48, 0x03, 0x52, 0x09, 0x6f, 0x6e,
	0x6c, 0x79, 0x41, 0x74, 0x74, 0x72, 0x73, 0x88, 0x01, 0x01, 0x12, 0x26, 0x0a, 0x0e, 0x6f, 0x6e,
	0x6c, 0x79, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x73, 0x18, 0x05, 0x20, 0x03,
	0x28, 0x09, 0x52, 0x0e, 0x6f, 0x6e, 0x6c, 0x79, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x49,
	0x64, 0x73, 0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x6f, 0x6e, 0x6c, 0x79, 0x55, 0x73, 0x65, 0x72, 0x49,
	0x64, 0x42, 0x0f, 0x0a, 0x0d, 0x5f, 0x6f, 0x6e, 0x6c, 0x79, 0x41, 0x74, 0x55, 0x73, 0x65, 0x72,
	0x49, 0x64, 0x42, 0x0c, 0x0a, 0x0a, 0x5f, 0x6f, 0x6e, 0x6c, 0x79, 0x53, 0x74, 0x61, 0x74, 0x65,
	0x42, 0x0c, 0x0a, 0x0a, 0x5f, 0x6f, 0x6e, 0x6c, 0x79, 0x41, 0x74, 0x74, 0x72, 0x73, 0x22, 0x14,
	0x0a, 0x12, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x4f, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x22, 0xdd, 0x01, 0x0a, 0x07, 0x53, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74,
	0x12, 0x1c, 0x0a, 0x09, 0x73, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x49, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x49, 0x64, 0x12, 0x16,
	0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x22, 0x0a, 0x0c, 0x54, 0x6f, 0x70, 0x43, 0x6f, 0x6d,
	0x6d, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x54, 0x6f,
	0x70, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x72, 0x6f,
	0x6f, 0x74, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x72,
	0x6f, 0x6f, 0x74, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x61, 0x6c, 0x6c, 0x43,
	0x6f, 0x75, 0x6e, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x61, 0x6c, 0x6c, 0x43,
	0x6f, 0x75, 0x6e, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x61, 0x74,
	0x74, 0x72, 0x73, 0x18, 0x07, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x61, 0x74, 0x74, 0x72, 0x73,
	0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04,
	0x74, 0x79, 0x70, 0x65, 0x22, 0x53, 0x0a, 0x05, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x12, 0x18, 0x0a,
	0x07, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x6c, 0x61, 0x62, 0x65, 0x6c, 0x49, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x66, 0x61, 0x74, 0x68, 0x65,
	0x72, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x61, 0x74, 0x68, 0x65,
	0x72, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0xae, 0x01, 0x0a, 0x08, 0x52, 0x65,
	0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1a, 0x0a, 0x08, 0x66, 0x72, 0x6f, 0x6d, 0x54, 0x79,
	0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x66, 0x72, 0x6f, 0x6d, 0x54, 0x79,
	0x70, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x66, 0x72, 0x6f, 0x6d, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x66, 0x72, 0x6f, 0x6d, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x74, 0x6f,
	0x54, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x74, 0x6f, 0x54, 0x79,
	0x70, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x6f, 0x49, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x74, 0x6f, 0x49, 0x64, 0x12, 0x22, 0x0a, 0x0c, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0c, 0x72, 0x65,
	0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x63, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a,
	0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x2a, 0x5a, 0x0a, 0x05, 0x41, 0x74,
	0x74, 0x72, 0x73, 0x12, 0x10, 0x0a, 0x0c, 0x55, 0x6e, 0x6b, 0x6e, 0x6f, 0x77, 0x6e, 0x41, 0x74,
	0x74, 0x72, 0x73, 0x10, 0x00, 0x12, 0x08, 0x0a, 0x04, 0x4e, 0x6f, 0x6e, 0x65, 0x10, 0x01, 0x12,
	0x0a, 0x0a, 0x06, 0x50, 0x69, 0x6e, 0x6e, 0x65, 0x64, 0x10, 0x02, 0x12, 0x0f, 0x0a, 0x0b, 0x48,
	0x69, 0x67, 0x68, 0x6c, 0x69, 0x67, 0x68, 0x74, 0x65, 0x64, 0x10, 0x03, 0x12, 0x18, 0x0a, 0x14,
	0x50, 0x69, 0x6e, 0x6e, 0x65, 0x64, 0x41, 0x6e, 0x64, 0x48, 0x69, 0x67, 0x68, 0x6c, 0x69, 0x67,
	0x68, 0x74, 0x65, 0x64, 0x10, 0x04, 0x2a, 0x31, 0x0a, 0x05, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12,
	0x10, 0x0a, 0x0c, 0x55, 0x6e, 0x6b, 0x6e, 0x6f, 0x77, 0x6e, 0x53, 0x74, 0x61, 0x74, 0x65, 0x10,
	0x00, 0x12, 0x0a, 0x0a, 0x06, 0x4e, 0x6f, 0x72, 0x6d, 0x61, 0x6c, 0x10, 0x01, 0x12, 0x0a, 0x0a,
	0x06, 0x48, 0x69, 0x64, 0x64, 0x65, 0x6e, 0x10, 0x02, 0x42, 0x49, 0x5a, 0x47, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x43, 0x6c, 0x6f, 0x75, 0x64, 0x53, 0x74, 0x72,
	0x69, 0x76, 0x65, 0x72, 0x2f, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x6d, 0x69, 0x6e, 0x64, 0x2d, 0x63,
	0x6f, 0x72, 0x65, 0x2d, 0x61, 0x70, 0x69, 0x2f, 0x62, 0x69, 0x7a, 0x2f, 0x61, 0x70, 0x70, 0x6c,
	0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x64, 0x74, 0x6f, 0x2f, 0x70, 0x6c, 0x61, 0x74,
	0x66, 0x6f, 0x72, 0x6d, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_platform_common_proto_rawDescOnce sync.Once
	file_platform_common_proto_rawDescData = file_platform_common_proto_rawDesc
)

func file_platform_common_proto_rawDescGZIP() []byte {
	file_platform_common_proto_rawDescOnce.Do(func() {
		file_platform_common_proto_rawDescData = protoimpl.X.CompressGZIP(file_platform_common_proto_rawDescData)
	})
	return file_platform_common_proto_rawDescData
}

var file_platform_common_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_platform_common_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_platform_common_proto_goTypes = []interface{}{
	(Attrs)(0),                   // 0: platform.Attrs
	(State)(0),                   // 1: platform.State
	(*Comment)(nil),              // 2: platform.Comment
	(*CommentBlock)(nil),         // 3: platform.CommentBlock
	(*ReplyList)(nil),            // 4: platform.ReplyList
	(*CommentFilterOptions)(nil), // 5: platform.CommentFilterOptions
	(*LabelFilterOptions)(nil),   // 6: platform.LabelFilterOptions
	(*Subject)(nil),              // 7: platform.Subject
	(*Label)(nil),                // 8: platform.Label
	(*Relation)(nil),             // 9: platform.Relation
}
var file_platform_common_proto_depIdxs = []int32{
	2, // 0: platform.CommentBlock.rootComment:type_name -> platform.Comment
	4, // 1: platform.CommentBlock.replyList:type_name -> platform.ReplyList
	2, // 2: platform.ReplyList.comments:type_name -> platform.Comment
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}


func file_platform_common_proto_init() {
	if File_platform_common_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_platform_common_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Comment); i {
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
		file_platform_common_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CommentBlock); i {
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
		file_platform_common_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReplyList); i {
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
		file_platform_common_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CommentFilterOptions); i {
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
		file_platform_common_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LabelFilterOptions); i {
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
		file_platform_common_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Subject); i {
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
		file_platform_common_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Label); i {
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
		file_platform_common_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Relation); i {
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
	file_platform_common_proto_msgTypes[3].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_platform_common_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_platform_common_proto_goTypes,
		DependencyIndexes: file_platform_common_proto_depIdxs,
		EnumInfos:         file_platform_common_proto_enumTypes,
		MessageInfos:      file_platform_common_proto_msgTypes,
	}.Build()
	File_platform_common_proto = out.File
	file_platform_common_proto_rawDesc = nil
	file_platform_common_proto_goTypes = nil
	file_platform_common_proto_depIdxs = nil
}
