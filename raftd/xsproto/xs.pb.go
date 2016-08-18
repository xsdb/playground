// Code generated by protoc-gen-go.
// source: xs.proto
// DO NOT EDIT!

/*
Package xsproto is a generated protocol buffer package.

It is generated from these files:
	xs.proto

It has these top-level messages:
	Request
	Response
*/
package xsproto

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Request struct {
	Key string `protobuf:"bytes,1,opt,name=key" json:"key,omitempty"`
	Val string `protobuf:"bytes,2,opt,name=val" json:"val,omitempty"`
}

func (m *Request) Reset()                    { *m = Request{} }
func (m *Request) String() string            { return proto.CompactTextString(m) }
func (*Request) ProtoMessage()               {}
func (*Request) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type Response struct {
	ReturnCode int32  `protobuf:"varint,1,opt,name=return_code,json=returnCode" json:"return_code,omitempty"`
	Value      string `protobuf:"bytes,2,opt,name=value" json:"value,omitempty"`
}

func (m *Response) Reset()                    { *m = Response{} }
func (m *Response) String() string            { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()               {}
func (*Response) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func init() {
	proto.RegisterType((*Request)(nil), "xsproto.Request")
	proto.RegisterType((*Response)(nil), "xsproto.Response")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion3

// Client API for XsProto service

type XsProtoClient interface {
	Get(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error)
	Put(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error)
}

type xsProtoClient struct {
	cc *grpc.ClientConn
}

func NewXsProtoClient(cc *grpc.ClientConn) XsProtoClient {
	return &xsProtoClient{cc}
}

func (c *xsProtoClient) Get(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := grpc.Invoke(ctx, "/xsproto.XsProto/Get", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *xsProtoClient) Put(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := grpc.Invoke(ctx, "/xsproto.XsProto/Put", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for XsProto service

type XsProtoServer interface {
	Get(context.Context, *Request) (*Response, error)
	Put(context.Context, *Request) (*Response, error)
}

func RegisterXsProtoServer(s *grpc.Server, srv XsProtoServer) {
	s.RegisterService(&_XsProto_serviceDesc, srv)
}

func _XsProto_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(XsProtoServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/xsproto.XsProto/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(XsProtoServer).Get(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _XsProto_Put_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(XsProtoServer).Put(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/xsproto.XsProto/Put",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(XsProtoServer).Put(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

var _XsProto_serviceDesc = grpc.ServiceDesc{
	ServiceName: "xsproto.XsProto",
	HandlerType: (*XsProtoServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Get",
			Handler:    _XsProto_Get_Handler,
		},
		{
			MethodName: "Put",
			Handler:    _XsProto_Put_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: fileDescriptor0,
}

func init() { proto.RegisterFile("xs.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 197 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0xa8, 0x28, 0xd6, 0x2b,
	0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0xaf, 0x28, 0x06, 0x33, 0x94, 0x74, 0xb9, 0xd8, 0x83, 0x52,
	0x0b, 0x4b, 0x53, 0x8b, 0x4b, 0x84, 0x04, 0xb8, 0x98, 0xb3, 0x53, 0x2b, 0x25, 0x18, 0x15, 0x18,
	0x35, 0x38, 0x83, 0x40, 0x4c, 0x90, 0x48, 0x59, 0x62, 0x8e, 0x04, 0x13, 0x44, 0x04, 0xc8, 0x54,
	0x72, 0xe4, 0xe2, 0x08, 0x4a, 0x2d, 0x2e, 0xc8, 0xcf, 0x2b, 0x4e, 0x15, 0x92, 0xe7, 0xe2, 0x2e,
	0x4a, 0x2d, 0x29, 0x2d, 0xca, 0x8b, 0x4f, 0xce, 0x4f, 0x49, 0x05, 0xeb, 0x63, 0x0d, 0xe2, 0x82,
	0x08, 0x39, 0x03, 0x45, 0x84, 0x44, 0xb8, 0x58, 0x81, 0x7a, 0x4a, 0x53, 0xa1, 0x06, 0x40, 0x38,
	0x46, 0xa9, 0x5c, 0xec, 0x11, 0xc5, 0x01, 0x60, 0x57, 0xe8, 0x70, 0x31, 0xbb, 0xa7, 0x02, 0x2d,
	0xd6, 0x83, 0xba, 0x46, 0x0f, 0xea, 0x14, 0x29, 0x41, 0x24, 0x11, 0x88, 0x6d, 0x4a, 0x0c, 0x20,
	0xd5, 0x01, 0xa5, 0xc4, 0xaa, 0x76, 0xd2, 0xe2, 0x92, 0xc8, 0xcc, 0xd7, 0x4b, 0x2f, 0x2a, 0x48,
	0xd6, 0x4b, 0xad, 0x48, 0xcc, 0x2d, 0xc8, 0x49, 0x2d, 0x86, 0x29, 0x73, 0xe2, 0x81, 0x3a, 0x00,
	0x42, 0x30, 0x26, 0xb1, 0x81, 0x85, 0x8d, 0x01, 0x01, 0x00, 0x00, 0xff, 0xff, 0x60, 0x2d, 0xe4,
	0x4c, 0x20, 0x01, 0x00, 0x00,
}
