// Code generated by protoc-gen-go. DO NOT EDIT.
// source: api/items.proto

package items

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type ItemList struct {
	Items                []int64  `protobuf:"varint,1,rep,packed,name=items,proto3" json:"items,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ItemList) Reset()         { *m = ItemList{} }
func (m *ItemList) String() string { return proto.CompactTextString(m) }
func (*ItemList) ProtoMessage()    {}
func (*ItemList) Descriptor() ([]byte, []int) {
	return fileDescriptor_78f55b0301ff2c5c, []int{0}
}

func (m *ItemList) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ItemList.Unmarshal(m, b)
}
func (m *ItemList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ItemList.Marshal(b, m, deterministic)
}
func (m *ItemList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ItemList.Merge(m, src)
}
func (m *ItemList) XXX_Size() int {
	return xxx_messageInfo_ItemList.Size(m)
}
func (m *ItemList) XXX_DiscardUnknown() {
	xxx_messageInfo_ItemList.DiscardUnknown(m)
}

var xxx_messageInfo_ItemList proto.InternalMessageInfo

func (m *ItemList) GetItems() []int64 {
	if m != nil {
		return m.Items
	}
	return nil
}

type Empty struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Empty) Reset()         { *m = Empty{} }
func (m *Empty) String() string { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()    {}
func (*Empty) Descriptor() ([]byte, []int) {
	return fileDescriptor_78f55b0301ff2c5c, []int{1}
}

func (m *Empty) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Empty.Unmarshal(m, b)
}
func (m *Empty) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Empty.Marshal(b, m, deterministic)
}
func (m *Empty) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Empty.Merge(m, src)
}
func (m *Empty) XXX_Size() int {
	return xxx_messageInfo_Empty.Size(m)
}
func (m *Empty) XXX_DiscardUnknown() {
	xxx_messageInfo_Empty.DiscardUnknown(m)
}

var xxx_messageInfo_Empty proto.InternalMessageInfo

type Text struct {
	Text                 string   `protobuf:"bytes,1,opt,name=text,proto3" json:"text,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Text) Reset()         { *m = Text{} }
func (m *Text) String() string { return proto.CompactTextString(m) }
func (*Text) ProtoMessage()    {}
func (*Text) Descriptor() ([]byte, []int) {
	return fileDescriptor_78f55b0301ff2c5c, []int{2}
}

func (m *Text) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Text.Unmarshal(m, b)
}
func (m *Text) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Text.Marshal(b, m, deterministic)
}
func (m *Text) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Text.Merge(m, src)
}
func (m *Text) XXX_Size() int {
	return xxx_messageInfo_Text.Size(m)
}
func (m *Text) XXX_DiscardUnknown() {
	xxx_messageInfo_Text.DiscardUnknown(m)
}

var xxx_messageInfo_Text proto.InternalMessageInfo

func (m *Text) GetText() string {
	if m != nil {
		return m.Text
	}
	return ""
}

func init() {
	proto.RegisterType((*ItemList)(nil), "ItemList")
	proto.RegisterType((*Empty)(nil), "Empty")
	proto.RegisterType((*Text)(nil), "Text")
}

func init() {
	proto.RegisterFile("api/items.proto", fileDescriptor_78f55b0301ff2c5c)
}

var fileDescriptor_78f55b0301ff2c5c = []byte{
	// 187 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4f, 0x2c, 0xc8, 0xd4,
	0xcf, 0x2c, 0x49, 0xcd, 0x2d, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x57, 0x52, 0xe0, 0xe2, 0xf0,
	0x2c, 0x49, 0xcd, 0xf5, 0xc9, 0x2c, 0x2e, 0x11, 0x12, 0xe1, 0x62, 0x05, 0x4b, 0x49, 0x30, 0x2a,
	0x30, 0x6b, 0x30, 0x07, 0x41, 0x38, 0x4a, 0xec, 0x5c, 0xac, 0xae, 0xb9, 0x05, 0x25, 0x95, 0x4a,
	0x52, 0x5c, 0x2c, 0x21, 0xa9, 0x15, 0x25, 0x42, 0x42, 0x5c, 0x2c, 0x25, 0xa9, 0x15, 0x25, 0x12,
	0x8c, 0x0a, 0x8c, 0x1a, 0x9c, 0x41, 0x60, 0xb6, 0x51, 0x3b, 0x23, 0x97, 0x40, 0x78, 0x62, 0x51,
	0x6a, 0x46, 0x7e, 0x69, 0x71, 0x6a, 0x70, 0x6a, 0x51, 0x59, 0x66, 0x72, 0xaa, 0x90, 0x2c, 0x17,
	0x47, 0x40, 0x69, 0x09, 0xc8, 0xf8, 0x62, 0x21, 0x4e, 0x3d, 0x98, 0x35, 0x52, 0x6c, 0x7a, 0x60,
	0xf3, 0x84, 0xe4, 0xb8, 0x38, 0x43, 0x12, 0xb3, 0x53, 0x71, 0xca, 0xcb, 0x72, 0x71, 0xb8, 0xa7,
	0x42, 0xb5, 0x43, 0xc5, 0xa4, 0x10, 0xca, 0x84, 0x44, 0xb9, 0x58, 0x3d, 0x52, 0x73, 0x72, 0xf2,
	0x85, 0x58, 0xf5, 0x40, 0xce, 0x92, 0x82, 0x50, 0x49, 0x6c, 0x60, 0x7f, 0x19, 0x03, 0x02, 0x00,
	0x00, 0xff, 0xff, 0xaf, 0xb2, 0xa4, 0x80, 0xea, 0x00, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// WarehouseServiceClient is the client API for WarehouseService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type WarehouseServiceClient interface {
	PutItems(ctx context.Context, in *ItemList, opts ...grpc.CallOption) (*Empty, error)
	TakeItems(ctx context.Context, in *ItemList, opts ...grpc.CallOption) (*Empty, error)
	GetItems(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*ItemList, error)
	Hello(ctx context.Context, in *Text, opts ...grpc.CallOption) (*Text, error)
}

type warehouseServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewWarehouseServiceClient(cc grpc.ClientConnInterface) WarehouseServiceClient {
	return &warehouseServiceClient{cc}
}

func (c *warehouseServiceClient) PutItems(ctx context.Context, in *ItemList, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/WarehouseService/PutItems", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *warehouseServiceClient) TakeItems(ctx context.Context, in *ItemList, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/WarehouseService/TakeItems", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *warehouseServiceClient) GetItems(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*ItemList, error) {
	out := new(ItemList)
	err := c.cc.Invoke(ctx, "/WarehouseService/GetItems", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *warehouseServiceClient) Hello(ctx context.Context, in *Text, opts ...grpc.CallOption) (*Text, error) {
	out := new(Text)
	err := c.cc.Invoke(ctx, "/WarehouseService/Hello", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// WarehouseServiceServer is the server API for WarehouseService service.
type WarehouseServiceServer interface {
	PutItems(context.Context, *ItemList) (*Empty, error)
	TakeItems(context.Context, *ItemList) (*Empty, error)
	GetItems(context.Context, *Empty) (*ItemList, error)
	Hello(context.Context, *Text) (*Text, error)
}

// UnimplementedWarehouseServiceServer can be embedded to have forward compatible implementations.
type UnimplementedWarehouseServiceServer struct {
}

func (*UnimplementedWarehouseServiceServer) PutItems(ctx context.Context, req *ItemList) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PutItems not implemented")
}
func (*UnimplementedWarehouseServiceServer) TakeItems(ctx context.Context, req *ItemList) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TakeItems not implemented")
}
func (*UnimplementedWarehouseServiceServer) GetItems(ctx context.Context, req *Empty) (*ItemList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetItems not implemented")
}
func (*UnimplementedWarehouseServiceServer) Hello(ctx context.Context, req *Text) (*Text, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Hello not implemented")
}

func RegisterWarehouseServiceServer(s *grpc.Server, srv WarehouseServiceServer) {
	s.RegisterService(&_WarehouseService_serviceDesc, srv)
}

func _WarehouseService_PutItems_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ItemList)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WarehouseServiceServer).PutItems(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/WarehouseService/PutItems",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WarehouseServiceServer).PutItems(ctx, req.(*ItemList))
	}
	return interceptor(ctx, in, info, handler)
}

func _WarehouseService_TakeItems_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ItemList)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WarehouseServiceServer).TakeItems(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/WarehouseService/TakeItems",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WarehouseServiceServer).TakeItems(ctx, req.(*ItemList))
	}
	return interceptor(ctx, in, info, handler)
}

func _WarehouseService_GetItems_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WarehouseServiceServer).GetItems(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/WarehouseService/GetItems",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WarehouseServiceServer).GetItems(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _WarehouseService_Hello_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Text)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WarehouseServiceServer).Hello(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/WarehouseService/Hello",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WarehouseServiceServer).Hello(ctx, req.(*Text))
	}
	return interceptor(ctx, in, info, handler)
}

var _WarehouseService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "WarehouseService",
	HandlerType: (*WarehouseServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PutItems",
			Handler:    _WarehouseService_PutItems_Handler,
		},
		{
			MethodName: "TakeItems",
			Handler:    _WarehouseService_TakeItems_Handler,
		},
		{
			MethodName: "GetItems",
			Handler:    _WarehouseService_GetItems_Handler,
		},
		{
			MethodName: "Hello",
			Handler:    _WarehouseService_Hello_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/items.proto",
}
