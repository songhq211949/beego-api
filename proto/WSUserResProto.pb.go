// Code generated by protoc-gen-go. DO NOT EDIT.
// source: WSUserResProto.proto

package protocol

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// 请求实体
type WSUserResProto struct {
	Uid                  uint64   `protobuf:"varint,1,opt,name=uid,proto3" json:"uid,omitempty"`
	Name                 string   `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Avatar               string   `protobuf:"bytes,3,opt,name=avatar,proto3" json:"avatar,omitempty"`
	Remark               string   `protobuf:"bytes,4,opt,name=remark,proto3" json:"remark,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *WSUserResProto) Reset()         { *m = WSUserResProto{} }
func (m *WSUserResProto) String() string { return proto.CompactTextString(m) }
func (*WSUserResProto) ProtoMessage()    {}
func (*WSUserResProto) Descriptor() ([]byte, []int) {
	return fileDescriptor_ccb052403a7075aa, []int{0}
}

func (m *WSUserResProto) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WSUserResProto.Unmarshal(m, b)
}
func (m *WSUserResProto) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WSUserResProto.Marshal(b, m, deterministic)
}
func (m *WSUserResProto) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WSUserResProto.Merge(m, src)
}
func (m *WSUserResProto) XXX_Size() int {
	return xxx_messageInfo_WSUserResProto.Size(m)
}
func (m *WSUserResProto) XXX_DiscardUnknown() {
	xxx_messageInfo_WSUserResProto.DiscardUnknown(m)
}

var xxx_messageInfo_WSUserResProto proto.InternalMessageInfo

func (m *WSUserResProto) GetUid() uint64 {
	if m != nil {
		return m.Uid
	}
	return 0
}

func (m *WSUserResProto) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *WSUserResProto) GetAvatar() string {
	if m != nil {
		return m.Avatar
	}
	return ""
}

func (m *WSUserResProto) GetRemark() string {
	if m != nil {
		return m.Remark
	}
	return ""
}

func init() {
	proto.RegisterType((*WSUserResProto)(nil), "protocol.WSUserResProto")
}

func init() {
	proto.RegisterFile("WSUserResProto.proto", fileDescriptor_ccb052403a7075aa)
}

var fileDescriptor_ccb052403a7075aa = []byte{
	// 127 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x09, 0x0f, 0x0e, 0x2d,
	0x4e, 0x2d, 0x0a, 0x4a, 0x2d, 0x0e, 0x28, 0xca, 0x2f, 0xc9, 0xd7, 0x2b, 0x00, 0x91, 0x42, 0x1c,
	0x60, 0x2a, 0x39, 0x3f, 0x47, 0x29, 0x8d, 0x8b, 0x0f, 0x55, 0x85, 0x90, 0x00, 0x17, 0x73, 0x69,
	0x66, 0x8a, 0x04, 0xa3, 0x02, 0xa3, 0x06, 0x4b, 0x10, 0x88, 0x29, 0x24, 0xc4, 0xc5, 0x92, 0x97,
	0x98, 0x9b, 0x2a, 0xc1, 0xa4, 0xc0, 0xa8, 0xc1, 0x19, 0x04, 0x66, 0x0b, 0x89, 0x71, 0xb1, 0x25,
	0x96, 0x25, 0x96, 0x24, 0x16, 0x49, 0x30, 0x83, 0x45, 0xa1, 0x3c, 0x90, 0x78, 0x51, 0x6a, 0x6e,
	0x62, 0x51, 0xb6, 0x04, 0x0b, 0x44, 0x1c, 0xc2, 0x73, 0x62, 0xf2, 0x60, 0x4c, 0x62, 0x03, 0xdb,
	0x6a, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff, 0x20, 0x83, 0x4c, 0x17, 0x94, 0x00, 0x00, 0x00,
}