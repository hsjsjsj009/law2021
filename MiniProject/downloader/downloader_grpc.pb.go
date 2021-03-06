// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package downloader

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

// MiniProjectServiceClient is the client API for MiniProjectService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MiniProjectServiceClient interface {
	Download(ctx context.Context, in *DownloadRequest, opts ...grpc.CallOption) (*DownloadResponse, error)
}

type miniProjectServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMiniProjectServiceClient(cc grpc.ClientConnInterface) MiniProjectServiceClient {
	return &miniProjectServiceClient{cc}
}

func (c *miniProjectServiceClient) Download(ctx context.Context, in *DownloadRequest, opts ...grpc.CallOption) (*DownloadResponse, error) {
	out := new(DownloadResponse)
	err := c.cc.Invoke(ctx, "/downloader.MiniProjectService/Download", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MiniProjectServiceServer is the server API for MiniProjectService service.
// All implementations must embed UnimplementedMiniProjectServiceServer
// for forward compatibility
type MiniProjectServiceServer interface {
	Download(context.Context, *DownloadRequest) (*DownloadResponse, error)
	mustEmbedUnimplementedMiniProjectServiceServer()
}

// UnimplementedMiniProjectServiceServer must be embedded to have forward compatible implementations.
type UnimplementedMiniProjectServiceServer struct {
}

func (UnimplementedMiniProjectServiceServer) Download(context.Context, *DownloadRequest) (*DownloadResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Download not implemented")
}
func (UnimplementedMiniProjectServiceServer) mustEmbedUnimplementedMiniProjectServiceServer() {}

// UnsafeMiniProjectServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MiniProjectServiceServer will
// result in compilation errors.
type UnsafeMiniProjectServiceServer interface {
	mustEmbedUnimplementedMiniProjectServiceServer()
}

func RegisterMiniProjectServiceServer(s grpc.ServiceRegistrar, srv MiniProjectServiceServer) {
	s.RegisterService(&MiniProjectService_ServiceDesc, srv)
}

func _MiniProjectService_Download_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DownloadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MiniProjectServiceServer).Download(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/downloader.MiniProjectService/Download",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MiniProjectServiceServer).Download(ctx, req.(*DownloadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// MiniProjectService_ServiceDesc is the grpc.ServiceDesc for MiniProjectService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MiniProjectService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "downloader.MiniProjectService",
	HandlerType: (*MiniProjectServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Download",
			Handler:    _MiniProjectService_Download_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "downloader/downloader.proto",
}
