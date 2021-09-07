// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
// source: plugin.proto

package proto

import (
	v1 "github.com/aserto-dev/go-grpc/aserto/api/v1"
	_ "google.golang.org/genproto/googleapis/rpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	_ "google.golang.org/protobuf/types/known/structpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Use config from: https://github.com/aserto-dev/proto/blob/main/public/aserto/api/v1/connection.proto#L71 ??
type ConfigElementKind int32

const (
	ConfigElementKind_CONFIG_ELEMENT_KIND_UNKNOWN     ConfigElementKind = 0 // Unknown configuration element kind
	ConfigElementKind_CONFIG_ELEMENT_KIND_ATTRIBUTE   ConfigElementKind = 1 // Normal attribute
	ConfigElementKind_CONFIG_ELEMENT_KIND_SECRET      ConfigElementKind = 2 // Secret
	ConfigElementKind_CONFIG_ELEMENT_KIND_CERTIFICATE ConfigElementKind = 3 // Certificate
)

// Enum value maps for ConfigElementKind.
var (
	ConfigElementKind_name = map[int32]string{
		0: "CONFIG_ELEMENT_KIND_UNKNOWN",
		1: "CONFIG_ELEMENT_KIND_ATTRIBUTE",
		2: "CONFIG_ELEMENT_KIND_SECRET",
		3: "CONFIG_ELEMENT_KIND_CERTIFICATE",
	}
	ConfigElementKind_value = map[string]int32{
		"CONFIG_ELEMENT_KIND_UNKNOWN":     0,
		"CONFIG_ELEMENT_KIND_ATTRIBUTE":   1,
		"CONFIG_ELEMENT_KIND_SECRET":      2,
		"CONFIG_ELEMENT_KIND_CERTIFICATE": 3,
	}
)

func (x ConfigElementKind) Enum() *ConfigElementKind {
	p := new(ConfigElementKind)
	*p = x
	return p
}

func (x ConfigElementKind) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ConfigElementKind) Descriptor() protoreflect.EnumDescriptor {
	return file_plugin_proto_enumTypes[0].Descriptor()
}

func (ConfigElementKind) Type() protoreflect.EnumType {
	return &file_plugin_proto_enumTypes[0]
}

func (x ConfigElementKind) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ConfigElementKind.Descriptor instead.
func (ConfigElementKind) EnumDescriptor() ([]byte, []int) {
	return file_plugin_proto_rawDescGZIP(), []int{0}
}

type ConfigElementType int32

const (
	ConfigElementType_CONFIG_ELEMENT_TYPE_UNKNOWN ConfigElementType = 0
	ConfigElementType_CONFIG_ELEMENT_TYPE_STRING  ConfigElementType = 1
	ConfigElementType_CONFIG_ELEMENT_TYPE_INTEGER ConfigElementType = 2
	ConfigElementType_CONFIG_ELEMENT_TYPE_BOOLEAN ConfigElementType = 3
)

// Enum value maps for ConfigElementType.
var (
	ConfigElementType_name = map[int32]string{
		0: "CONFIG_ELEMENT_TYPE_UNKNOWN",
		1: "CONFIG_ELEMENT_TYPE_STRING",
		2: "CONFIG_ELEMENT_TYPE_INTEGER",
		3: "CONFIG_ELEMENT_TYPE_BOOLEAN",
	}
	ConfigElementType_value = map[string]int32{
		"CONFIG_ELEMENT_TYPE_UNKNOWN": 0,
		"CONFIG_ELEMENT_TYPE_STRING":  1,
		"CONFIG_ELEMENT_TYPE_INTEGER": 2,
		"CONFIG_ELEMENT_TYPE_BOOLEAN": 3,
	}
)

func (x ConfigElementType) Enum() *ConfigElementType {
	p := new(ConfigElementType)
	*p = x
	return p
}

func (x ConfigElementType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ConfigElementType) Descriptor() protoreflect.EnumDescriptor {
	return file_plugin_proto_enumTypes[1].Descriptor()
}

func (ConfigElementType) Type() protoreflect.EnumType {
	return &file_plugin_proto_enumTypes[1]
}

func (x ConfigElementType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ConfigElementType.Descriptor instead.
func (ConfigElementType) EnumDescriptor() ([]byte, []int) {
	return file_plugin_proto_rawDescGZIP(), []int{1}
}

type DisplayMode int32

const (
	DisplayMode_DISPLAY_MODE_UNKNOWN DisplayMode = 0
	DisplayMode_DISPLAY_MODE_NORMAL  DisplayMode = 1
	DisplayMode_DISPLAY_MODE_MASKED  DisplayMode = 2
)

// Enum value maps for DisplayMode.
var (
	DisplayMode_name = map[int32]string{
		0: "DISPLAY_MODE_UNKNOWN",
		1: "DISPLAY_MODE_NORMAL",
		2: "DISPLAY_MODE_MASKED",
	}
	DisplayMode_value = map[string]int32{
		"DISPLAY_MODE_UNKNOWN": 0,
		"DISPLAY_MODE_NORMAL":  1,
		"DISPLAY_MODE_MASKED":  2,
	}
)

func (x DisplayMode) Enum() *DisplayMode {
	p := new(DisplayMode)
	*p = x
	return p
}

func (x DisplayMode) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (DisplayMode) Descriptor() protoreflect.EnumDescriptor {
	return file_plugin_proto_enumTypes[2].Descriptor()
}

func (DisplayMode) Type() protoreflect.EnumType {
	return &file_plugin_proto_enumTypes[2]
}

func (x DisplayMode) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use DisplayMode.Descriptor instead.
func (DisplayMode) EnumDescriptor() ([]byte, []int) {
	return file_plugin_proto_rawDescGZIP(), []int{2}
}

type InfoRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *InfoRequest) Reset() {
	*x = InfoRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_plugin_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InfoRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InfoRequest) ProtoMessage() {}

func (x *InfoRequest) ProtoReflect() protoreflect.Message {
	mi := &file_plugin_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InfoRequest.ProtoReflect.Descriptor instead.
func (*InfoRequest) Descriptor() ([]byte, []int) {
	return file_plugin_proto_rawDescGZIP(), []int{0}
}

type InfoResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Use info from https://github.com/aserto-dev/proto/blob/main/public/aserto/common/info/v1/info.proto#L37 ??
	System  string           `protobuf:"bytes,1,opt,name=system,proto3" json:"system,omitempty"`
	Version string           `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`
	Build   string           `protobuf:"bytes,3,opt,name=build,proto3" json:"build,omitempty"`
	Config  []*ConfigElement `protobuf:"bytes,4,rep,name=config,proto3" json:"config,omitempty"`
}

func (x *InfoResponse) Reset() {
	*x = InfoResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_plugin_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InfoResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InfoResponse) ProtoMessage() {}

func (x *InfoResponse) ProtoReflect() protoreflect.Message {
	mi := &file_plugin_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InfoResponse.ProtoReflect.Descriptor instead.
func (*InfoResponse) Descriptor() ([]byte, []int) {
	return file_plugin_proto_rawDescGZIP(), []int{1}
}

func (x *InfoResponse) GetSystem() string {
	if x != nil {
		return x.System
	}
	return ""
}

func (x *InfoResponse) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *InfoResponse) GetBuild() string {
	if x != nil {
		return x.Build
	}
	return ""
}

func (x *InfoResponse) GetConfig() []*ConfigElement {
	if x != nil {
		return x.Config
	}
	return nil
}

type ConfigElement struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          int32             `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Kind        ConfigElementKind `protobuf:"varint,2,opt,name=kind,proto3,enum=proto.ConfigElementKind" json:"kind,omitempty"`
	Type        ConfigElementType `protobuf:"varint,3,opt,name=type,proto3,enum=proto.ConfigElementType" json:"type,omitempty"`
	Name        string            `protobuf:"bytes,4,opt,name=name,proto3" json:"name,omitempty"`
	Description string            `protobuf:"bytes,5,opt,name=description,proto3" json:"description,omitempty"`
	Usage       string            `protobuf:"bytes,6,opt,name=usage,proto3" json:"usage,omitempty"`
	Mode        DisplayMode       `protobuf:"varint,7,opt,name=mode,proto3,enum=proto.DisplayMode" json:"mode,omitempty"`
	ReadOnly    bool              `protobuf:"varint,8,opt,name=read_only,json=readOnly,proto3" json:"read_only,omitempty"`
}

func (x *ConfigElement) Reset() {
	*x = ConfigElement{}
	if protoimpl.UnsafeEnabled {
		mi := &file_plugin_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConfigElement) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConfigElement) ProtoMessage() {}

func (x *ConfigElement) ProtoReflect() protoreflect.Message {
	mi := &file_plugin_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConfigElement.ProtoReflect.Descriptor instead.
func (*ConfigElement) Descriptor() ([]byte, []int) {
	return file_plugin_proto_rawDescGZIP(), []int{2}
}

func (x *ConfigElement) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *ConfigElement) GetKind() ConfigElementKind {
	if x != nil {
		return x.Kind
	}
	return ConfigElementKind_CONFIG_ELEMENT_KIND_UNKNOWN
}

func (x *ConfigElement) GetType() ConfigElementType {
	if x != nil {
		return x.Type
	}
	return ConfigElementType_CONFIG_ELEMENT_TYPE_UNKNOWN
}

func (x *ConfigElement) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ConfigElement) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *ConfigElement) GetUsage() string {
	if x != nil {
		return x.Usage
	}
	return ""
}

func (x *ConfigElement) GetMode() DisplayMode {
	if x != nil {
		return x.Mode
	}
	return DisplayMode_DISPLAY_MODE_UNKNOWN
}

func (x *ConfigElement) GetReadOnly() bool {
	if x != nil {
		return x.ReadOnly
	}
	return false
}

type ImportRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// needs to be struct
	Options map[string]string `protobuf:"bytes,1,rep,name=options,proto3" json:"options,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// Types that are assignable to Data:
	//	*ImportRequest_User
	//	*ImportRequest_UserExt
	Data isImportRequest_Data `protobuf_oneof:"data"`
}

func (x *ImportRequest) Reset() {
	*x = ImportRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_plugin_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ImportRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ImportRequest) ProtoMessage() {}

func (x *ImportRequest) ProtoReflect() protoreflect.Message {
	mi := &file_plugin_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ImportRequest.ProtoReflect.Descriptor instead.
func (*ImportRequest) Descriptor() ([]byte, []int) {
	return file_plugin_proto_rawDescGZIP(), []int{3}
}

func (x *ImportRequest) GetOptions() map[string]string {
	if x != nil {
		return x.Options
	}
	return nil
}

func (m *ImportRequest) GetData() isImportRequest_Data {
	if m != nil {
		return m.Data
	}
	return nil
}

func (x *ImportRequest) GetUser() *v1.User {
	if x, ok := x.GetData().(*ImportRequest_User); ok {
		return x.User
	}
	return nil
}

func (x *ImportRequest) GetUserExt() *v1.UserExt {
	if x, ok := x.GetData().(*ImportRequest_UserExt); ok {
		return x.UserExt
	}
	return nil
}

type isImportRequest_Data interface {
	isImportRequest_Data()
}

type ImportRequest_User struct {
	User *v1.User `protobuf:"bytes,3,opt,name=user,proto3,oneof"`
}

type ImportRequest_UserExt struct {
	UserExt *v1.UserExt `protobuf:"bytes,4,opt,name=user_ext,json=userExt,proto3,oneof"`
}

func (*ImportRequest_User) isImportRequest_Data() {}

func (*ImportRequest_UserExt) isImportRequest_Data() {}

type ImportResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SuccededCount int32 `protobuf:"varint,1,opt,name=succededCount,proto3" json:"succededCount,omitempty"`
	FailCount     int32 `protobuf:"varint,2,opt,name=failCount,proto3" json:"failCount,omitempty"`
}

func (x *ImportResponse) Reset() {
	*x = ImportResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_plugin_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ImportResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ImportResponse) ProtoMessage() {}

func (x *ImportResponse) ProtoReflect() protoreflect.Message {
	mi := &file_plugin_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ImportResponse.ProtoReflect.Descriptor instead.
func (*ImportResponse) Descriptor() ([]byte, []int) {
	return file_plugin_proto_rawDescGZIP(), []int{4}
}

func (x *ImportResponse) GetSuccededCount() int32 {
	if x != nil {
		return x.SuccededCount
	}
	return 0
}

func (x *ImportResponse) GetFailCount() int32 {
	if x != nil {
		return x.FailCount
	}
	return 0
}

type ExportRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Options map[string]string `protobuf:"bytes,1,rep,name=options,proto3" json:"options,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *ExportRequest) Reset() {
	*x = ExportRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_plugin_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ExportRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExportRequest) ProtoMessage() {}

func (x *ExportRequest) ProtoReflect() protoreflect.Message {
	mi := &file_plugin_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExportRequest.ProtoReflect.Descriptor instead.
func (*ExportRequest) Descriptor() ([]byte, []int) {
	return file_plugin_proto_rawDescGZIP(), []int{5}
}

func (x *ExportRequest) GetOptions() map[string]string {
	if x != nil {
		return x.Options
	}
	return nil
}

type ExportResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Data:
	//	*ExportResponse_User
	//	*ExportResponse_UserExt
	Data isExportResponse_Data `protobuf_oneof:"data"`
}

func (x *ExportResponse) Reset() {
	*x = ExportResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_plugin_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ExportResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExportResponse) ProtoMessage() {}

func (x *ExportResponse) ProtoReflect() protoreflect.Message {
	mi := &file_plugin_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExportResponse.ProtoReflect.Descriptor instead.
func (*ExportResponse) Descriptor() ([]byte, []int) {
	return file_plugin_proto_rawDescGZIP(), []int{6}
}

func (m *ExportResponse) GetData() isExportResponse_Data {
	if m != nil {
		return m.Data
	}
	return nil
}

func (x *ExportResponse) GetUser() *v1.User {
	if x, ok := x.GetData().(*ExportResponse_User); ok {
		return x.User
	}
	return nil
}

func (x *ExportResponse) GetUserExt() *v1.UserExt {
	if x, ok := x.GetData().(*ExportResponse_UserExt); ok {
		return x.UserExt
	}
	return nil
}

type isExportResponse_Data interface {
	isExportResponse_Data()
}

type ExportResponse_User struct {
	User *v1.User `protobuf:"bytes,3,opt,name=user,proto3,oneof"`
}

type ExportResponse_UserExt struct {
	UserExt *v1.UserExt `protobuf:"bytes,4,opt,name=user_ext,json=userExt,proto3,oneof"`
}

func (*ExportResponse_User) isExportResponse_Data() {}

func (*ExportResponse_UserExt) isExportResponse_Data() {}

var File_plugin_proto protoreflect.FileDescriptor

var file_plugin_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x18, 0x61, 0x73, 0x65, 0x72, 0x74, 0x6f, 0x2f, 0x61, 0x70,
	0x69, 0x2f, 0x76, 0x31, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x72, 0x70, 0x63, 0x2f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x0d, 0x0a, 0x0b, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x84, 0x01, 0x0a, 0x0c, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x12, 0x18,
	0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x62, 0x75, 0x69, 0x6c,
	0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x12, 0x2c,
	0x0a, 0x06, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x45, 0x6c, 0x65,
	0x6d, 0x65, 0x6e, 0x74, 0x52, 0x06, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x22, 0x8c, 0x02, 0x0a,
	0x0d, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x45, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x69, 0x64, 0x12, 0x2c,
	0x0a, 0x04, 0x6b, 0x69, 0x6e, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x18, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x45, 0x6c, 0x65, 0x6d, 0x65,
	0x6e, 0x74, 0x4b, 0x69, 0x6e, 0x64, 0x52, 0x04, 0x6b, 0x69, 0x6e, 0x64, 0x12, 0x2c, 0x0a, 0x04,
	0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x18, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x45, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74,
	0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x20,
	0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x14, 0x0a, 0x05, 0x75, 0x73, 0x61, 0x67, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x75, 0x73, 0x61, 0x67, 0x65, 0x12, 0x26, 0x0a, 0x04, 0x6d, 0x6f, 0x64, 0x65, 0x18, 0x07,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x12, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x44, 0x69, 0x73,
	0x70, 0x6c, 0x61, 0x79, 0x4d, 0x6f, 0x64, 0x65, 0x52, 0x04, 0x6d, 0x6f, 0x64, 0x65, 0x12, 0x1b,
	0x0a, 0x09, 0x72, 0x65, 0x61, 0x64, 0x5f, 0x6f, 0x6e, 0x6c, 0x79, 0x18, 0x08, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x08, 0x72, 0x65, 0x61, 0x64, 0x4f, 0x6e, 0x6c, 0x79, 0x22, 0xf0, 0x01, 0x0a, 0x0d,
	0x49, 0x6d, 0x70, 0x6f, 0x72, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x3b, 0x0a,
	0x07, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x21,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x49, 0x6d, 0x70, 0x6f, 0x72, 0x74, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x2e, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x45, 0x6e, 0x74, 0x72,
	0x79, 0x52, 0x07, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x29, 0x0a, 0x04, 0x75, 0x73,
	0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x61, 0x73, 0x65, 0x72, 0x74,
	0x6f, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x48, 0x00, 0x52,
	0x04, 0x75, 0x73, 0x65, 0x72, 0x12, 0x33, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x65, 0x78,
	0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x61, 0x73, 0x65, 0x72, 0x74, 0x6f,
	0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x45, 0x78, 0x74, 0x48,
	0x00, 0x52, 0x07, 0x75, 0x73, 0x65, 0x72, 0x45, 0x78, 0x74, 0x1a, 0x3a, 0x0a, 0x0c, 0x4f, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65,
	0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x42, 0x06, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x54,
	0x0a, 0x0e, 0x49, 0x6d, 0x70, 0x6f, 0x72, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x24, 0x0a, 0x0d, 0x73, 0x75, 0x63, 0x63, 0x65, 0x64, 0x65, 0x64, 0x43, 0x6f, 0x75, 0x6e,
	0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0d, 0x73, 0x75, 0x63, 0x63, 0x65, 0x64, 0x65,
	0x64, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x66, 0x61, 0x69, 0x6c, 0x43, 0x6f,
	0x75, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x66, 0x61, 0x69, 0x6c, 0x43,
	0x6f, 0x75, 0x6e, 0x74, 0x22, 0x88, 0x01, 0x0a, 0x0d, 0x45, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x3b, 0x0a, 0x07, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x45, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x4f, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x07, 0x6f, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x1a, 0x3a, 0x0a, 0x0c, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22,
	0x78, 0x0a, 0x0e, 0x45, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x29, 0x0a, 0x04, 0x75, 0x73, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x13, 0x2e, 0x61, 0x73, 0x65, 0x72, 0x74, 0x6f, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e,
	0x55, 0x73, 0x65, 0x72, 0x48, 0x00, 0x52, 0x04, 0x75, 0x73, 0x65, 0x72, 0x12, 0x33, 0x0a, 0x08,
	0x75, 0x73, 0x65, 0x72, 0x5f, 0x65, 0x78, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16,
	0x2e, 0x61, 0x73, 0x65, 0x72, 0x74, 0x6f, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x55,
	0x73, 0x65, 0x72, 0x45, 0x78, 0x74, 0x48, 0x00, 0x52, 0x07, 0x75, 0x73, 0x65, 0x72, 0x45, 0x78,
	0x74, 0x42, 0x06, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x2a, 0x9c, 0x01, 0x0a, 0x11, 0x43, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x45, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x4b, 0x69, 0x6e, 0x64, 0x12,
	0x1f, 0x0a, 0x1b, 0x43, 0x4f, 0x4e, 0x46, 0x49, 0x47, 0x5f, 0x45, 0x4c, 0x45, 0x4d, 0x45, 0x4e,
	0x54, 0x5f, 0x4b, 0x49, 0x4e, 0x44, 0x5f, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00,
	0x12, 0x21, 0x0a, 0x1d, 0x43, 0x4f, 0x4e, 0x46, 0x49, 0x47, 0x5f, 0x45, 0x4c, 0x45, 0x4d, 0x45,
	0x4e, 0x54, 0x5f, 0x4b, 0x49, 0x4e, 0x44, 0x5f, 0x41, 0x54, 0x54, 0x52, 0x49, 0x42, 0x55, 0x54,
	0x45, 0x10, 0x01, 0x12, 0x1e, 0x0a, 0x1a, 0x43, 0x4f, 0x4e, 0x46, 0x49, 0x47, 0x5f, 0x45, 0x4c,
	0x45, 0x4d, 0x45, 0x4e, 0x54, 0x5f, 0x4b, 0x49, 0x4e, 0x44, 0x5f, 0x53, 0x45, 0x43, 0x52, 0x45,
	0x54, 0x10, 0x02, 0x12, 0x23, 0x0a, 0x1f, 0x43, 0x4f, 0x4e, 0x46, 0x49, 0x47, 0x5f, 0x45, 0x4c,
	0x45, 0x4d, 0x45, 0x4e, 0x54, 0x5f, 0x4b, 0x49, 0x4e, 0x44, 0x5f, 0x43, 0x45, 0x52, 0x54, 0x49,
	0x46, 0x49, 0x43, 0x41, 0x54, 0x45, 0x10, 0x03, 0x2a, 0x96, 0x01, 0x0a, 0x11, 0x43, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x45, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1f,
	0x0a, 0x1b, 0x43, 0x4f, 0x4e, 0x46, 0x49, 0x47, 0x5f, 0x45, 0x4c, 0x45, 0x4d, 0x45, 0x4e, 0x54,
	0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12,
	0x1e, 0x0a, 0x1a, 0x43, 0x4f, 0x4e, 0x46, 0x49, 0x47, 0x5f, 0x45, 0x4c, 0x45, 0x4d, 0x45, 0x4e,
	0x54, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x53, 0x54, 0x52, 0x49, 0x4e, 0x47, 0x10, 0x01, 0x12,
	0x1f, 0x0a, 0x1b, 0x43, 0x4f, 0x4e, 0x46, 0x49, 0x47, 0x5f, 0x45, 0x4c, 0x45, 0x4d, 0x45, 0x4e,
	0x54, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x49, 0x4e, 0x54, 0x45, 0x47, 0x45, 0x52, 0x10, 0x02,
	0x12, 0x1f, 0x0a, 0x1b, 0x43, 0x4f, 0x4e, 0x46, 0x49, 0x47, 0x5f, 0x45, 0x4c, 0x45, 0x4d, 0x45,
	0x4e, 0x54, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x42, 0x4f, 0x4f, 0x4c, 0x45, 0x41, 0x4e, 0x10,
	0x03, 0x2a, 0x59, 0x0a, 0x0b, 0x44, 0x69, 0x73, 0x70, 0x6c, 0x61, 0x79, 0x4d, 0x6f, 0x64, 0x65,
	0x12, 0x18, 0x0a, 0x14, 0x44, 0x49, 0x53, 0x50, 0x4c, 0x41, 0x59, 0x5f, 0x4d, 0x4f, 0x44, 0x45,
	0x5f, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x17, 0x0a, 0x13, 0x44, 0x49,
	0x53, 0x50, 0x4c, 0x41, 0x59, 0x5f, 0x4d, 0x4f, 0x44, 0x45, 0x5f, 0x4e, 0x4f, 0x52, 0x4d, 0x41,
	0x4c, 0x10, 0x01, 0x12, 0x17, 0x0a, 0x13, 0x44, 0x49, 0x53, 0x50, 0x4c, 0x41, 0x59, 0x5f, 0x4d,
	0x4f, 0x44, 0x45, 0x5f, 0x4d, 0x41, 0x53, 0x4b, 0x45, 0x44, 0x10, 0x02, 0x32, 0xab, 0x01, 0x0a,
	0x06, 0x50, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x12, 0x2f, 0x0a, 0x04, 0x49, 0x6e, 0x66, 0x6f, 0x12,
	0x12, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x13, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x49, 0x6e, 0x66, 0x6f,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x37, 0x0a, 0x06, 0x49, 0x6d, 0x70, 0x6f,
	0x72, 0x74, 0x12, 0x14, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x49, 0x6d, 0x70, 0x6f, 0x72,
	0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x49, 0x6d, 0x70, 0x6f, 0x72, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x28,
	0x01, 0x12, 0x37, 0x0a, 0x06, 0x45, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x12, 0x14, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x45, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x15, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x45, 0x78, 0x70, 0x6f, 0x72, 0x74,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x30, 0x01, 0x42, 0x32, 0x5a, 0x30, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x73, 0x65, 0x72, 0x74, 0x6f, 0x2d,
	0x64, 0x65, 0x76, 0x2f, 0x61, 0x73, 0x65, 0x72, 0x74, 0x6f, 0x2d, 0x69, 0x64, 0x70, 0x2f, 0x70,
	0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x3b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_plugin_proto_rawDescOnce sync.Once
	file_plugin_proto_rawDescData = file_plugin_proto_rawDesc
)

func file_plugin_proto_rawDescGZIP() []byte {
	file_plugin_proto_rawDescOnce.Do(func() {
		file_plugin_proto_rawDescData = protoimpl.X.CompressGZIP(file_plugin_proto_rawDescData)
	})
	return file_plugin_proto_rawDescData
}

var file_plugin_proto_enumTypes = make([]protoimpl.EnumInfo, 3)
var file_plugin_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_plugin_proto_goTypes = []interface{}{
	(ConfigElementKind)(0), // 0: proto.ConfigElementKind
	(ConfigElementType)(0), // 1: proto.ConfigElementType
	(DisplayMode)(0),       // 2: proto.DisplayMode
	(*InfoRequest)(nil),    // 3: proto.InfoRequest
	(*InfoResponse)(nil),   // 4: proto.InfoResponse
	(*ConfigElement)(nil),  // 5: proto.ConfigElement
	(*ImportRequest)(nil),  // 6: proto.ImportRequest
	(*ImportResponse)(nil), // 7: proto.ImportResponse
	(*ExportRequest)(nil),  // 8: proto.ExportRequest
	(*ExportResponse)(nil), // 9: proto.ExportResponse
	nil,                    // 10: proto.ImportRequest.OptionsEntry
	nil,                    // 11: proto.ExportRequest.OptionsEntry
	(*v1.User)(nil),        // 12: aserto.api.v1.User
	(*v1.UserExt)(nil),     // 13: aserto.api.v1.UserExt
}
var file_plugin_proto_depIdxs = []int32{
	5,  // 0: proto.InfoResponse.config:type_name -> proto.ConfigElement
	0,  // 1: proto.ConfigElement.kind:type_name -> proto.ConfigElementKind
	1,  // 2: proto.ConfigElement.type:type_name -> proto.ConfigElementType
	2,  // 3: proto.ConfigElement.mode:type_name -> proto.DisplayMode
	10, // 4: proto.ImportRequest.options:type_name -> proto.ImportRequest.OptionsEntry
	12, // 5: proto.ImportRequest.user:type_name -> aserto.api.v1.User
	13, // 6: proto.ImportRequest.user_ext:type_name -> aserto.api.v1.UserExt
	11, // 7: proto.ExportRequest.options:type_name -> proto.ExportRequest.OptionsEntry
	12, // 8: proto.ExportResponse.user:type_name -> aserto.api.v1.User
	13, // 9: proto.ExportResponse.user_ext:type_name -> aserto.api.v1.UserExt
	3,  // 10: proto.Plugin.Info:input_type -> proto.InfoRequest
	6,  // 11: proto.Plugin.Import:input_type -> proto.ImportRequest
	8,  // 12: proto.Plugin.Export:input_type -> proto.ExportRequest
	4,  // 13: proto.Plugin.Info:output_type -> proto.InfoResponse
	7,  // 14: proto.Plugin.Import:output_type -> proto.ImportResponse
	9,  // 15: proto.Plugin.Export:output_type -> proto.ExportResponse
	13, // [13:16] is the sub-list for method output_type
	10, // [10:13] is the sub-list for method input_type
	10, // [10:10] is the sub-list for extension type_name
	10, // [10:10] is the sub-list for extension extendee
	0,  // [0:10] is the sub-list for field type_name
}

func init() { file_plugin_proto_init() }
func file_plugin_proto_init() {
	if File_plugin_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_plugin_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InfoRequest); i {
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
		file_plugin_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InfoResponse); i {
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
		file_plugin_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConfigElement); i {
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
		file_plugin_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ImportRequest); i {
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
		file_plugin_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ImportResponse); i {
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
		file_plugin_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ExportRequest); i {
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
		file_plugin_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ExportResponse); i {
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
	file_plugin_proto_msgTypes[3].OneofWrappers = []interface{}{
		(*ImportRequest_User)(nil),
		(*ImportRequest_UserExt)(nil),
	}
	file_plugin_proto_msgTypes[6].OneofWrappers = []interface{}{
		(*ExportResponse_User)(nil),
		(*ExportResponse_UserExt)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_plugin_proto_rawDesc,
			NumEnums:      3,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_plugin_proto_goTypes,
		DependencyIndexes: file_plugin_proto_depIdxs,
		EnumInfos:         file_plugin_proto_enumTypes,
		MessageInfos:      file_plugin_proto_msgTypes,
	}.Build()
	File_plugin_proto = out.File
	file_plugin_proto_rawDesc = nil
	file_plugin_proto_goTypes = nil
	file_plugin_proto_depIdxs = nil
}
