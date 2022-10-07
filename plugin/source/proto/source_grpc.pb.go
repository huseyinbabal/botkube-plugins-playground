// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.6
// source: plugin/source/proto/source.proto

package source

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// SourceClient is the client API for Source service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SourceClient interface {
	Consume(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (Source_ConsumeClient, error)
}

type sourceClient struct {
	cc grpc.ClientConnInterface
}

func NewSourceClient(cc grpc.ClientConnInterface) SourceClient {
	return &sourceClient{cc}
}

func (c *sourceClient) Consume(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (Source_ConsumeClient, error) {
	stream, err := c.cc.NewStream(ctx, &Source_ServiceDesc.Streams[0], "/source.Source/Consume", opts...)
	if err != nil {
		return nil, err
	}
	x := &sourceConsumeClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Source_ConsumeClient interface {
	Recv() (*ConsumeResponse, error)
	grpc.ClientStream
}

type sourceConsumeClient struct {
	grpc.ClientStream
}

func (x *sourceConsumeClient) Recv() (*ConsumeResponse, error) {
	m := new(ConsumeResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// SourceServer is the server API for Source service.
// All implementations must embed UnimplementedSourceServer
// for forward compatibility
type SourceServer interface {
	Consume(*emptypb.Empty, Source_ConsumeServer) error
	mustEmbedUnimplementedSourceServer()
}

// UnimplementedSourceServer must be embedded to have forward compatible implementations.
type UnimplementedSourceServer struct {
}

func (UnimplementedSourceServer) Consume(*emptypb.Empty, Source_ConsumeServer) error {
	return status.Errorf(codes.Unimplemented, "method Consume not implemented")
}
func (UnimplementedSourceServer) mustEmbedUnimplementedSourceServer() {}

// UnsafeSourceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SourceServer will
// result in compilation errors.
type UnsafeSourceServer interface {
	mustEmbedUnimplementedSourceServer()
}

func RegisterSourceServer(s grpc.ServiceRegistrar, srv SourceServer) {
	s.RegisterService(&Source_ServiceDesc, srv)
}

func _Source_Consume_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(emptypb.Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(SourceServer).Consume(m, &sourceConsumeServer{stream})
}

type Source_ConsumeServer interface {
	Send(*ConsumeResponse) error
	grpc.ServerStream
}

type sourceConsumeServer struct {
	grpc.ServerStream
}

func (x *sourceConsumeServer) Send(m *ConsumeResponse) error {
	return x.ServerStream.SendMsg(m)
}

// Source_ServiceDesc is the grpc.ServiceDesc for Source service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Source_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "source.Source",
	HandlerType: (*SourceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Consume",
			Handler:       _Source_Consume_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "plugin/source/proto/source.proto",
}
