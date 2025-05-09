// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v6.30.0
// source: scanner.proto

package __

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	ScannerService_ExecuteScan_FullMethodName   = "/byakugan.scanner.ScannerService/ExecuteScan"
	ScannerService_StreamResults_FullMethodName = "/byakugan.scanner.ScannerService/StreamResults"
	ScannerService_GetTaskStatus_FullMethodName = "/byakugan.scanner.ScannerService/GetTaskStatus"
	ScannerService_GetScanStatus_FullMethodName = "/byakugan.scanner.ScannerService/GetScanStatus"
	ScannerService_Heartbeat_FullMethodName     = "/byakugan.scanner.ScannerService/Heartbeat"
)

// ScannerServiceClient is the client API for ScannerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// Scanner service definition
type ScannerServiceClient interface {
	ExecuteScan(ctx context.Context, in *ScanTask, opts ...grpc.CallOption) (*ScanResult, error)
	StreamResults(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[ScanResult, StreamAck], error)
	GetTaskStatus(ctx context.Context, in *TaskStatusRequest, opts ...grpc.CallOption) (*TaskStatus, error)
	GetScanStatus(ctx context.Context, in *StatusRequest, opts ...grpc.CallOption) (*ScanStatus, error)
	Heartbeat(ctx context.Context, in *HeartbeatRequest, opts ...grpc.CallOption) (*HeartbeatResponse, error)
}

type scannerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewScannerServiceClient(cc grpc.ClientConnInterface) ScannerServiceClient {
	return &scannerServiceClient{cc}
}

func (c *scannerServiceClient) ExecuteScan(ctx context.Context, in *ScanTask, opts ...grpc.CallOption) (*ScanResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ScanResult)
	err := c.cc.Invoke(ctx, ScannerService_ExecuteScan_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *scannerServiceClient) StreamResults(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[ScanResult, StreamAck], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &ScannerService_ServiceDesc.Streams[0], ScannerService_StreamResults_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[ScanResult, StreamAck]{ClientStream: stream}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type ScannerService_StreamResultsClient = grpc.BidiStreamingClient[ScanResult, StreamAck]

func (c *scannerServiceClient) GetTaskStatus(ctx context.Context, in *TaskStatusRequest, opts ...grpc.CallOption) (*TaskStatus, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(TaskStatus)
	err := c.cc.Invoke(ctx, ScannerService_GetTaskStatus_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *scannerServiceClient) GetScanStatus(ctx context.Context, in *StatusRequest, opts ...grpc.CallOption) (*ScanStatus, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ScanStatus)
	err := c.cc.Invoke(ctx, ScannerService_GetScanStatus_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *scannerServiceClient) Heartbeat(ctx context.Context, in *HeartbeatRequest, opts ...grpc.CallOption) (*HeartbeatResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(HeartbeatResponse)
	err := c.cc.Invoke(ctx, ScannerService_Heartbeat_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ScannerServiceServer is the server API for ScannerService service.
// All implementations must embed UnimplementedScannerServiceServer
// for forward compatibility.
//
// Scanner service definition
type ScannerServiceServer interface {
	ExecuteScan(context.Context, *ScanTask) (*ScanResult, error)
	StreamResults(grpc.BidiStreamingServer[ScanResult, StreamAck]) error
	GetTaskStatus(context.Context, *TaskStatusRequest) (*TaskStatus, error)
	GetScanStatus(context.Context, *StatusRequest) (*ScanStatus, error)
	Heartbeat(context.Context, *HeartbeatRequest) (*HeartbeatResponse, error)
	mustEmbedUnimplementedScannerServiceServer()
}

// UnimplementedScannerServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedScannerServiceServer struct{}

func (UnimplementedScannerServiceServer) ExecuteScan(context.Context, *ScanTask) (*ScanResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExecuteScan not implemented")
}
func (UnimplementedScannerServiceServer) StreamResults(grpc.BidiStreamingServer[ScanResult, StreamAck]) error {
	return status.Errorf(codes.Unimplemented, "method StreamResults not implemented")
}
func (UnimplementedScannerServiceServer) GetTaskStatus(context.Context, *TaskStatusRequest) (*TaskStatus, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTaskStatus not implemented")
}
func (UnimplementedScannerServiceServer) GetScanStatus(context.Context, *StatusRequest) (*ScanStatus, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetScanStatus not implemented")
}
func (UnimplementedScannerServiceServer) Heartbeat(context.Context, *HeartbeatRequest) (*HeartbeatResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Heartbeat not implemented")
}
func (UnimplementedScannerServiceServer) mustEmbedUnimplementedScannerServiceServer() {}
func (UnimplementedScannerServiceServer) testEmbeddedByValue()                        {}

// UnsafeScannerServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ScannerServiceServer will
// result in compilation errors.
type UnsafeScannerServiceServer interface {
	mustEmbedUnimplementedScannerServiceServer()
}

func RegisterScannerServiceServer(s grpc.ServiceRegistrar, srv ScannerServiceServer) {
	// If the following call pancis, it indicates UnimplementedScannerServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&ScannerService_ServiceDesc, srv)
}

func _ScannerService_ExecuteScan_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ScanTask)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ScannerServiceServer).ExecuteScan(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ScannerService_ExecuteScan_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ScannerServiceServer).ExecuteScan(ctx, req.(*ScanTask))
	}
	return interceptor(ctx, in, info, handler)
}

func _ScannerService_StreamResults_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ScannerServiceServer).StreamResults(&grpc.GenericServerStream[ScanResult, StreamAck]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type ScannerService_StreamResultsServer = grpc.BidiStreamingServer[ScanResult, StreamAck]

func _ScannerService_GetTaskStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TaskStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ScannerServiceServer).GetTaskStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ScannerService_GetTaskStatus_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ScannerServiceServer).GetTaskStatus(ctx, req.(*TaskStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ScannerService_GetScanStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ScannerServiceServer).GetScanStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ScannerService_GetScanStatus_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ScannerServiceServer).GetScanStatus(ctx, req.(*StatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ScannerService_Heartbeat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HeartbeatRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ScannerServiceServer).Heartbeat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ScannerService_Heartbeat_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ScannerServiceServer).Heartbeat(ctx, req.(*HeartbeatRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ScannerService_ServiceDesc is the grpc.ServiceDesc for ScannerService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ScannerService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "byakugan.scanner.ScannerService",
	HandlerType: (*ScannerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ExecuteScan",
			Handler:    _ScannerService_ExecuteScan_Handler,
		},
		{
			MethodName: "GetTaskStatus",
			Handler:    _ScannerService_GetTaskStatus_Handler,
		},
		{
			MethodName: "GetScanStatus",
			Handler:    _ScannerService_GetScanStatus_Handler,
		},
		{
			MethodName: "Heartbeat",
			Handler:    _ScannerService_Heartbeat_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamResults",
			Handler:       _ScannerService_StreamResults_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "scanner.proto",
}
