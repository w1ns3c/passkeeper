// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.12.4
// source: users_change_pwd.proto

package protofiles

import (
	empty "github.com/golang/protobuf/ptypes/empty"
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

type UserChangePassReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OldPass string `protobuf:"bytes,1,opt,name=OldPass,proto3" json:"OldPass,omitempty"`
	NewPass string `protobuf:"bytes,2,opt,name=NewPass,proto3" json:"NewPass,omitempty"`
	Repeat  string `protobuf:"bytes,3,opt,name=Repeat,proto3" json:"Repeat,omitempty"`
}

func (x *UserChangePassReq) Reset() {
	*x = UserChangePassReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_users_change_pwd_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserChangePassReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserChangePassReq) ProtoMessage() {}

func (x *UserChangePassReq) ProtoReflect() protoreflect.Message {
	mi := &file_users_change_pwd_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserChangePassReq.ProtoReflect.Descriptor instead.
func (*UserChangePassReq) Descriptor() ([]byte, []int) {
	return file_users_change_pwd_proto_rawDescGZIP(), []int{0}
}

func (x *UserChangePassReq) GetOldPass() string {
	if x != nil {
		return x.OldPass
	}
	return ""
}

func (x *UserChangePassReq) GetNewPass() string {
	if x != nil {
		return x.NewPass
	}
	return ""
}

func (x *UserChangePassReq) GetRepeat() string {
	if x != nil {
		return x.Repeat
	}
	return ""
}

var File_users_change_pwd_proto protoreflect.FileDescriptor

var file_users_change_pwd_proto_rawDesc = []byte{
	0x0a, 0x16, 0x75, 0x73, 0x65, 0x72, 0x73, 0x5f, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x5f, 0x70,
	0x77, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x5f, 0x0a, 0x11, 0x55, 0x73, 0x65, 0x72, 0x43, 0x68, 0x61,
	0x6e, 0x67, 0x65, 0x50, 0x61, 0x73, 0x73, 0x52, 0x65, 0x71, 0x12, 0x18, 0x0a, 0x07, 0x4f, 0x6c,
	0x64, 0x50, 0x61, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x4f, 0x6c, 0x64,
	0x50, 0x61, 0x73, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x4e, 0x65, 0x77, 0x50, 0x61, 0x73, 0x73, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x4e, 0x65, 0x77, 0x50, 0x61, 0x73, 0x73, 0x12, 0x16,
	0x0a, 0x06, 0x52, 0x65, 0x70, 0x65, 0x61, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x52, 0x65, 0x70, 0x65, 0x61, 0x74, 0x32, 0x47, 0x0a, 0x0b, 0x55, 0x73, 0x65, 0x72, 0x50, 0x61,
	0x73, 0x73, 0x53, 0x76, 0x63, 0x12, 0x38, 0x0a, 0x0a, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x50,
	0x61, 0x73, 0x73, 0x12, 0x12, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65,
	0x50, 0x61, 0x73, 0x73, 0x52, 0x65, 0x71, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x42,
	0x41, 0x5a, 0x3f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x77, 0x31,
	0x6e, 0x73, 0x33, 0x63, 0x2f, 0x70, 0x61, 0x73, 0x73, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x2f,
	0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f,
	0x72, 0x74, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x66, 0x69, 0x6c,
	0x65, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_users_change_pwd_proto_rawDescOnce sync.Once
	file_users_change_pwd_proto_rawDescData = file_users_change_pwd_proto_rawDesc
)

func file_users_change_pwd_proto_rawDescGZIP() []byte {
	file_users_change_pwd_proto_rawDescOnce.Do(func() {
		file_users_change_pwd_proto_rawDescData = protoimpl.X.CompressGZIP(file_users_change_pwd_proto_rawDescData)
	})
	return file_users_change_pwd_proto_rawDescData
}

var file_users_change_pwd_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_users_change_pwd_proto_goTypes = []interface{}{
	(*UserChangePassReq)(nil), // 0: UserChangePassReq
	(*empty.Empty)(nil),       // 1: google.protobuf.Empty
}
var file_users_change_pwd_proto_depIdxs = []int32{
	0, // 0: UserPassSvc.ChangePass:input_type -> UserChangePassReq
	1, // 1: UserPassSvc.ChangePass:output_type -> google.protobuf.Empty
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_users_change_pwd_proto_init() }
func file_users_change_pwd_proto_init() {
	if File_users_change_pwd_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_users_change_pwd_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserChangePassReq); i {
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
			RawDescriptor: file_users_change_pwd_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_users_change_pwd_proto_goTypes,
		DependencyIndexes: file_users_change_pwd_proto_depIdxs,
		MessageInfos:      file_users_change_pwd_proto_msgTypes,
	}.Build()
	File_users_change_pwd_proto = out.File
	file_users_change_pwd_proto_rawDesc = nil
	file_users_change_pwd_proto_goTypes = nil
	file_users_change_pwd_proto_depIdxs = nil
}
