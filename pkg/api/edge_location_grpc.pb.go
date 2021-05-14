// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package api

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// EdgeLocationsClient is the client API for EdgeLocations service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EdgeLocationsClient interface {
	Register(ctx context.Context, in *EdgeLocation, opts ...grpc.CallOption) (*emptypb.Empty, error)
	Get(ctx context.Context, in *EdgeLocation, opts ...grpc.CallOption) (*EdgeLocation, error)
	List(ctx context.Context, in *ListEdgeLocationParams, opts ...grpc.CallOption) (EdgeLocations_ListClient, error)
}

type edgeLocationsClient struct {
	cc grpc.ClientConnInterface
}

func NewEdgeLocationsClient(cc grpc.ClientConnInterface) EdgeLocationsClient {
	return &edgeLocationsClient{cc}
}

func (c *edgeLocationsClient) Register(ctx context.Context, in *EdgeLocation, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/edge_location.EdgeLocations/Register", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *edgeLocationsClient) Get(ctx context.Context, in *EdgeLocation, opts ...grpc.CallOption) (*EdgeLocation, error) {
	out := new(EdgeLocation)
	err := c.cc.Invoke(ctx, "/edge_location.EdgeLocations/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *edgeLocationsClient) List(ctx context.Context, in *ListEdgeLocationParams, opts ...grpc.CallOption) (EdgeLocations_ListClient, error) {
	stream, err := c.cc.NewStream(ctx, &_EdgeLocations_serviceDesc.Streams[0], "/edge_location.EdgeLocations/List", opts...)
	if err != nil {
		return nil, err
	}
	x := &edgeLocationsListClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type EdgeLocations_ListClient interface {
	Recv() (*EdgeLocation, error)
	grpc.ClientStream
}

type edgeLocationsListClient struct {
	grpc.ClientStream
}

func (x *edgeLocationsListClient) Recv() (*EdgeLocation, error) {
	m := new(EdgeLocation)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// EdgeLocationsServer is the server API for EdgeLocations service.
// All implementations must embed UnimplementedEdgeLocationsServer
// for forward compatibility
type EdgeLocationsServer interface {
	Register(context.Context, *EdgeLocation) (*emptypb.Empty, error)
	Get(context.Context, *EdgeLocation) (*EdgeLocation, error)
	List(*ListEdgeLocationParams, EdgeLocations_ListServer) error
	mustEmbedUnimplementedEdgeLocationsServer()
}

// UnimplementedEdgeLocationsServer must be embedded to have forward compatible implementations.
type UnimplementedEdgeLocationsServer struct {
}

func (UnimplementedEdgeLocationsServer) Register(context.Context, *EdgeLocation) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (UnimplementedEdgeLocationsServer) Get(context.Context, *EdgeLocation) (*EdgeLocation, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedEdgeLocationsServer) List(*ListEdgeLocationParams, EdgeLocations_ListServer) error {
	return status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (UnimplementedEdgeLocationsServer) mustEmbedUnimplementedEdgeLocationsServer() {}

// UnsafeEdgeLocationsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EdgeLocationsServer will
// result in compilation errors.
type UnsafeEdgeLocationsServer interface {
	mustEmbedUnimplementedEdgeLocationsServer()
}

func RegisterEdgeLocationsServer(s grpc.ServiceRegistrar, srv EdgeLocationsServer) {
	s.RegisterService(&_EdgeLocations_serviceDesc, srv)
}

func _EdgeLocations_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EdgeLocation)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EdgeLocationsServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/edge_location.EdgeLocations/Register",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EdgeLocationsServer).Register(ctx, req.(*EdgeLocation))
	}
	return interceptor(ctx, in, info, handler)
}

func _EdgeLocations_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EdgeLocation)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EdgeLocationsServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/edge_location.EdgeLocations/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EdgeLocationsServer).Get(ctx, req.(*EdgeLocation))
	}
	return interceptor(ctx, in, info, handler)
}

func _EdgeLocations_List_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ListEdgeLocationParams)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(EdgeLocationsServer).List(m, &edgeLocationsListServer{stream})
}

type EdgeLocations_ListServer interface {
	Send(*EdgeLocation) error
	grpc.ServerStream
}

type edgeLocationsListServer struct {
	grpc.ServerStream
}

func (x *edgeLocationsListServer) Send(m *EdgeLocation) error {
	return x.ServerStream.SendMsg(m)
}

var _EdgeLocations_serviceDesc = grpc.ServiceDesc{
	ServiceName: "edge_location.EdgeLocations",
	HandlerType: (*EdgeLocationsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Register",
			Handler:    _EdgeLocations_Register_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _EdgeLocations_Get_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "List",
			Handler:       _EdgeLocations_List_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "api/edge_location.proto",
}