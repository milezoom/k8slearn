// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.1
// source: contract/printsvc.proto

package contract

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	PrintService_PrintGreeting_FullMethodName = "/printsvc.PrintService/PrintGreeting"
)

// PrintServiceClient is the client API for PrintService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PrintServiceClient interface {
	PrintGreeting(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*PrintGreetingResponse, error)
}

type printServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPrintServiceClient(cc grpc.ClientConnInterface) PrintServiceClient {
	return &printServiceClient{cc}
}

func (c *printServiceClient) PrintGreeting(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*PrintGreetingResponse, error) {
	out := new(PrintGreetingResponse)
	err := c.cc.Invoke(ctx, PrintService_PrintGreeting_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PrintServiceServer is the server API for PrintService service.
// All implementations must embed UnimplementedPrintServiceServer
// for forward compatibility
type PrintServiceServer interface {
	PrintGreeting(context.Context, *EmptyRequest) (*PrintGreetingResponse, error)
	mustEmbedUnimplementedPrintServiceServer()
}

// UnimplementedPrintServiceServer must be embedded to have forward compatible implementations.
type UnimplementedPrintServiceServer struct {
}

func (UnimplementedPrintServiceServer) PrintGreeting(context.Context, *EmptyRequest) (*PrintGreetingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PrintGreeting not implemented")
}
func (UnimplementedPrintServiceServer) mustEmbedUnimplementedPrintServiceServer() {}

// UnsafePrintServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PrintServiceServer will
// result in compilation errors.
type UnsafePrintServiceServer interface {
	mustEmbedUnimplementedPrintServiceServer()
}

func RegisterPrintServiceServer(s grpc.ServiceRegistrar, srv PrintServiceServer) {
	s.RegisterService(&PrintService_ServiceDesc, srv)
}

func _PrintService_PrintGreeting_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmptyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PrintServiceServer).PrintGreeting(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PrintService_PrintGreeting_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PrintServiceServer).PrintGreeting(ctx, req.(*EmptyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PrintService_ServiceDesc is the grpc.ServiceDesc for PrintService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PrintService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "printsvc.PrintService",
	HandlerType: (*PrintServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PrintGreeting",
			Handler:    _PrintService_PrintGreeting_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "contract/printsvc.proto",
}
