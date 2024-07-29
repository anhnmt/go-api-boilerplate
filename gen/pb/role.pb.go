// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        (unknown)
// source: role.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	descriptorpb "google.golang.org/protobuf/types/descriptorpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Role int32

const (
	Role_USER      Role = 0
	Role_ADMIN     Role = 1
	Role_ALL_ROLES Role = 100
)

// Enum value maps for Role.
var (
	Role_name = map[int32]string{
		0:   "USER",
		1:   "ADMIN",
		100: "ALL_ROLES",
	}
	Role_value = map[string]int32{
		"USER":      0,
		"ADMIN":     1,
		"ALL_ROLES": 100,
	}
)

func (x Role) Enum() *Role {
	p := new(Role)
	*p = x
	return p
}

func (x Role) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Role) Descriptor() protoreflect.EnumDescriptor {
	return file_role_proto_enumTypes[0].Descriptor()
}

func (Role) Type() protoreflect.EnumType {
	return &file_role_proto_enumTypes[0]
}

func (x Role) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Role.Descriptor instead.
func (Role) EnumDescriptor() ([]byte, []int) {
	return file_role_proto_rawDescGZIP(), []int{0}
}

type RoleOptions struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Defaults  []Role `protobuf:"varint,1,rep,packed,name=defaults,proto3,enum=role.v1.Role" json:"defaults,omitempty"`
	Abilities []Role `protobuf:"varint,2,rep,packed,name=abilities,proto3,enum=role.v1.Role" json:"abilities,omitempty"`
}

func (x *RoleOptions) Reset() {
	*x = RoleOptions{}
	if protoimpl.UnsafeEnabled {
		mi := &file_role_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RoleOptions) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RoleOptions) ProtoMessage() {}

func (x *RoleOptions) ProtoReflect() protoreflect.Message {
	mi := &file_role_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RoleOptions.ProtoReflect.Descriptor instead.
func (*RoleOptions) Descriptor() ([]byte, []int) {
	return file_role_proto_rawDescGZIP(), []int{0}
}

func (x *RoleOptions) GetDefaults() []Role {
	if x != nil {
		return x.Defaults
	}
	return nil
}

func (x *RoleOptions) GetAbilities() []Role {
	if x != nil {
		return x.Abilities
	}
	return nil
}

var file_role_proto_extTypes = []protoimpl.ExtensionInfo{
	{
		ExtendedType:  (*descriptorpb.MethodOptions)(nil),
		ExtensionType: (*RoleOptions)(nil),
		Field:         12102,
		Name:          "role.v1.roles",
		Tag:           "bytes,12102,opt,name=roles",
		Filename:      "role.proto",
	},
}

// Extension fields to descriptorpb.MethodOptions.
var (
	// optional role.v1.RoleOptions roles = 12102;
	E_Roles = &file_role_proto_extTypes[0]
)

var File_role_proto protoreflect.FileDescriptor

var file_role_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x72, 0x6f, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x72, 0x6f,
	0x6c, 0x65, 0x2e, 0x76, 0x31, 0x1a, 0x20, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f,
	0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x65, 0x0a, 0x0b, 0x52, 0x6f, 0x6c, 0x65, 0x4f,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x29, 0x0a, 0x08, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c,
	0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0e, 0x32, 0x0d, 0x2e, 0x72, 0x6f, 0x6c, 0x65, 0x2e,
	0x76, 0x31, 0x2e, 0x52, 0x6f, 0x6c, 0x65, 0x52, 0x08, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74,
	0x73, 0x12, 0x2b, 0x0a, 0x09, 0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x69, 0x65, 0x73, 0x18, 0x02,
	0x20, 0x03, 0x28, 0x0e, 0x32, 0x0d, 0x2e, 0x72, 0x6f, 0x6c, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x52,
	0x6f, 0x6c, 0x65, 0x52, 0x09, 0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x69, 0x65, 0x73, 0x2a, 0x2a,
	0x0a, 0x04, 0x52, 0x6f, 0x6c, 0x65, 0x12, 0x08, 0x0a, 0x04, 0x55, 0x53, 0x45, 0x52, 0x10, 0x00,
	0x12, 0x09, 0x0a, 0x05, 0x41, 0x44, 0x4d, 0x49, 0x4e, 0x10, 0x01, 0x12, 0x0d, 0x0a, 0x09, 0x41,
	0x4c, 0x4c, 0x5f, 0x52, 0x4f, 0x4c, 0x45, 0x53, 0x10, 0x64, 0x32, 0x0d, 0x0a, 0x0b, 0x52, 0x6f,
	0x6c, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x3a, 0x4b, 0x0a, 0x05, 0x72, 0x6f, 0x6c,
	0x65, 0x73, 0x12, 0x1e, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x18, 0xc6, 0x5e, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x72, 0x6f, 0x6c, 0x65,
	0x2e, 0x76, 0x31, 0x2e, 0x52, 0x6f, 0x6c, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52,
	0x05, 0x72, 0x6f, 0x6c, 0x65, 0x73, 0x42, 0x82, 0x01, 0x0a, 0x0b, 0x63, 0x6f, 0x6d, 0x2e, 0x72,
	0x6f, 0x6c, 0x65, 0x2e, 0x76, 0x31, 0x42, 0x09, 0x52, 0x6f, 0x6c, 0x65, 0x50, 0x72, 0x6f, 0x74,
	0x6f, 0x50, 0x01, 0x5a, 0x2b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x61, 0x6e, 0x68, 0x6e, 0x6d, 0x74, 0x2f, 0x67, 0x6f, 0x2d, 0x61, 0x70, 0x69, 0x2d, 0x62, 0x6f,
	0x69, 0x6c, 0x65, 0x72, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x70, 0x62,
	0xa2, 0x02, 0x03, 0x52, 0x58, 0x58, 0xaa, 0x02, 0x07, 0x52, 0x6f, 0x6c, 0x65, 0x2e, 0x56, 0x31,
	0xca, 0x02, 0x07, 0x52, 0x6f, 0x6c, 0x65, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x13, 0x52, 0x6f, 0x6c,
	0x65, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61,
	0xea, 0x02, 0x08, 0x52, 0x6f, 0x6c, 0x65, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_role_proto_rawDescOnce sync.Once
	file_role_proto_rawDescData = file_role_proto_rawDesc
)

func file_role_proto_rawDescGZIP() []byte {
	file_role_proto_rawDescOnce.Do(func() {
		file_role_proto_rawDescData = protoimpl.X.CompressGZIP(file_role_proto_rawDescData)
	})
	return file_role_proto_rawDescData
}

var file_role_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_role_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_role_proto_goTypes = []any{
	(Role)(0),                          // 0: role.v1.Role
	(*RoleOptions)(nil),                // 1: role.v1.RoleOptions
	(*descriptorpb.MethodOptions)(nil), // 2: google.protobuf.MethodOptions
}
var file_role_proto_depIdxs = []int32{
	0, // 0: role.v1.RoleOptions.defaults:type_name -> role.v1.Role
	0, // 1: role.v1.RoleOptions.abilities:type_name -> role.v1.Role
	2, // 2: role.v1.roles:extendee -> google.protobuf.MethodOptions
	1, // 3: role.v1.roles:type_name -> role.v1.RoleOptions
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	3, // [3:4] is the sub-list for extension type_name
	2, // [2:3] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_role_proto_init() }
func file_role_proto_init() {
	if File_role_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_role_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*RoleOptions); i {
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
			RawDescriptor: file_role_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 1,
			NumServices:   1,
		},
		GoTypes:           file_role_proto_goTypes,
		DependencyIndexes: file_role_proto_depIdxs,
		EnumInfos:         file_role_proto_enumTypes,
		MessageInfos:      file_role_proto_msgTypes,
		ExtensionInfos:    file_role_proto_extTypes,
	}.Build()
	File_role_proto = out.File
	file_role_proto_rawDesc = nil
	file_role_proto_goTypes = nil
	file_role_proto_depIdxs = nil
}