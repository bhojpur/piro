// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package v1

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

// PiroServiceClient is the client API for PiroService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PiroServiceClient interface {
	// StartLocalJob starts a Kubernetes Job by uploading Bhojpur.NET Platform
	// application content directly. The incoming requests are expected in the
	// following order:
	//   1. metadata
	//   2. all bytes constituting the piro/config.yaml
	//   3. all bytes constituting the job YAML that will be executed (that the config.yaml points to)
	//   4. all bytes constituting the gzipped Bhojpur.NET Platform application tar stream
	//   5. the Bhojpur.NET Platform application tar stream done marker
	StartLocalJob(ctx context.Context, opts ...grpc.CallOption) (PiroService_StartLocalJobClient, error)
	// StartGitHubJob starts a job on a Git context, possibly with a custom job.
	StartGitHubJob(ctx context.Context, in *StartGitHubJobRequest, opts ...grpc.CallOption) (*StartJobResponse, error)
	// StartFromPreviousJob starts a new job based on a previous one.
	// If the previous job does not have the can-replay condition set this call will result in an error.
	StartFromPreviousJob(ctx context.Context, in *StartFromPreviousJobRequest, opts ...grpc.CallOption) (*StartJobResponse, error)
	// StartJobRequest starts a new job based on its specification.
	StartJob(ctx context.Context, in *StartJobRequest, opts ...grpc.CallOption) (*StartJobResponse, error)
	// StartJob2 starts a new job based on its specification.
	StartJob2(ctx context.Context, in *StartJobRequest2, opts ...grpc.CallOption) (*StartJobResponse, error)
	// Searches for jobs known to this instance
	ListJobs(ctx context.Context, in *ListJobsRequest, opts ...grpc.CallOption) (*ListJobsResponse, error)
	// Subscribe listens to new jobs/job updates
	Subscribe(ctx context.Context, in *SubscribeRequest, opts ...grpc.CallOption) (PiroService_SubscribeClient, error)
	// GetJob retrieves details of a single job
	GetJob(ctx context.Context, in *GetJobRequest, opts ...grpc.CallOption) (*GetJobResponse, error)
	// Listen listens to job updates and log output of a running job
	Listen(ctx context.Context, in *ListenRequest, opts ...grpc.CallOption) (PiroService_ListenClient, error)
	// StopJob stops a currently running job
	StopJob(ctx context.Context, in *StopJobRequest, opts ...grpc.CallOption) (*StopJobResponse, error)
}

type piroServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPiroServiceClient(cc grpc.ClientConnInterface) PiroServiceClient {
	return &piroServiceClient{cc}
}

func (c *piroServiceClient) StartLocalJob(ctx context.Context, opts ...grpc.CallOption) (PiroService_StartLocalJobClient, error) {
	stream, err := c.cc.NewStream(ctx, &PiroService_ServiceDesc.Streams[0], "/v1.PiroService/StartLocalJob", opts...)
	if err != nil {
		return nil, err
	}
	x := &piroServiceStartLocalJobClient{stream}
	return x, nil
}

type PiroService_StartLocalJobClient interface {
	Send(*StartLocalJobRequest) error
	CloseAndRecv() (*StartJobResponse, error)
	grpc.ClientStream
}

type piroServiceStartLocalJobClient struct {
	grpc.ClientStream
}

func (x *piroServiceStartLocalJobClient) Send(m *StartLocalJobRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *piroServiceStartLocalJobClient) CloseAndRecv() (*StartJobResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(StartJobResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *piroServiceClient) StartGitHubJob(ctx context.Context, in *StartGitHubJobRequest, opts ...grpc.CallOption) (*StartJobResponse, error) {
	out := new(StartJobResponse)
	err := c.cc.Invoke(ctx, "/v1.PiroService/StartGitHubJob", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *piroServiceClient) StartFromPreviousJob(ctx context.Context, in *StartFromPreviousJobRequest, opts ...grpc.CallOption) (*StartJobResponse, error) {
	out := new(StartJobResponse)
	err := c.cc.Invoke(ctx, "/v1.PiroService/StartFromPreviousJob", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *piroServiceClient) StartJob(ctx context.Context, in *StartJobRequest, opts ...grpc.CallOption) (*StartJobResponse, error) {
	out := new(StartJobResponse)
	err := c.cc.Invoke(ctx, "/v1.PiroService/StartJob", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *piroServiceClient) StartJob2(ctx context.Context, in *StartJobRequest2, opts ...grpc.CallOption) (*StartJobResponse, error) {
	out := new(StartJobResponse)
	err := c.cc.Invoke(ctx, "/v1.PiroService/StartJob2", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *piroServiceClient) ListJobs(ctx context.Context, in *ListJobsRequest, opts ...grpc.CallOption) (*ListJobsResponse, error) {
	out := new(ListJobsResponse)
	err := c.cc.Invoke(ctx, "/v1.PiroService/ListJobs", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *piroServiceClient) Subscribe(ctx context.Context, in *SubscribeRequest, opts ...grpc.CallOption) (PiroService_SubscribeClient, error) {
	stream, err := c.cc.NewStream(ctx, &PiroService_ServiceDesc.Streams[1], "/v1.PiroService/Subscribe", opts...)
	if err != nil {
		return nil, err
	}
	x := &piroServiceSubscribeClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type PiroService_SubscribeClient interface {
	Recv() (*SubscribeResponse, error)
	grpc.ClientStream
}

type piroServiceSubscribeClient struct {
	grpc.ClientStream
}

func (x *piroServiceSubscribeClient) Recv() (*SubscribeResponse, error) {
	m := new(SubscribeResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *piroServiceClient) GetJob(ctx context.Context, in *GetJobRequest, opts ...grpc.CallOption) (*GetJobResponse, error) {
	out := new(GetJobResponse)
	err := c.cc.Invoke(ctx, "/v1.PiroService/GetJob", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *piroServiceClient) Listen(ctx context.Context, in *ListenRequest, opts ...grpc.CallOption) (PiroService_ListenClient, error) {
	stream, err := c.cc.NewStream(ctx, &PiroService_ServiceDesc.Streams[2], "/v1.PiroService/Listen", opts...)
	if err != nil {
		return nil, err
	}
	x := &piroServiceListenClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type PiroService_ListenClient interface {
	Recv() (*ListenResponse, error)
	grpc.ClientStream
}

type piroServiceListenClient struct {
	grpc.ClientStream
}

func (x *piroServiceListenClient) Recv() (*ListenResponse, error) {
	m := new(ListenResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *piroServiceClient) StopJob(ctx context.Context, in *StopJobRequest, opts ...grpc.CallOption) (*StopJobResponse, error) {
	out := new(StopJobResponse)
	err := c.cc.Invoke(ctx, "/v1.PiroService/StopJob", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PiroServiceServer is the server API for PiroService service.
// All implementations must embed UnimplementedPiroServiceServer
// for forward compatibility
type PiroServiceServer interface {
	// StartLocalJob starts a Kubernetes Job by uploading Bhojpur.NET Platform
	// application content directly. The incoming requests are expected in the
	// following order:
	//   1. metadata
	//   2. all bytes constituting the piro/config.yaml
	//   3. all bytes constituting the job YAML that will be executed (that the config.yaml points to)
	//   4. all bytes constituting the gzipped Bhojpur.NET Platform application tar stream
	//   5. the Bhojpur.NET Platform application tar stream done marker
	StartLocalJob(PiroService_StartLocalJobServer) error
	// StartGitHubJob starts a job on a Git context, possibly with a custom job.
	StartGitHubJob(context.Context, *StartGitHubJobRequest) (*StartJobResponse, error)
	// StartFromPreviousJob starts a new job based on a previous one.
	// If the previous job does not have the can-replay condition set this call will result in an error.
	StartFromPreviousJob(context.Context, *StartFromPreviousJobRequest) (*StartJobResponse, error)
	// StartJobRequest starts a new job based on its specification.
	StartJob(context.Context, *StartJobRequest) (*StartJobResponse, error)
	// StartJob2 starts a new job based on its specification.
	StartJob2(context.Context, *StartJobRequest2) (*StartJobResponse, error)
	// Searches for jobs known to this instance
	ListJobs(context.Context, *ListJobsRequest) (*ListJobsResponse, error)
	// Subscribe listens to new jobs/job updates
	Subscribe(*SubscribeRequest, PiroService_SubscribeServer) error
	// GetJob retrieves details of a single job
	GetJob(context.Context, *GetJobRequest) (*GetJobResponse, error)
	// Listen listens to job updates and log output of a running job
	Listen(*ListenRequest, PiroService_ListenServer) error
	// StopJob stops a currently running job
	StopJob(context.Context, *StopJobRequest) (*StopJobResponse, error)
	mustEmbedUnimplementedPiroServiceServer()
}

// UnimplementedPiroServiceServer must be embedded to have forward compatible implementations.
type UnimplementedPiroServiceServer struct {
}

func (UnimplementedPiroServiceServer) StartLocalJob(PiroService_StartLocalJobServer) error {
	return status.Errorf(codes.Unimplemented, "method StartLocalJob not implemented")
}
func (UnimplementedPiroServiceServer) StartGitHubJob(context.Context, *StartGitHubJobRequest) (*StartJobResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StartGitHubJob not implemented")
}
func (UnimplementedPiroServiceServer) StartFromPreviousJob(context.Context, *StartFromPreviousJobRequest) (*StartJobResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StartFromPreviousJob not implemented")
}
func (UnimplementedPiroServiceServer) StartJob(context.Context, *StartJobRequest) (*StartJobResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StartJob not implemented")
}
func (UnimplementedPiroServiceServer) StartJob2(context.Context, *StartJobRequest2) (*StartJobResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StartJob2 not implemented")
}
func (UnimplementedPiroServiceServer) ListJobs(context.Context, *ListJobsRequest) (*ListJobsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListJobs not implemented")
}
func (UnimplementedPiroServiceServer) Subscribe(*SubscribeRequest, PiroService_SubscribeServer) error {
	return status.Errorf(codes.Unimplemented, "method Subscribe not implemented")
}
func (UnimplementedPiroServiceServer) GetJob(context.Context, *GetJobRequest) (*GetJobResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetJob not implemented")
}
func (UnimplementedPiroServiceServer) Listen(*ListenRequest, PiroService_ListenServer) error {
	return status.Errorf(codes.Unimplemented, "method Listen not implemented")
}
func (UnimplementedPiroServiceServer) StopJob(context.Context, *StopJobRequest) (*StopJobResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StopJob not implemented")
}
func (UnimplementedPiroServiceServer) mustEmbedUnimplementedPiroServiceServer() {}

// UnsafePiroServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PiroServiceServer will
// result in compilation errors.
type UnsafePiroServiceServer interface {
	mustEmbedUnimplementedPiroServiceServer()
}

func RegisterPiroServiceServer(s grpc.ServiceRegistrar, srv PiroServiceServer) {
	s.RegisterService(&PiroService_ServiceDesc, srv)
}

func _PiroService_StartLocalJob_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(PiroServiceServer).StartLocalJob(&piroServiceStartLocalJobServer{stream})
}

type PiroService_StartLocalJobServer interface {
	SendAndClose(*StartJobResponse) error
	Recv() (*StartLocalJobRequest, error)
	grpc.ServerStream
}

type piroServiceStartLocalJobServer struct {
	grpc.ServerStream
}

func (x *piroServiceStartLocalJobServer) SendAndClose(m *StartJobResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *piroServiceStartLocalJobServer) Recv() (*StartLocalJobRequest, error) {
	m := new(StartLocalJobRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _PiroService_StartGitHubJob_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StartGitHubJobRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PiroServiceServer).StartGitHubJob(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.PiroService/StartGitHubJob",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PiroServiceServer).StartGitHubJob(ctx, req.(*StartGitHubJobRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PiroService_StartFromPreviousJob_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StartFromPreviousJobRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PiroServiceServer).StartFromPreviousJob(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.PiroService/StartFromPreviousJob",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PiroServiceServer).StartFromPreviousJob(ctx, req.(*StartFromPreviousJobRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PiroService_StartJob_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StartJobRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PiroServiceServer).StartJob(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.PiroService/StartJob",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PiroServiceServer).StartJob(ctx, req.(*StartJobRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PiroService_StartJob2_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StartJobRequest2)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PiroServiceServer).StartJob2(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.PiroService/StartJob2",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PiroServiceServer).StartJob2(ctx, req.(*StartJobRequest2))
	}
	return interceptor(ctx, in, info, handler)
}

func _PiroService_ListJobs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListJobsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PiroServiceServer).ListJobs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.PiroService/ListJobs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PiroServiceServer).ListJobs(ctx, req.(*ListJobsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PiroService_Subscribe_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SubscribeRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(PiroServiceServer).Subscribe(m, &piroServiceSubscribeServer{stream})
}

type PiroService_SubscribeServer interface {
	Send(*SubscribeResponse) error
	grpc.ServerStream
}

type piroServiceSubscribeServer struct {
	grpc.ServerStream
}

func (x *piroServiceSubscribeServer) Send(m *SubscribeResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _PiroService_GetJob_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetJobRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PiroServiceServer).GetJob(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.PiroService/GetJob",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PiroServiceServer).GetJob(ctx, req.(*GetJobRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PiroService_Listen_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ListenRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(PiroServiceServer).Listen(m, &piroServiceListenServer{stream})
}

type PiroService_ListenServer interface {
	Send(*ListenResponse) error
	grpc.ServerStream
}

type piroServiceListenServer struct {
	grpc.ServerStream
}

func (x *piroServiceListenServer) Send(m *ListenResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _PiroService_StopJob_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StopJobRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PiroServiceServer).StopJob(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.PiroService/StopJob",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PiroServiceServer).StopJob(ctx, req.(*StopJobRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PiroService_ServiceDesc is the grpc.ServiceDesc for PiroService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PiroService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "v1.PiroService",
	HandlerType: (*PiroServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "StartGitHubJob",
			Handler:    _PiroService_StartGitHubJob_Handler,
		},
		{
			MethodName: "StartFromPreviousJob",
			Handler:    _PiroService_StartFromPreviousJob_Handler,
		},
		{
			MethodName: "StartJob",
			Handler:    _PiroService_StartJob_Handler,
		},
		{
			MethodName: "StartJob2",
			Handler:    _PiroService_StartJob2_Handler,
		},
		{
			MethodName: "ListJobs",
			Handler:    _PiroService_ListJobs_Handler,
		},
		{
			MethodName: "GetJob",
			Handler:    _PiroService_GetJob_Handler,
		},
		{
			MethodName: "StopJob",
			Handler:    _PiroService_StopJob_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StartLocalJob",
			Handler:       _PiroService_StartLocalJob_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "Subscribe",
			Handler:       _PiroService_Subscribe_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "Listen",
			Handler:       _PiroService_Listen_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "piro.proto",
}
