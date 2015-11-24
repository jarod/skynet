// Code generated by protoc-gen-go.
// source: skynet.proto
// DO NOT EDIT!

/*
Package skynet is a generated protocol buffer package.

It is generated from these files:
	skynet.proto

It has these top-level messages:
	Psint32
	Psint64
	Pstring
	AppServer
	AppMsg
*/
package skynet

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type SkynetMsg int32

const (
	SkynetMsg_SM_APP_INFO          SkynetMsg = 0
	SkynetMsg_SM_APP_DISCONNECTED  SkynetMsg = 1
	SkynetMsg_SM_SEND_TO_APP       SkynetMsg = 16
	SkynetMsg_SM_AGENT_EXECUTE_CMD SkynetMsg = 32
	SkynetMsg_SM_AGENT_FIND_APPS   SkynetMsg = 33
	SkynetMsg_SM_AGENT_PING        SkynetMsg = 34
)

var SkynetMsg_name = map[int32]string{
	0:  "SM_APP_INFO",
	1:  "SM_APP_DISCONNECTED",
	16: "SM_SEND_TO_APP",
	32: "SM_AGENT_EXECUTE_CMD",
	33: "SM_AGENT_FIND_APPS",
	34: "SM_AGENT_PING",
}
var SkynetMsg_value = map[string]int32{
	"SM_APP_INFO":          0,
	"SM_APP_DISCONNECTED":  1,
	"SM_SEND_TO_APP":       16,
	"SM_AGENT_EXECUTE_CMD": 32,
	"SM_AGENT_FIND_APPS":   33,
	"SM_AGENT_PING":        34,
}

func (x SkynetMsg) Enum() *SkynetMsg {
	p := new(SkynetMsg)
	*p = x
	return p
}
func (x SkynetMsg) String() string {
	return proto.EnumName(SkynetMsg_name, int32(x))
}
func (x *SkynetMsg) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(SkynetMsg_value, data, "SkynetMsg")
	if err != nil {
		return err
	}
	*x = SkynetMsg(value)
	return nil
}
func (SkynetMsg) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type Psint32 struct {
	Value            *int32 `protobuf:"zigzag32,1,opt,name=value" json:"value,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *Psint32) Reset()                    { *m = Psint32{} }
func (m *Psint32) String() string            { return proto.CompactTextString(m) }
func (*Psint32) ProtoMessage()               {}
func (*Psint32) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Psint32) GetValue() int32 {
	if m != nil && m.Value != nil {
		return *m.Value
	}
	return 0
}

type Psint64 struct {
	Value            *int64 `protobuf:"zigzag64,2,opt,name=value" json:"value,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *Psint64) Reset()                    { *m = Psint64{} }
func (m *Psint64) String() string            { return proto.CompactTextString(m) }
func (*Psint64) ProtoMessage()               {}
func (*Psint64) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Psint64) GetValue() int64 {
	if m != nil && m.Value != nil {
		return *m.Value
	}
	return 0
}

type Pstring struct {
	Value            *string `protobuf:"bytes,3,opt,name=value" json:"value,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *Pstring) Reset()                    { *m = Pstring{} }
func (m *Pstring) String() string            { return proto.CompactTextString(m) }
func (*Pstring) ProtoMessage()               {}
func (*Pstring) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *Pstring) GetValue() string {
	if m != nil && m.Value != nil {
		return *m.Value
	}
	return ""
}

type AppServer struct {
	Id               *string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Host             *string `protobuf:"bytes,2,opt,name=host" json:"host,omitempty"`
	Port             *uint32 `protobuf:"varint,3,opt,name=port" json:"port,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *AppServer) Reset()                    { *m = AppServer{} }
func (m *AppServer) String() string            { return proto.CompactTextString(m) }
func (*AppServer) ProtoMessage()               {}
func (*AppServer) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *AppServer) GetId() string {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return ""
}

func (m *AppServer) GetHost() string {
	if m != nil && m.Host != nil {
		return *m.Host
	}
	return ""
}

func (m *AppServer) GetPort() uint32 {
	if m != nil && m.Port != nil {
		return *m.Port
	}
	return 0
}

type AppMsg struct {
	AppId            *string `protobuf:"bytes,1,req,name=app_id" json:"app_id,omitempty"`
	Head             *uint32 `protobuf:"varint,2,req,name=head" json:"head,omitempty"`
	Data             []byte  `protobuf:"bytes,3,opt,name=data" json:"data,omitempty"`
	FromApp          *string `protobuf:"bytes,4,opt,name=from_app" json:"from_app,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *AppMsg) Reset()                    { *m = AppMsg{} }
func (m *AppMsg) String() string            { return proto.CompactTextString(m) }
func (*AppMsg) ProtoMessage()               {}
func (*AppMsg) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *AppMsg) GetAppId() string {
	if m != nil && m.AppId != nil {
		return *m.AppId
	}
	return ""
}

func (m *AppMsg) GetHead() uint32 {
	if m != nil && m.Head != nil {
		return *m.Head
	}
	return 0
}

func (m *AppMsg) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *AppMsg) GetFromApp() string {
	if m != nil && m.FromApp != nil {
		return *m.FromApp
	}
	return ""
}

func init() {
	proto.RegisterType((*Psint32)(nil), "skynet.Psint32")
	proto.RegisterType((*Psint64)(nil), "skynet.Psint64")
	proto.RegisterType((*Pstring)(nil), "skynet.Pstring")
	proto.RegisterType((*AppServer)(nil), "skynet.AppServer")
	proto.RegisterType((*AppMsg)(nil), "skynet.AppMsg")
	proto.RegisterEnum("skynet.SkynetMsg", SkynetMsg_name, SkynetMsg_value)
}

var fileDescriptor0 = []byte{
	// 287 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x54, 0x90, 0x41, 0x4f, 0xc2, 0x30,
	0x18, 0x86, 0xdd, 0x44, 0x74, 0x9f, 0x0c, 0xcb, 0xa7, 0xd1, 0x1d, 0x71, 0x27, 0xe2, 0xc1, 0x83,
	0x18, 0xef, 0xb8, 0x15, 0xb2, 0x03, 0x65, 0x49, 0x67, 0xe2, 0xad, 0x59, 0xc2, 0xc4, 0x45, 0xdd,
	0x9a, 0xae, 0x92, 0xf8, 0x27, 0xfc, 0xcd, 0x76, 0x85, 0x8c, 0x78, 0x7c, 0xbf, 0xe7, 0x7d, 0xde,
	0xa4, 0x85, 0x41, 0xf3, 0xf1, 0x53, 0x15, 0xfa, 0x5e, 0xaa, 0x5a, 0xd7, 0xd8, 0xdf, 0xa5, 0x30,
	0x80, 0xd3, 0xb4, 0x29, 0x2b, 0x3d, 0x7d, 0x40, 0x1f, 0x4e, 0xb6, 0xf9, 0xe7, 0x77, 0x11, 0x38,
	0x63, 0x67, 0x32, 0xea, 0xc8, 0xd3, 0xe3, 0x81, 0xb8, 0x86, 0xe0, 0x8e, 0x68, 0x55, 0x56, 0x9b,
	0x03, 0x39, 0x36, 0xc4, 0x0b, 0xa7, 0xe0, 0xcd, 0xa4, 0xe4, 0x85, 0xda, 0x16, 0x0a, 0x01, 0xdc,
	0x72, 0x6d, 0xc7, 0x3c, 0x1c, 0x40, 0xef, 0xbd, 0x6e, 0xb4, 0x1d, 0xb0, 0x49, 0xd6, 0x4a, 0x5b,
	0xc9, 0x0f, 0xe7, 0xd0, 0x37, 0xd2, 0xb2, 0xd9, 0xe0, 0x10, 0xfa, 0xb9, 0x94, 0xc2, 0x5a, 0xee,
	0xde, 0x2a, 0xf2, 0xb5, 0xb1, 0xdc, 0x89, 0xdf, 0xa6, 0x75, 0xae, 0x73, 0x6b, 0x0d, 0x90, 0xc0,
	0xd9, 0x9b, 0xaa, 0xbf, 0x84, 0x11, 0x82, 0x5e, 0xbb, 0x7a, 0xf7, 0xeb, 0x80, 0xc7, 0xed, 0xab,
	0xda, 0xad, 0x0b, 0x38, 0xe7, 0x4b, 0x31, 0x4b, 0x53, 0x91, 0xb0, 0xf9, 0x8a, 0x1c, 0xe1, 0x0d,
	0x5c, 0xee, 0x0f, 0x71, 0xc2, 0xa3, 0x15, 0x63, 0x34, 0xca, 0x68, 0x4c, 0x1c, 0x44, 0x18, 0x1a,
	0xc0, 0x29, 0x8b, 0x45, 0xb6, 0x6a, 0x0b, 0x84, 0x60, 0x00, 0x57, 0x6d, 0x79, 0x41, 0x59, 0x26,
	0xe8, 0x2b, 0x8d, 0x5e, 0x32, 0x2a, 0xa2, 0x65, 0x4c, 0xc6, 0x78, 0x0d, 0xd8, 0x91, 0x79, 0x62,
	0x1c, 0x23, 0x70, 0x72, 0x8b, 0x23, 0xf0, 0xbb, 0x7b, 0x9a, 0xb0, 0x05, 0x09, 0x9f, 0x87, 0xff,
	0xff, 0xfc, 0x2f, 0x00, 0x00, 0xff, 0xff, 0x9e, 0x13, 0x42, 0x39, 0x82, 0x01, 0x00, 0x00,
}
