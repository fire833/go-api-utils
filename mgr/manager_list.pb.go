/*
*	Copyright (C) 2024 Kendall Tauser
*
*	This program is free software; you can redistribute it and/or modify
*	it under the terms of the GNU General Public License as published by
*	the Free Software Foundation; either version 2 of the License, or
*	(at your option) any later version.
*
*	This program is distributed in the hope that it will be useful,
*	but WITHOUT ANY WARRANTY; without even the implied warranty of
*	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
*	GNU General Public License for more details.
*
*	You should have received a copy of the GNU General Public License along
*	with this program; if not, write to the Free Software Foundation, Inc.,
*	51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.
 */

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.5
// source: manager_list.proto

package mgr

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type SubsystemStatusList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Items []*SubsystemStatus `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
}

func (x *SubsystemStatusList) Reset() {
	*x = SubsystemStatusList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_manager_list_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SubsystemStatusList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubsystemStatusList) ProtoMessage() {}

func (x *SubsystemStatusList) ProtoReflect() protoreflect.Message {
	mi := &file_manager_list_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubsystemStatusList.ProtoReflect.Descriptor instead.
func (*SubsystemStatusList) Descriptor() ([]byte, []int) {
	return file_manager_list_proto_rawDescGZIP(), []int{0}
}

func (x *SubsystemStatusList) GetItems() []*SubsystemStatus {
	if x != nil {
		return x.Items
	}
	return nil
}

type BuildInfoList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Items []*BuildInfo `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
}

func (x *BuildInfoList) Reset() {
	*x = BuildInfoList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_manager_list_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BuildInfoList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BuildInfoList) ProtoMessage() {}

func (x *BuildInfoList) ProtoReflect() protoreflect.Message {
	mi := &file_manager_list_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BuildInfoList.ProtoReflect.Descriptor instead.
func (*BuildInfoList) Descriptor() ([]byte, []int) {
	return file_manager_list_proto_rawDescGZIP(), []int{1}
}

func (x *BuildInfoList) GetItems() []*BuildInfo {
	if x != nil {
		return x.Items
	}
	return nil
}

var File_manager_list_proto protoreflect.FileDescriptor

var file_manager_list_proto_rawDesc = []byte{
	0x0a, 0x12, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x5f, 0x6c, 0x69, 0x73, 0x74, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x1a, 0x0d, 0x6d,
	0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x45, 0x0a, 0x13,
	0x53, 0x75, 0x62, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x4c,
	0x69, 0x73, 0x74, 0x12, 0x2e, 0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x18, 0x2e, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x53, 0x75, 0x62,
	0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x05, 0x69, 0x74,
	0x65, 0x6d, 0x73, 0x22, 0x39, 0x0a, 0x0d, 0x42, 0x75, 0x69, 0x6c, 0x64, 0x49, 0x6e, 0x66, 0x6f,
	0x4c, 0x69, 0x73, 0x74, 0x12, 0x28, 0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x42, 0x75,
	0x69, 0x6c, 0x64, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x42, 0x0c,
	0x5a, 0x0a, 0x2e, 0x2e, 0x2f, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_manager_list_proto_rawDescOnce sync.Once
	file_manager_list_proto_rawDescData = file_manager_list_proto_rawDesc
)

func file_manager_list_proto_rawDescGZIP() []byte {
	file_manager_list_proto_rawDescOnce.Do(func() {
		file_manager_list_proto_rawDescData = protoimpl.X.CompressGZIP(file_manager_list_proto_rawDescData)
	})
	return file_manager_list_proto_rawDescData
}

var (
	file_manager_list_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
	file_manager_list_proto_goTypes  = []interface{}{
		(*SubsystemStatusList)(nil), // 0: manager.SubsystemStatusList
		(*BuildInfoList)(nil),       // 1: manager.BuildInfoList
		(*SubsystemStatus)(nil),     // 2: manager.SubsystemStatus
		(*BuildInfo)(nil),           // 3: manager.BuildInfo
	}
)

var file_manager_list_proto_depIdxs = []int32{
	2, // 0: manager.SubsystemStatusList.items:type_name -> manager.SubsystemStatus
	3, // 1: manager.BuildInfoList.items:type_name -> manager.BuildInfo
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_manager_list_proto_init() }
func file_manager_list_proto_init() {
	if File_manager_list_proto != nil {
		return
	}
	file_manager_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_manager_list_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SubsystemStatusList); i {
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
		file_manager_list_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BuildInfoList); i {
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
			RawDescriptor: file_manager_list_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_manager_list_proto_goTypes,
		DependencyIndexes: file_manager_list_proto_depIdxs,
		MessageInfos:      file_manager_list_proto_msgTypes,
	}.Build()
	File_manager_list_proto = out.File
	file_manager_list_proto_rawDesc = nil
	file_manager_list_proto_goTypes = nil
	file_manager_list_proto_depIdxs = nil
}
