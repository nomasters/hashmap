// Code generated by protoc-gen-go. DO NOT EDIT.
// source: payload.proto

package pb

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	duration "github.com/golang/protobuf/ptypes/duration"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
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

type Payload_Version int32

const (
	Payload_V0 Payload_Version = 0
	Payload_V1 Payload_Version = 1
)

var Payload_Version_name = map[int32]string{
	0: "V0",
	1: "V1",
}

var Payload_Version_value = map[string]int32{
	"V0": 0,
	"V1": 1,
}

func (x Payload_Version) String() string {
	return proto.EnumName(Payload_Version_name, int32(x))
}

func (Payload_Version) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_678c914f1bee6d56, []int{0, 0}
}

type Payload_SigBundle_Alg int32

const (
	Payload_SigBundle_UNKNOWN   Payload_SigBundle_Alg = 0
	Payload_SigBundle_NACL_SIGN Payload_SigBundle_Alg = 1
)

var Payload_SigBundle_Alg_name = map[int32]string{
	0: "UNKNOWN",
	1: "NACL_SIGN",
}

var Payload_SigBundle_Alg_value = map[string]int32{
	"UNKNOWN":   0,
	"NACL_SIGN": 1,
}

func (x Payload_SigBundle_Alg) String() string {
	return proto.EnumName(Payload_SigBundle_Alg_name, int32(x))
}

func (Payload_SigBundle_Alg) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_678c914f1bee6d56, []int{0, 0, 0}
}

type Payload struct {
	Version              Payload_Version      `protobuf:"varint,1,opt,name=version,proto3,enum=pb.Payload_Version" json:"version,omitempty"`
	Timestamp            *timestamp.Timestamp `protobuf:"bytes,2,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Ttl                  *duration.Duration   `protobuf:"bytes,3,opt,name=ttl,proto3" json:"ttl,omitempty"`
	SigBundles           []*Payload_SigBundle `protobuf:"bytes,4,rep,name=sig_bundles,json=sigBundles,proto3" json:"sig_bundles,omitempty"`
	Len                  uint32               `protobuf:"varint,5,opt,name=len,proto3" json:"len,omitempty"`
	Data                 []byte               `protobuf:"bytes,6,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *Payload) Reset()         { *m = Payload{} }
func (m *Payload) String() string { return proto.CompactTextString(m) }
func (*Payload) ProtoMessage()    {}
func (*Payload) Descriptor() ([]byte, []int) {
	return fileDescriptor_678c914f1bee6d56, []int{0}
}

func (m *Payload) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Payload.Unmarshal(m, b)
}
func (m *Payload) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Payload.Marshal(b, m, deterministic)
}
func (m *Payload) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Payload.Merge(m, src)
}
func (m *Payload) XXX_Size() int {
	return xxx_messageInfo_Payload.Size(m)
}
func (m *Payload) XXX_DiscardUnknown() {
	xxx_messageInfo_Payload.DiscardUnknown(m)
}

var xxx_messageInfo_Payload proto.InternalMessageInfo

func (m *Payload) GetVersion() Payload_Version {
	if m != nil {
		return m.Version
	}
	return Payload_V0
}

func (m *Payload) GetTimestamp() *timestamp.Timestamp {
	if m != nil {
		return m.Timestamp
	}
	return nil
}

func (m *Payload) GetTtl() *duration.Duration {
	if m != nil {
		return m.Ttl
	}
	return nil
}

func (m *Payload) GetSigBundles() []*Payload_SigBundle {
	if m != nil {
		return m.SigBundles
	}
	return nil
}

func (m *Payload) GetLen() uint32 {
	if m != nil {
		return m.Len
	}
	return 0
}

func (m *Payload) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

type Payload_SigBundle struct {
	Alg                  Payload_SigBundle_Alg `protobuf:"varint,1,opt,name=alg,proto3,enum=pb.Payload_SigBundle_Alg" json:"alg,omitempty"`
	Pub                  []byte                `protobuf:"bytes,2,opt,name=pub,proto3" json:"pub,omitempty"`
	Sig                  []byte                `protobuf:"bytes,3,opt,name=sig,proto3" json:"sig,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *Payload_SigBundle) Reset()         { *m = Payload_SigBundle{} }
func (m *Payload_SigBundle) String() string { return proto.CompactTextString(m) }
func (*Payload_SigBundle) ProtoMessage()    {}
func (*Payload_SigBundle) Descriptor() ([]byte, []int) {
	return fileDescriptor_678c914f1bee6d56, []int{0, 0}
}

func (m *Payload_SigBundle) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Payload_SigBundle.Unmarshal(m, b)
}
func (m *Payload_SigBundle) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Payload_SigBundle.Marshal(b, m, deterministic)
}
func (m *Payload_SigBundle) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Payload_SigBundle.Merge(m, src)
}
func (m *Payload_SigBundle) XXX_Size() int {
	return xxx_messageInfo_Payload_SigBundle.Size(m)
}
func (m *Payload_SigBundle) XXX_DiscardUnknown() {
	xxx_messageInfo_Payload_SigBundle.DiscardUnknown(m)
}

var xxx_messageInfo_Payload_SigBundle proto.InternalMessageInfo

func (m *Payload_SigBundle) GetAlg() Payload_SigBundle_Alg {
	if m != nil {
		return m.Alg
	}
	return Payload_SigBundle_UNKNOWN
}

func (m *Payload_SigBundle) GetPub() []byte {
	if m != nil {
		return m.Pub
	}
	return nil
}

func (m *Payload_SigBundle) GetSig() []byte {
	if m != nil {
		return m.Sig
	}
	return nil
}

func init() {
	proto.RegisterEnum("pb.Payload_Version", Payload_Version_name, Payload_Version_value)
	proto.RegisterEnum("pb.Payload_SigBundle_Alg", Payload_SigBundle_Alg_name, Payload_SigBundle_Alg_value)
	proto.RegisterType((*Payload)(nil), "pb.Payload")
	proto.RegisterType((*Payload_SigBundle)(nil), "pb.Payload.SigBundle")
}

func init() { proto.RegisterFile("payload.proto", fileDescriptor_678c914f1bee6d56) }

var fileDescriptor_678c914f1bee6d56 = []byte{
	// 332 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x90, 0x51, 0x4b, 0xfb, 0x30,
	0x14, 0xc5, 0xd7, 0x75, 0xff, 0x95, 0xdd, 0x6e, 0x7f, 0x4a, 0x44, 0xe8, 0xfa, 0xa0, 0x75, 0x4f,
	0x05, 0x31, 0xd3, 0x09, 0xe2, 0xeb, 0x54, 0x10, 0x51, 0xaa, 0x64, 0x3a, 0x1f, 0x47, 0x4a, 0x6b,
	0x28, 0x64, 0x4d, 0x59, 0x52, 0xc1, 0x27, 0x3f, 0x87, 0xdf, 0x56, 0x92, 0xb6, 0x53, 0xd4, 0xa7,
	0x1c, 0x72, 0x7e, 0xf7, 0x9e, 0x7b, 0x2f, 0x8c, 0x4a, 0xfa, 0xc6, 0x05, 0x4d, 0x71, 0xb9, 0x11,
	0x4a, 0xa0, 0x6e, 0x99, 0x04, 0xfb, 0x4c, 0x08, 0xc6, 0xb3, 0xa9, 0xf9, 0x49, 0xaa, 0x97, 0xa9,
	0xca, 0xd7, 0x99, 0x54, 0x74, 0x5d, 0xd6, 0x50, 0xb0, 0xf7, 0x13, 0x48, 0xab, 0x0d, 0x55, 0xb9,
	0x28, 0x6a, 0x7f, 0xf2, 0x61, 0x83, 0xf3, 0x50, 0xb7, 0x45, 0x47, 0xe0, 0xbc, 0x66, 0x1b, 0x99,
	0x8b, 0xc2, 0xb7, 0x42, 0x2b, 0xfa, 0x3f, 0xdb, 0xc1, 0x65, 0x82, 0x1b, 0x17, 0x2f, 0x6b, 0x8b,
	0xb4, 0x0c, 0x3a, 0x87, 0xc1, 0x36, 0xcd, 0xef, 0x86, 0x56, 0xe4, 0xce, 0x02, 0x5c, 0xc7, 0xe1,
	0x36, 0x0e, 0x3f, 0xb6, 0x04, 0xf9, 0x82, 0xd1, 0x21, 0xd8, 0x4a, 0x71, 0xdf, 0x36, 0x35, 0xe3,
	0x5f, 0x35, 0x57, 0xcd, 0x88, 0x44, 0x53, 0xe8, 0x0c, 0x5c, 0x99, 0xb3, 0x55, 0x52, 0x15, 0x29,
	0xcf, 0xa4, 0xdf, 0x0b, 0xed, 0xc8, 0x9d, 0xed, 0x7e, 0x9f, 0x6c, 0x91, 0xb3, 0x0b, 0xe3, 0x12,
	0x90, 0xad, 0x94, 0xc8, 0x03, 0x9b, 0x67, 0x85, 0xff, 0x2f, 0xb4, 0xa2, 0x11, 0xd1, 0x12, 0x21,
	0xe8, 0xa5, 0x54, 0x51, 0xbf, 0x1f, 0x5a, 0xd1, 0x90, 0x18, 0x1d, 0xbc, 0xc3, 0x60, 0x5b, 0xae,
	0xe7, 0xa2, 0x9c, 0x35, 0xcb, 0x8f, 0xff, 0x8c, 0xc0, 0x73, 0xce, 0x88, 0xa6, 0x74, 0xff, 0xb2,
	0x4a, 0xcc, 0xe2, 0x43, 0xa2, 0xa5, 0xfe, 0x91, 0x39, 0x33, 0x6b, 0x0d, 0x89, 0x96, 0x93, 0x03,
	0xb0, 0xe7, 0x9c, 0x21, 0x17, 0x9c, 0xa7, 0xf8, 0x36, 0xbe, 0x7f, 0x8e, 0xbd, 0x0e, 0x1a, 0xc1,
	0x20, 0x9e, 0x5f, 0xde, 0xad, 0x16, 0x37, 0xd7, 0xb1, 0x67, 0x4d, 0xc6, 0xe0, 0x34, 0x97, 0x45,
	0x7d, 0xe8, 0x2e, 0x8f, 0xbd, 0x8e, 0x79, 0x4f, 0x3c, 0x2b, 0xe9, 0x9b, 0x8b, 0x9c, 0x7e, 0x06,
	0x00, 0x00, 0xff, 0xff, 0x03, 0x01, 0xc1, 0x33, 0xf8, 0x01, 0x00, 0x00,
}