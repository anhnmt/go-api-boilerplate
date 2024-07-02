// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             (unknown)
// source: auth.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.62.0 or later.
const _ = grpc.SupportPackageIsVersion8

const (
	AuthService_Login_FullMethodName          = "/auth.v1.AuthService/Login"
	AuthService_Info_FullMethodName           = "/auth.v1.AuthService/Info"
	AuthService_RevokeToken_FullMethodName    = "/auth.v1.AuthService/RevokeToken"
	AuthService_RefreshToken_FullMethodName   = "/auth.v1.AuthService/RefreshToken"
	AuthService_ActiveSessions_FullMethodName = "/auth.v1.AuthService/ActiveSessions"
)

// AuthServiceClient is the client API for AuthService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AuthServiceClient interface {
	Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error)
	Info(ctx context.Context, in *InfoRequest, opts ...grpc.CallOption) (*InfoResponse, error)
	RevokeToken(ctx context.Context, in *RevokeTokenRequest, opts ...grpc.CallOption) (*RevokeTokenResponse, error)
	RefreshToken(ctx context.Context, in *RefreshTokenRequest, opts ...grpc.CallOption) (*RefreshTokenResponse, error)
	ActiveSessions(ctx context.Context, in *ActiveSessionsRequest, opts ...grpc.CallOption) (*ActiveSessionsResponse, error)
}

type authServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAuthServiceClient(cc grpc.ClientConnInterface) AuthServiceClient {
	return &authServiceClient{cc}
}

func (c *authServiceClient) Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LoginResponse)
	err := c.cc.Invoke(ctx, AuthService_Login_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) Info(ctx context.Context, in *InfoRequest, opts ...grpc.CallOption) (*InfoResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(InfoResponse)
	err := c.cc.Invoke(ctx, AuthService_Info_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) RevokeToken(ctx context.Context, in *RevokeTokenRequest, opts ...grpc.CallOption) (*RevokeTokenResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RevokeTokenResponse)
	err := c.cc.Invoke(ctx, AuthService_RevokeToken_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) RefreshToken(ctx context.Context, in *RefreshTokenRequest, opts ...grpc.CallOption) (*RefreshTokenResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RefreshTokenResponse)
	err := c.cc.Invoke(ctx, AuthService_RefreshToken_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) ActiveSessions(ctx context.Context, in *ActiveSessionsRequest, opts ...grpc.CallOption) (*ActiveSessionsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ActiveSessionsResponse)
	err := c.cc.Invoke(ctx, AuthService_ActiveSessions_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuthServiceServer is the server API for AuthService service.
// All implementations must embed UnimplementedAuthServiceServer
// for forward compatibility
type AuthServiceServer interface {
	Login(context.Context, *LoginRequest) (*LoginResponse, error)
	Info(context.Context, *InfoRequest) (*InfoResponse, error)
	RevokeToken(context.Context, *RevokeTokenRequest) (*RevokeTokenResponse, error)
	RefreshToken(context.Context, *RefreshTokenRequest) (*RefreshTokenResponse, error)
	ActiveSessions(context.Context, *ActiveSessionsRequest) (*ActiveSessionsResponse, error)
	mustEmbedUnimplementedAuthServiceServer()
}

// UnimplementedAuthServiceServer must be embedded to have forward compatible implementations.
type UnimplementedAuthServiceServer struct {
}

func (UnimplementedAuthServiceServer) Login(context.Context, *LoginRequest) (*LoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedAuthServiceServer) Info(context.Context, *InfoRequest) (*InfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Info not implemented")
}
func (UnimplementedAuthServiceServer) RevokeToken(context.Context, *RevokeTokenRequest) (*RevokeTokenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RevokeToken not implemented")
}
func (UnimplementedAuthServiceServer) RefreshToken(context.Context, *RefreshTokenRequest) (*RefreshTokenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RefreshToken not implemented")
}
func (UnimplementedAuthServiceServer) ActiveSessions(context.Context, *ActiveSessionsRequest) (*ActiveSessionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ActiveSessions not implemented")
}
func (UnimplementedAuthServiceServer) mustEmbedUnimplementedAuthServiceServer() {}

// UnsafeAuthServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AuthServiceServer will
// result in compilation errors.
type UnsafeAuthServiceServer interface {
	mustEmbedUnimplementedAuthServiceServer()
}

func RegisterAuthServiceServer(s grpc.ServiceRegistrar, srv AuthServiceServer) {
	s.RegisterService(&AuthService_ServiceDesc, srv)
}

func _AuthService_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuthService_Login_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).Login(ctx, req.(*LoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_Info_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).Info(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuthService_Info_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).Info(ctx, req.(*InfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_RevokeToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RevokeTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).RevokeToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuthService_RevokeToken_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).RevokeToken(ctx, req.(*RevokeTokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_RefreshToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RefreshTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).RefreshToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuthService_RefreshToken_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).RefreshToken(ctx, req.(*RefreshTokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_ActiveSessions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ActiveSessionsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).ActiveSessions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuthService_ActiveSessions_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).ActiveSessions(ctx, req.(*ActiveSessionsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AuthService_ServiceDesc is the grpc.ServiceDesc for AuthService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AuthService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "auth.v1.AuthService",
	HandlerType: (*AuthServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Login",
			Handler:    _AuthService_Login_Handler,
		},
		{
			MethodName: "Info",
			Handler:    _AuthService_Info_Handler,
		},
		{
			MethodName: "RevokeToken",
			Handler:    _AuthService_RevokeToken_Handler,
		},
		{
			MethodName: "RefreshToken",
			Handler:    _AuthService_RefreshToken_Handler,
		},
		{
			MethodName: "ActiveSessions",
			Handler:    _AuthService_ActiveSessions_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "auth.proto",
}
