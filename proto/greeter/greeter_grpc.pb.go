// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package greeter

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// GreeterClient is the client API for Greeter service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GreeterClient interface {
	Hello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloResponse, error)
	Goodbye(ctx context.Context, in *GoodbyeRequest, opts ...grpc.CallOption) (*GoodbyeResponse, error)
}

type greeterClient struct {
	cc grpc.ClientConnInterface
}

func NewGreeterClient(cc grpc.ClientConnInterface) GreeterClient {
	return &greeterClient{cc}
}

var greeterHelloStreamDesc = &grpc.StreamDesc{
	StreamName: "Hello",
}

func (c *greeterClient) Hello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloResponse, error) {
	out := new(HelloResponse)
	err := c.cc.Invoke(ctx, "/greeter.Greeter/Hello", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

var greeterGoodbyeStreamDesc = &grpc.StreamDesc{
	StreamName: "Goodbye",
}

func (c *greeterClient) Goodbye(ctx context.Context, in *GoodbyeRequest, opts ...grpc.CallOption) (*GoodbyeResponse, error) {
	out := new(GoodbyeResponse)
	err := c.cc.Invoke(ctx, "/greeter.Greeter/Goodbye", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GreeterService is the service API for Greeter service.
// Fields should be assigned to their respective handler implementations only before
// RegisterGreeterService is called.  Any unassigned fields will result in the
// handler for that method returning an Unimplemented error.
type GreeterService struct {
	Hello   func(context.Context, *HelloRequest) (*HelloResponse, error)
	Goodbye func(context.Context, *GoodbyeRequest) (*GoodbyeResponse, error)
}

func (s *GreeterService) hello(_ interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HelloRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return s.Hello(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     s,
		FullMethod: "/greeter.Greeter/Hello",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return s.Hello(ctx, req.(*HelloRequest))
	}
	return interceptor(ctx, in, info, handler)
}
func (s *GreeterService) goodbye(_ interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GoodbyeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return s.Goodbye(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     s,
		FullMethod: "/greeter.Greeter/Goodbye",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return s.Goodbye(ctx, req.(*GoodbyeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RegisterGreeterService registers a service implementation with a gRPC server.
func RegisterGreeterService(s grpc.ServiceRegistrar, srv *GreeterService) {
	srvCopy := *srv
	if srvCopy.Hello == nil {
		srvCopy.Hello = func(context.Context, *HelloRequest) (*HelloResponse, error) {
			return nil, status.Errorf(codes.Unimplemented, "method Hello not implemented")
		}
	}
	if srvCopy.Goodbye == nil {
		srvCopy.Goodbye = func(context.Context, *GoodbyeRequest) (*GoodbyeResponse, error) {
			return nil, status.Errorf(codes.Unimplemented, "method Goodbye not implemented")
		}
	}
	sd := grpc.ServiceDesc{
		ServiceName: "greeter.Greeter",
		Methods: []grpc.MethodDesc{
			{
				MethodName: "Hello",
				Handler:    srvCopy.hello,
			},
			{
				MethodName: "Goodbye",
				Handler:    srvCopy.goodbye,
			},
		},
		Streams:  []grpc.StreamDesc{},
		Metadata: "greeter.proto",
	}

	s.RegisterService(&sd, nil)
}

// NewGreeterService creates a new GreeterService containing the
// implemented methods of the Greeter service in s.  Any unimplemented
// methods will result in the gRPC server returning an UNIMPLEMENTED status to the client.
// This includes situations where the method handler is misspelled or has the wrong
// signature.  For this reason, this function should be used with great care and
// is not recommended to be used by most users.
func NewGreeterService(s interface{}) *GreeterService {
	ns := &GreeterService{}
	if h, ok := s.(interface {
		Hello(context.Context, *HelloRequest) (*HelloResponse, error)
	}); ok {
		ns.Hello = h.Hello
	}
	if h, ok := s.(interface {
		Goodbye(context.Context, *GoodbyeRequest) (*GoodbyeResponse, error)
	}); ok {
		ns.Goodbye = h.Goodbye
	}
	return ns
}

// UnstableGreeterService is the service API for Greeter service.
// New methods may be added to this interface if they are added to the service
// definition, which is not a backward-compatible change.  For this reason,
// use of this type is not recommended.
type UnstableGreeterService interface {
	Hello(context.Context, *HelloRequest) (*HelloResponse, error)
	Goodbye(context.Context, *GoodbyeRequest) (*GoodbyeResponse, error)
}