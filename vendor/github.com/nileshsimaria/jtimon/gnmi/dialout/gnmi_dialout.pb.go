// Code generated by protoc-gen-go. DO NOT EDIT.
// source: gnmi_dialout.proto

/*
Package gnmi_dialout is a generated protocol buffer package.

It is generated from these files:
	gnmi_dialout.proto

It has these top-level messages:
*/
package gnmi_dialout

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import gnmi "github.com/nileshsimaria/jtimon/gnmi/gnmi"

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

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Subscriber service

type SubscriberClient interface {
	DialOutSubscriber(ctx context.Context, opts ...grpc.CallOption) (Subscriber_DialOutSubscriberClient, error)
}

type subscriberClient struct {
	cc *grpc.ClientConn
}

func NewSubscriberClient(cc *grpc.ClientConn) SubscriberClient {
	return &subscriberClient{cc}
}

func (c *subscriberClient) DialOutSubscriber(ctx context.Context, opts ...grpc.CallOption) (Subscriber_DialOutSubscriberClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_Subscriber_serviceDesc.Streams[0], c.cc, "/Subscriber/DialOutSubscriber", opts...)
	if err != nil {
		return nil, err
	}
	x := &subscriberDialOutSubscriberClient{stream}
	return x, nil
}

type Subscriber_DialOutSubscriberClient interface {
	Send(*gnmi.SubscribeResponse) error
	Recv() (*gnmi.SubscribeRequest, error)
	grpc.ClientStream
}

type subscriberDialOutSubscriberClient struct {
	grpc.ClientStream
}

func (x *subscriberDialOutSubscriberClient) Send(m *gnmi.SubscribeResponse) error {
	return x.ClientStream.SendMsg(m)
}

func (x *subscriberDialOutSubscriberClient) Recv() (*gnmi.SubscribeRequest, error) {
	m := new(gnmi.SubscribeRequest)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for Subscriber service

type SubscriberServer interface {
	DialOutSubscriber(Subscriber_DialOutSubscriberServer) error
}

func RegisterSubscriberServer(s *grpc.Server, srv SubscriberServer) {
	s.RegisterService(&_Subscriber_serviceDesc, srv)
}

func _Subscriber_DialOutSubscriber_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(SubscriberServer).DialOutSubscriber(&subscriberDialOutSubscriberServer{stream})
}

type Subscriber_DialOutSubscriberServer interface {
	Send(*gnmi.SubscribeRequest) error
	Recv() (*gnmi.SubscribeResponse, error)
	grpc.ServerStream
}

type subscriberDialOutSubscriberServer struct {
	grpc.ServerStream
}

func (x *subscriberDialOutSubscriberServer) Send(m *gnmi.SubscribeRequest) error {
	return x.ServerStream.SendMsg(m)
}

func (x *subscriberDialOutSubscriberServer) Recv() (*gnmi.SubscribeResponse, error) {
	m := new(gnmi.SubscribeResponse)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _Subscriber_serviceDesc = grpc.ServiceDesc{
	ServiceName: "Subscriber",
	HandlerType: (*SubscriberServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "DialOutSubscriber",
			Handler:       _Subscriber_DialOutSubscriber_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "gnmi_dialout.proto",
}

func init() { proto.RegisterFile("gnmi_dialout.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 143 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x4a, 0xcf, 0xcb, 0xcd,
	0x8c, 0x4f, 0xc9, 0x4c, 0xcc, 0xc9, 0x2f, 0x2d, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x97, 0xe2,
	0x02, 0x89, 0x41, 0xd8, 0x46, 0x61, 0x5c, 0x5c, 0xc1, 0xa5, 0x49, 0xc5, 0xc9, 0x45, 0x99, 0x49,
	0xa9, 0x45, 0x42, 0x1e, 0x5c, 0x82, 0x2e, 0x99, 0x89, 0x39, 0xfe, 0xa5, 0x25, 0x48, 0x82, 0xe2,
	0x7a, 0x60, 0xf5, 0x70, 0x91, 0xa0, 0xd4, 0xe2, 0x82, 0xfc, 0xbc, 0xe2, 0x54, 0x29, 0x31, 0x0c,
	0x89, 0xc2, 0xd2, 0xd4, 0xe2, 0x12, 0x0d, 0x46, 0x03, 0x46, 0x27, 0xfd, 0x28, 0xdd, 0xf4, 0xcc,
	0x92, 0x8c, 0xd2, 0x24, 0xbd, 0xe4, 0xfc, 0x5c, 0xfd, 0xfc, 0x82, 0xd4, 0xbc, 0xe4, 0xfc, 0xbc,
	0xb4, 0xcc, 0x74, 0x7d, 0x90, 0x16, 0x7d, 0xb0, 0xdd, 0xfa, 0xc8, 0x4e, 0x4b, 0x62, 0x03, 0x8b,
	0x19, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0x5a, 0xa4, 0x0d, 0x6b, 0xb1, 0x00, 0x00, 0x00,
}
