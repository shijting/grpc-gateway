// Code generated by protoc-gen-go. DO NOT EDIT.
// source: enum.proto

package proto

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

// 运营商
type Operator int32

const (
	Operator_OPERATOR_NONE    Operator = 0
	Operator_OPERATOR_UNICOM  Operator = 1
	Operator_OPERATOR_MOBILE  Operator = 2
	Operator_OPERATOR_TELECOM Operator = 3
)

var Operator_name = map[int32]string{
	0: "OPERATOR_NONE",
	1: "OPERATOR_UNICOM",
	2: "OPERATOR_MOBILE",
	3: "OPERATOR_TELECOM",
}

var Operator_value = map[string]int32{
	"OPERATOR_NONE":    0,
	"OPERATOR_UNICOM":  1,
	"OPERATOR_MOBILE":  2,
	"OPERATOR_TELECOM": 3,
}

func (x Operator) String() string {
	return proto.EnumName(Operator_name, int32(x))
}

func (Operator) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_13a9f1b5947140c8, []int{0}
}

// 设备状态
type Status int32

const (
	Status_STATUS_NONE        Status = 0
	Status_STATUS_TEST        Status = 1
	Status_STATUS_INVENTORY   Status = 2
	Status_STATUS_READY       Status = 3
	Status_STATUS_ACTIVATED   Status = 4
	Status_STATUS_DEACTIVATED Status = 5
	Status_STATUS_RETIRED     Status = 6
)

var Status_name = map[int32]string{
	0: "STATUS_NONE",
	1: "STATUS_TEST",
	2: "STATUS_INVENTORY",
	3: "STATUS_READY",
	4: "STATUS_ACTIVATED",
	5: "STATUS_DEACTIVATED",
	6: "STATUS_RETIRED",
}

var Status_value = map[string]int32{
	"STATUS_NONE":        0,
	"STATUS_TEST":        1,
	"STATUS_INVENTORY":   2,
	"STATUS_READY":       3,
	"STATUS_ACTIVATED":   4,
	"STATUS_DEACTIVATED": 5,
	"STATUS_RETIRED":     6,
}

func (x Status) String() string {
	return proto.EnumName(Status_name, int32(x))
}

func (Status) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_13a9f1b5947140c8, []int{1}
}

// 服务状态
type ServiceStatus int32

const (
	ServiceStatus_SERVICE_NONE     ServiceStatus = 0
	ServiceStatus_SERVICE_ENABLED  ServiceStatus = 1
	ServiceStatus_SERVICE_DISABLED ServiceStatus = 2
)

var ServiceStatus_name = map[int32]string{
	0: "SERVICE_NONE",
	1: "SERVICE_ENABLED",
	2: "SERVICE_DISABLED",
}

var ServiceStatus_value = map[string]int32{
	"SERVICE_NONE":     0,
	"SERVICE_ENABLED":  1,
	"SERVICE_DISABLED": 2,
}

func (x ServiceStatus) String() string {
	return proto.EnumName(ServiceStatus_name, int32(x))
}

func (ServiceStatus) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_13a9f1b5947140c8, []int{2}
}

// 对设备操作状态
type SwitchStatus int32

const (
	SwitchStatus_SWITCH_NONE        SwitchStatus = 0
	SwitchStatus_SWITCH_AVAILABLE   SwitchStatus = 1
	SwitchStatus_SWITCH_ACTIVATE    SwitchStatus = 2
	SwitchStatus_SWITCH_DEACTIVATE  SwitchStatus = 4
	SwitchStatus_SWITCH_OPEN_DATA   SwitchStatus = 8
	SwitchStatus_SWITCH_CLOSE_DATA  SwitchStatus = 16
	SwitchStatus_SWITCH_OPEN_VOICE  SwitchStatus = 32
	SwitchStatus_SWITCH_CLOSE_VOICE SwitchStatus = 64
)

var SwitchStatus_name = map[int32]string{
	0:  "SWITCH_NONE",
	1:  "SWITCH_AVAILABLE",
	2:  "SWITCH_ACTIVATE",
	4:  "SWITCH_DEACTIVATE",
	8:  "SWITCH_OPEN_DATA",
	16: "SWITCH_CLOSE_DATA",
	32: "SWITCH_OPEN_VOICE",
	64: "SWITCH_CLOSE_VOICE",
}

var SwitchStatus_value = map[string]int32{
	"SWITCH_NONE":        0,
	"SWITCH_AVAILABLE":   1,
	"SWITCH_ACTIVATE":    2,
	"SWITCH_DEACTIVATE":  4,
	"SWITCH_OPEN_DATA":   8,
	"SWITCH_CLOSE_DATA":  16,
	"SWITCH_OPEN_VOICE":  32,
	"SWITCH_CLOSE_VOICE": 64,
}

func (x SwitchStatus) String() string {
	return proto.EnumName(SwitchStatus_name, int32(x))
}

func (SwitchStatus) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_13a9f1b5947140c8, []int{3}
}

// 设备限速等级
type SpeedLevel int32

const (
	SpeedLevel_SPEED_NONE  SpeedLevel = 0
	SpeedLevel_SPEED_512KB SpeedLevel = 1
	SpeedLevel_SPEED_1MB   SpeedLevel = 2
	SpeedLevel_SPEED_2MB   SpeedLevel = 3
	SpeedLevel_SPEED_7MB   SpeedLevel = 4
)

var SpeedLevel_name = map[int32]string{
	0: "SPEED_NONE",
	1: "SPEED_512KB",
	2: "SPEED_1MB",
	3: "SPEED_2MB",
	4: "SPEED_7MB",
}

var SpeedLevel_value = map[string]int32{
	"SPEED_NONE":  0,
	"SPEED_512KB": 1,
	"SPEED_1MB":   2,
	"SPEED_2MB":   3,
	"SPEED_7MB":   4,
}

func (x SpeedLevel) String() string {
	return proto.EnumName(SpeedLevel_name, int32(x))
}

func (SpeedLevel) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_13a9f1b5947140c8, []int{4}
}

type Error_Code int32

const (
	Error_ERR_OK                 Error_Code = 0
	Error_ERR_NOT_FOUND          Error_Code = 1
	Error_ERR_METHOD_NOT_ALLOWED Error_Code = 2
	Error_ERR_INTERNAL_SERVER    Error_Code = 3
	Error_ERR_UNAUTHORIZED       Error_Code = 4
	Error_ERR_FORBIDDEN          Error_Code = 5
	Error_ERR_RATE_LIMIT         Error_Code = 6
	Error_ERR_ILLEGAL_ACCESS     Error_Code = 7
	Error_ERR_DATABASE           Error_Code = 8
	Error_ERR_INVALID_ARGS       Error_Code = 9
	Error_ERR_INVALID_METHOD     Error_Code = 10
	Error_ERR_OPERATOR           Error_Code = 11
	Error_ERR_TASK_IN_PROGRESS   Error_Code = 12
	Error_ERR_INVALID_DEVICE     Error_Code = 13
)

var Error_Code_name = map[int32]string{
	0:  "ERR_OK",
	1:  "ERR_NOT_FOUND",
	2:  "ERR_METHOD_NOT_ALLOWED",
	3:  "ERR_INTERNAL_SERVER",
	4:  "ERR_UNAUTHORIZED",
	5:  "ERR_FORBIDDEN",
	6:  "ERR_RATE_LIMIT",
	7:  "ERR_ILLEGAL_ACCESS",
	8:  "ERR_DATABASE",
	9:  "ERR_INVALID_ARGS",
	10: "ERR_INVALID_METHOD",
	11: "ERR_OPERATOR",
	12: "ERR_TASK_IN_PROGRESS",
	13: "ERR_INVALID_DEVICE",
}

var Error_Code_value = map[string]int32{
	"ERR_OK":                 0,
	"ERR_NOT_FOUND":          1,
	"ERR_METHOD_NOT_ALLOWED": 2,
	"ERR_INTERNAL_SERVER":    3,
	"ERR_UNAUTHORIZED":       4,
	"ERR_FORBIDDEN":          5,
	"ERR_RATE_LIMIT":         6,
	"ERR_ILLEGAL_ACCESS":     7,
	"ERR_DATABASE":           8,
	"ERR_INVALID_ARGS":       9,
	"ERR_INVALID_METHOD":     10,
	"ERR_OPERATOR":           11,
	"ERR_TASK_IN_PROGRESS":   12,
	"ERR_INVALID_DEVICE":     13,
}

func (x Error_Code) String() string {
	return proto.EnumName(Error_Code_name, int32(x))
}

func (Error_Code) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_13a9f1b5947140c8, []int{0, 0}
}

// 错误码
type Error struct {
	Code                 Error_Code `protobuf:"varint,1,opt,name=code,proto3,enum=iotpb.Error_Code" json:"code,omitempty"`
	Message              string     `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *Error) Reset()         { *m = Error{} }
func (m *Error) String() string { return proto.CompactTextString(m) }
func (*Error) ProtoMessage()    {}
func (*Error) Descriptor() ([]byte, []int) {
	return fileDescriptor_13a9f1b5947140c8, []int{0}
}

func (m *Error) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Error.Unmarshal(m, b)
}
func (m *Error) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Error.Marshal(b, m, deterministic)
}
func (m *Error) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Error.Merge(m, src)
}
func (m *Error) XXX_Size() int {
	return xxx_messageInfo_Error.Size(m)
}
func (m *Error) XXX_DiscardUnknown() {
	xxx_messageInfo_Error.DiscardUnknown(m)
}

var xxx_messageInfo_Error proto.InternalMessageInfo

func (m *Error) GetCode() Error_Code {
	if m != nil {
		return m.Code
	}
	return Error_ERR_OK
}

func (m *Error) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func init() {
	proto.RegisterEnum("iotpb.Operator", Operator_name, Operator_value)
	proto.RegisterEnum("iotpb.Status", Status_name, Status_value)
	proto.RegisterEnum("iotpb.ServiceStatus", ServiceStatus_name, ServiceStatus_value)
	proto.RegisterEnum("iotpb.SwitchStatus", SwitchStatus_name, SwitchStatus_value)
	proto.RegisterEnum("iotpb.SpeedLevel", SpeedLevel_name, SpeedLevel_value)
	proto.RegisterEnum("iotpb.Error_Code", Error_Code_name, Error_Code_value)
	proto.RegisterType((*Error)(nil), "iotpb.Error")
}

func init() { proto.RegisterFile("enum.proto", fileDescriptor_13a9f1b5947140c8) }

var fileDescriptor_13a9f1b5947140c8 = []byte{
	// 592 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x93, 0xdd, 0x4e, 0x1b, 0x3b,
	0x10, 0xc7, 0xd9, 0x25, 0x04, 0x18, 0x08, 0x0c, 0x86, 0xc3, 0x89, 0xce, 0x15, 0x42, 0x3a, 0x12,
	0xca, 0x45, 0x24, 0x38, 0x3a, 0xea, 0x6d, 0xbd, 0xeb, 0x01, 0x2c, 0x1c, 0x3b, 0xb2, 0x9d, 0x20,
	0x90, 0xaa, 0x15, 0x1f, 0xab, 0x16, 0xa9, 0x34, 0x51, 0x08, 0xf4, 0x45, 0xfa, 0x26, 0xbd, 0xee,
	0xab, 0x55, 0x95, 0xf7, 0x23, 0xd9, 0xf6, 0x2e, 0xf3, 0xfb, 0xcf, 0xfc, 0x3d, 0x1f, 0x59, 0x80,
	0xfc, 0xcb, 0xeb, 0x73, 0x7f, 0x3a, 0x9b, 0xcc, 0x27, 0x6c, 0xed, 0x69, 0x32, 0x9f, 0xde, 0x1f,
	0xff, 0x8c, 0x61, 0x8d, 0x66, 0xb3, 0xc9, 0x8c, 0xfd, 0x0b, 0xad, 0x87, 0xc9, 0x63, 0xde, 0x8d,
	0x8e, 0xa2, 0x93, 0x9d, 0xb3, 0xbd, 0x7e, 0xa1, 0xf7, 0x0b, 0xad, 0x9f, 0x4e, 0x1e, 0x73, 0x5b,
	0xc8, 0xac, 0x0b, 0xeb, 0xcf, 0xf9, 0xcb, 0xcb, 0xdd, 0xc7, 0xbc, 0x1b, 0x1f, 0x45, 0x27, 0x9b,
	0xb6, 0x0e, 0x8f, 0xbf, 0xc7, 0xd0, 0x0a, 0x89, 0x0c, 0xa0, 0x4d, 0xd6, 0x66, 0xe6, 0x0a, 0x57,
	0xd8, 0x1e, 0x74, 0xc2, 0x6f, 0x6d, 0x7c, 0x76, 0x6e, 0x46, 0x5a, 0x60, 0xc4, 0xfe, 0x81, 0xc3,
	0x80, 0x06, 0xe4, 0x2f, 0x8d, 0x28, 0x14, 0xae, 0x94, 0xb9, 0x26, 0x81, 0x31, 0xfb, 0x1b, 0xf6,
	0x83, 0x26, 0xb5, 0x27, 0xab, 0xb9, 0xca, 0x1c, 0xd9, 0x31, 0x59, 0x5c, 0x65, 0x07, 0x80, 0x41,
	0x18, 0x69, 0x3e, 0xf2, 0x97, 0xc6, 0xca, 0x5b, 0x12, 0xd8, 0xaa, 0xdd, 0xcf, 0x8d, 0x4d, 0xa4,
	0x10, 0xa4, 0x71, 0x8d, 0x31, 0xd8, 0x09, 0xc8, 0x72, 0x4f, 0x99, 0x92, 0x03, 0xe9, 0xb1, 0xcd,
	0x0e, 0x81, 0x15, 0xae, 0x4a, 0xd1, 0x05, 0x57, 0x19, 0x4f, 0x53, 0x72, 0x0e, 0xd7, 0x19, 0xc2,
	0x76, 0xe0, 0x82, 0x7b, 0x9e, 0x70, 0x47, 0xb8, 0x51, 0x3f, 0x23, 0xf5, 0x98, 0x2b, 0x29, 0x32,
	0x6e, 0x2f, 0x1c, 0x6e, 0x2e, 0xea, 0x2b, 0x5a, 0x76, 0x8e, 0x50, 0xd7, 0x9b, 0x21, 0x59, 0xee,
	0x8d, 0xc5, 0x2d, 0xd6, 0x85, 0x83, 0x40, 0x3c, 0x77, 0x57, 0x99, 0xd4, 0xd9, 0xd0, 0x9a, 0x0b,
	0x1b, 0xde, 0xda, 0xfe, 0xd3, 0x43, 0xd0, 0x58, 0xa6, 0x84, 0x9d, 0xde, 0x07, 0xd8, 0x30, 0xd3,
	0x7c, 0x76, 0x37, 0x9f, 0xcc, 0xc2, 0x38, 0xb5, 0x57, 0xa6, 0x8d, 0x26, 0x5c, 0x61, 0xfb, 0xb0,
	0xbb, 0x40, 0x23, 0x2d, 0x53, 0x33, 0xc0, 0xe8, 0x37, 0x38, 0x30, 0x89, 0x54, 0x84, 0x71, 0x68,
	0x7d, 0x01, 0x3d, 0x29, 0x0a, 0xa9, 0xab, 0xbd, 0x6f, 0x11, 0xb4, 0xdd, 0xfc, 0x6e, 0xfe, 0xfa,
	0xc2, 0x76, 0x61, 0xcb, 0x79, 0xee, 0x47, 0xae, 0xf6, 0x5e, 0x02, 0x4f, 0xce, 0x63, 0x14, 0x2c,
	0x2a, 0x20, 0xf5, 0x98, 0xb4, 0x37, 0xf6, 0x06, 0xe3, 0x30, 0x65, 0x45, 0x2d, 0x71, 0x71, 0x53,
	0x1e, 0xa3, 0x22, 0x3c, 0xf5, 0x72, 0xcc, 0x7d, 0x71, 0x8c, 0x43, 0x60, 0x15, 0x15, 0xb4, 0xe4,
	0xc5, 0x45, 0x16, 0xf5, 0x5e, 0x5a, 0x12, 0xd8, 0xee, 0x29, 0xe8, 0xb8, 0x7c, 0xf6, 0xf6, 0xf4,
	0x90, 0x57, 0xcd, 0x85, 0x47, 0xc8, 0x86, 0x9d, 0x34, 0x26, 0xaf, 0x09, 0x69, 0x9e, 0x28, 0x12,
	0x55, 0x87, 0x15, 0x14, 0xd2, 0x95, 0x34, 0xee, 0xfd, 0x88, 0x60, 0xdb, 0x7d, 0x7d, 0x9a, 0x3f,
	0x7c, 0x6a, 0x8c, 0x7a, 0x2d, 0x7d, 0x7a, 0x59, 0x9b, 0x85, 0xba, 0x12, 0xf0, 0x31, 0x97, 0x2a,
	0x14, 0x96, 0x7b, 0xac, 0x69, 0xd5, 0x2f, 0xc6, 0xec, 0x2f, 0xd8, 0xab, 0xe0, 0x72, 0x0c, 0x6c,
	0x35, 0x1c, 0xcc, 0x90, 0x74, 0xf1, 0x9f, 0xc1, 0x8d, 0x46, 0x72, 0xaa, 0x8c, 0xa3, 0x12, 0x63,
	0x03, 0x17, 0xc9, 0x63, 0x13, 0x6e, 0x7d, 0x54, 0x6c, 0xa8, 0x99, 0x5d, 0xf2, 0xf7, 0xbd, 0x5b,
	0x00, 0x37, 0xcd, 0xf3, 0x47, 0x95, 0xbf, 0xe5, 0x9f, 0xd9, 0x0e, 0x80, 0x1b, 0x12, 0x89, 0xe6,
	0x99, 0x8a, 0xf8, 0xff, 0xd3, 0xb3, 0xab, 0x04, 0x23, 0xd6, 0x81, 0xcd, 0x12, 0x9c, 0x0e, 0x12,
	0x8c, 0x97, 0xe1, 0xd9, 0x20, 0xc1, 0xd5, 0x65, 0xf8, 0x6e, 0x90, 0x60, 0xeb, 0xbe, 0x5d, 0x7c,
	0xee, 0xff, 0xfd, 0x0a, 0x00, 0x00, 0xff, 0xff, 0x74, 0xdd, 0xca, 0x43, 0xfc, 0x03, 0x00, 0x00,
}