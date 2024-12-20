// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1-devel
// 	protoc        v3.19.1
// source: cfg.proto

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

// 测试message's struct tag
// @StructTagOfExample
type Example struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Int32Field  int32             `protobuf:"varint,1,opt,name=Int32Field,proto3" json:"Int32Field,omitempty"`
	StringField string            `protobuf:"bytes,2,opt,name=StringField,proto3" json:"StringField,omitempty"`
	FloatField  float32           `protobuf:"fixed32,3,opt,name=FloatField,proto3" json:"FloatField,omitempty"`
	Int32Slice  []int32           `protobuf:"varint,4,rep,packed,name=Int32Slice,proto3" json:"Int32Slice,omitempty"`
	StringSlice []string          `protobuf:"bytes,5,rep,name=StringSlice,proto3" json:"StringSlice,omitempty"`
	SingleChild *Child            `protobuf:"bytes,6,opt,name=SingleChild,proto3" json:"SingleChild,omitempty"`
	ChildSlice  []*Child          `protobuf:"bytes,7,rep,name=ChildSlice,proto3" json:"ChildSlice,omitempty"`
	NormalMap   map[int32]string  `protobuf:"bytes,8,rep,name=NormalMap,proto3" json:"NormalMap,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	ChildMap    map[string]*Child `protobuf:"bytes,9,rep,name=ChildMap,proto3" json:"ChildMap,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *Example) Reset() {
	*x = Example{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cfg_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Example) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Example) ProtoMessage() {}

func (x *Example) ProtoReflect() protoreflect.Message {
	mi := &file_cfg_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Example.ProtoReflect.Descriptor instead.
func (*Example) Descriptor() ([]byte, []int) {
	return file_cfg_proto_rawDescGZIP(), []int{0}
}

func (x *Example) GetInt32Field() int32 {
	if x != nil {
		return x.Int32Field
	}
	return 0
}

func (x *Example) GetStringField() string {
	if x != nil {
		return x.StringField
	}
	return ""
}

func (x *Example) GetFloatField() float32 {
	if x != nil {
		return x.FloatField
	}
	return 0
}

func (x *Example) GetInt32Slice() []int32 {
	if x != nil {
		return x.Int32Slice
	}
	return nil
}

func (x *Example) GetStringSlice() []string {
	if x != nil {
		return x.StringSlice
	}
	return nil
}

func (x *Example) GetSingleChild() *Child {
	if x != nil {
		return x.SingleChild
	}
	return nil
}

func (x *Example) GetChildSlice() []*Child {
	if x != nil {
		return x.ChildSlice
	}
	return nil
}

func (x *Example) GetNormalMap() map[int32]string {
	if x != nil {
		return x.NormalMap
	}
	return nil
}

func (x *Example) GetChildMap() map[string]*Child {
	if x != nil {
		return x.ChildMap
	}
	return nil
}

// 测试message's struct tag with value
// @StructTagOfExample:testValue
type Example2 struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Int32Field  int32  `protobuf:"varint,1,opt,name=Int32Field,proto3" json:"Int32Field,omitempty"`
	StringField string `protobuf:"bytes,2,opt,name=StringField,proto3" json:"StringField,omitempty"`
}

func (x *Example2) Reset() {
	*x = Example2{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cfg_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Example2) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Example2) ProtoMessage() {}

func (x *Example2) ProtoReflect() protoreflect.Message {
	mi := &file_cfg_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Example2.ProtoReflect.Descriptor instead.
func (*Example2) Descriptor() ([]byte, []int) {
	return file_cfg_proto_rawDescGZIP(), []int{1}
}

func (x *Example2) GetInt32Field() int32 {
	if x != nil {
		return x.Int32Field
	}
	return 0
}

func (x *Example2) GetStringField() string {
	if x != nil {
		return x.StringField
	}
	return ""
}

type ExampleWithoutTag struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Int32Field  int32  `protobuf:"varint,1,opt,name=Int32Field,proto3" json:"Int32Field,omitempty"`
	StringField string `protobuf:"bytes,2,opt,name=StringField,proto3" json:"StringField,omitempty"`
}

func (x *ExampleWithoutTag) Reset() {
	*x = ExampleWithoutTag{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cfg_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ExampleWithoutTag) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExampleWithoutTag) ProtoMessage() {}

func (x *ExampleWithoutTag) ProtoReflect() protoreflect.Message {
	mi := &file_cfg_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExampleWithoutTag.ProtoReflect.Descriptor instead.
func (*ExampleWithoutTag) Descriptor() ([]byte, []int) {
	return file_cfg_proto_rawDescGZIP(), []int{2}
}

func (x *ExampleWithoutTag) GetInt32Field() int32 {
	if x != nil {
		return x.Int32Field
	}
	return 0
}

func (x *ExampleWithoutTag) GetStringField() string {
	if x != nil {
		return x.StringField
	}
	return ""
}

type Child struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Int32Field  int32  `protobuf:"varint,1,opt,name=Int32Field,proto3" json:"Int32Field,omitempty"`
	StringField string `protobuf:"bytes,2,opt,name=StringField,proto3" json:"StringField,omitempty"`
}

func (x *Child) Reset() {
	*x = Child{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cfg_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Child) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Child) ProtoMessage() {}

func (x *Child) ProtoReflect() protoreflect.Message {
	mi := &file_cfg_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Child.ProtoReflect.Descriptor instead.
func (*Child) Descriptor() ([]byte, []int) {
	return file_cfg_proto_rawDescGZIP(), []int{3}
}

func (x *Child) GetInt32Field() int32 {
	if x != nil {
		return x.Int32Field
	}
	return 0
}

func (x *Child) GetStringField() string {
	if x != nil {
		return x.StringField
	}
	return ""
}

var File_cfg_proto protoreflect.FileDescriptor

var file_cfg_proto_rawDesc = []byte{
	0x0a, 0x09, 0x63, 0x66, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02, 0x70, 0x62, 0x22,
	0xfc, 0x03, 0x0a, 0x07, 0x45, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x49,
	0x6e, 0x74, 0x33, 0x32, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x0a, 0x49, 0x6e, 0x74, 0x33, 0x32, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x12, 0x20, 0x0a, 0x0b, 0x53,
	0x74, 0x72, 0x69, 0x6e, 0x67, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0b, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x12, 0x1e, 0x0a,
	0x0a, 0x46, 0x6c, 0x6f, 0x61, 0x74, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x02, 0x52, 0x0a, 0x46, 0x6c, 0x6f, 0x61, 0x74, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x12, 0x1e, 0x0a,
	0x0a, 0x49, 0x6e, 0x74, 0x33, 0x32, 0x53, 0x6c, 0x69, 0x63, 0x65, 0x18, 0x04, 0x20, 0x03, 0x28,
	0x05, 0x52, 0x0a, 0x49, 0x6e, 0x74, 0x33, 0x32, 0x53, 0x6c, 0x69, 0x63, 0x65, 0x12, 0x20, 0x0a,
	0x0b, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x53, 0x6c, 0x69, 0x63, 0x65, 0x18, 0x05, 0x20, 0x03,
	0x28, 0x09, 0x52, 0x0b, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x53, 0x6c, 0x69, 0x63, 0x65, 0x12,
	0x2b, 0x0a, 0x0b, 0x53, 0x69, 0x6e, 0x67, 0x6c, 0x65, 0x43, 0x68, 0x69, 0x6c, 0x64, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x09, 0x2e, 0x70, 0x62, 0x2e, 0x43, 0x68, 0x69, 0x6c, 0x64, 0x52,
	0x0b, 0x53, 0x69, 0x6e, 0x67, 0x6c, 0x65, 0x43, 0x68, 0x69, 0x6c, 0x64, 0x12, 0x29, 0x0a, 0x0a,
	0x43, 0x68, 0x69, 0x6c, 0x64, 0x53, 0x6c, 0x69, 0x63, 0x65, 0x18, 0x07, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x09, 0x2e, 0x70, 0x62, 0x2e, 0x43, 0x68, 0x69, 0x6c, 0x64, 0x52, 0x0a, 0x43, 0x68, 0x69,
	0x6c, 0x64, 0x53, 0x6c, 0x69, 0x63, 0x65, 0x12, 0x38, 0x0a, 0x09, 0x4e, 0x6f, 0x72, 0x6d, 0x61,
	0x6c, 0x4d, 0x61, 0x70, 0x18, 0x08, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x70, 0x62, 0x2e,
	0x45, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x4e, 0x6f, 0x72, 0x6d, 0x61, 0x6c, 0x4d, 0x61,
	0x70, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x09, 0x4e, 0x6f, 0x72, 0x6d, 0x61, 0x6c, 0x4d, 0x61,
	0x70, 0x12, 0x35, 0x0a, 0x08, 0x43, 0x68, 0x69, 0x6c, 0x64, 0x4d, 0x61, 0x70, 0x18, 0x09, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x70, 0x62, 0x2e, 0x45, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65,
	0x2e, 0x43, 0x68, 0x69, 0x6c, 0x64, 0x4d, 0x61, 0x70, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x08,
	0x43, 0x68, 0x69, 0x6c, 0x64, 0x4d, 0x61, 0x70, 0x1a, 0x3c, 0x0a, 0x0e, 0x4e, 0x6f, 0x72, 0x6d,
	0x61, 0x6c, 0x4d, 0x61, 0x70, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65,
	0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0x46, 0x0a, 0x0d, 0x43, 0x68, 0x69, 0x6c, 0x64, 0x4d,
	0x61, 0x70, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x1f, 0x0a, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x09, 0x2e, 0x70, 0x62, 0x2e, 0x43, 0x68,
	0x69, 0x6c, 0x64, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x4c,
	0x0a, 0x08, 0x45, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x32, 0x12, 0x1e, 0x0a, 0x0a, 0x49, 0x6e,
	0x74, 0x33, 0x32, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a,
	0x49, 0x6e, 0x74, 0x33, 0x32, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x12, 0x20, 0x0a, 0x0b, 0x53, 0x74,
	0x72, 0x69, 0x6e, 0x67, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0b, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x22, 0x55, 0x0a, 0x11,
	0x45, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x57, 0x69, 0x74, 0x68, 0x6f, 0x75, 0x74, 0x54, 0x61,
	0x67, 0x12, 0x1e, 0x0a, 0x0a, 0x49, 0x6e, 0x74, 0x33, 0x32, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x49, 0x6e, 0x74, 0x33, 0x32, 0x46, 0x69, 0x65, 0x6c,
	0x64, 0x12, 0x20, 0x0a, 0x0b, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x46, 0x69, 0x65, 0x6c, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x46, 0x69,
	0x65, 0x6c, 0x64, 0x22, 0x49, 0x0a, 0x05, 0x43, 0x68, 0x69, 0x6c, 0x64, 0x12, 0x1e, 0x0a, 0x0a,
	0x49, 0x6e, 0x74, 0x33, 0x32, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x0a, 0x49, 0x6e, 0x74, 0x33, 0x32, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x12, 0x20, 0x0a, 0x0b,
	0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0b, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x42, 0x05,
	0x5a, 0x03, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_cfg_proto_rawDescOnce sync.Once
	file_cfg_proto_rawDescData = file_cfg_proto_rawDesc
)

func file_cfg_proto_rawDescGZIP() []byte {
	file_cfg_proto_rawDescOnce.Do(func() {
		file_cfg_proto_rawDescData = protoimpl.X.CompressGZIP(file_cfg_proto_rawDescData)
	})
	return file_cfg_proto_rawDescData
}

var file_cfg_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_cfg_proto_goTypes = []interface{}{
	(*Example)(nil),           // 0: pb.Example
	(*Example2)(nil),          // 1: pb.Example2
	(*ExampleWithoutTag)(nil), // 2: pb.ExampleWithoutTag
	(*Child)(nil),             // 3: pb.Child
	nil,                       // 4: pb.Example.NormalMapEntry
	nil,                       // 5: pb.Example.ChildMapEntry
}
var file_cfg_proto_depIdxs = []int32{
	3, // 0: pb.Example.SingleChild:type_name -> pb.Child
	3, // 1: pb.Example.ChildSlice:type_name -> pb.Child
	4, // 2: pb.Example.NormalMap:type_name -> pb.Example.NormalMapEntry
	5, // 3: pb.Example.ChildMap:type_name -> pb.Example.ChildMapEntry
	3, // 4: pb.Example.ChildMapEntry.value:type_name -> pb.Child
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_cfg_proto_init() }
func file_cfg_proto_init() {
	if File_cfg_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_cfg_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Example); i {
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
		file_cfg_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Example2); i {
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
		file_cfg_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ExampleWithoutTag); i {
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
		file_cfg_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Child); i {
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
			RawDescriptor: file_cfg_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_cfg_proto_goTypes,
		DependencyIndexes: file_cfg_proto_depIdxs,
		MessageInfos:      file_cfg_proto_msgTypes,
	}.Build()
	File_cfg_proto = out.File
	file_cfg_proto_rawDesc = nil
	file_cfg_proto_goTypes = nil
	file_cfg_proto_depIdxs = nil
}
