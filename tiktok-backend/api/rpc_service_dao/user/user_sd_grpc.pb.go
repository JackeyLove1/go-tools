// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.4
// source: user_sd.proto

package user

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// UserDaoInfoClient is the client API for UserDaoInfo service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserDaoInfoClient interface {
	AddUser(ctx context.Context, in *UserDaoPost, opts ...grpc.CallOption) (*wrapperspb.BoolValue, error)
	GetUserInfoByUserName(ctx context.Context, in *UserDaoPost, opts ...grpc.CallOption) (*UserDaoInfoResp, error)
	GetUserInfoByUserId(ctx context.Context, in *UserDaoPost, opts ...grpc.CallOption) (*UserDaoInfoResp, error)
	GetUserInfoByUserNameAndPassword(ctx context.Context, in *UserDaoPost, opts ...grpc.CallOption) (*UserDaoInfoResp, error)
}

type userDaoInfoClient struct {
	cc grpc.ClientConnInterface
}

func NewUserDaoInfoClient(cc grpc.ClientConnInterface) UserDaoInfoClient {
	return &userDaoInfoClient{cc}
}

func (c *userDaoInfoClient) AddUser(ctx context.Context, in *UserDaoPost, opts ...grpc.CallOption) (*wrapperspb.BoolValue, error) {
	out := new(wrapperspb.BoolValue)
	err := c.cc.Invoke(ctx, "/user.UserDaoInfo/addUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userDaoInfoClient) GetUserInfoByUserName(ctx context.Context, in *UserDaoPost, opts ...grpc.CallOption) (*UserDaoInfoResp, error) {
	out := new(UserDaoInfoResp)
	err := c.cc.Invoke(ctx, "/user.UserDaoInfo/getUserInfoByUserName", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userDaoInfoClient) GetUserInfoByUserId(ctx context.Context, in *UserDaoPost, opts ...grpc.CallOption) (*UserDaoInfoResp, error) {
	out := new(UserDaoInfoResp)
	err := c.cc.Invoke(ctx, "/user.UserDaoInfo/getUserInfoByUserId", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userDaoInfoClient) GetUserInfoByUserNameAndPassword(ctx context.Context, in *UserDaoPost, opts ...grpc.CallOption) (*UserDaoInfoResp, error) {
	out := new(UserDaoInfoResp)
	err := c.cc.Invoke(ctx, "/user.UserDaoInfo/getUserInfoByUserNameAndPassword", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserDaoInfoServer is the server API for UserDaoInfo service.
// All implementations must embed UnimplementedUserDaoInfoServer
// for forward compatibility
type UserDaoInfoServer interface {
	AddUser(context.Context, *UserDaoPost) (*wrapperspb.BoolValue, error)
	GetUserInfoByUserName(context.Context, *UserDaoPost) (*UserDaoInfoResp, error)
	GetUserInfoByUserId(context.Context, *UserDaoPost) (*UserDaoInfoResp, error)
	GetUserInfoByUserNameAndPassword(context.Context, *UserDaoPost) (*UserDaoInfoResp, error)
	mustEmbedUnimplementedUserDaoInfoServer()
}

// UnimplementedUserDaoInfoServer must be embedded to have forward compatible implementations.
type UnimplementedUserDaoInfoServer struct {
}

func (UnimplementedUserDaoInfoServer) AddUser(context.Context, *UserDaoPost) (*wrapperspb.BoolValue, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddUser not implemented")
}
func (UnimplementedUserDaoInfoServer) GetUserInfoByUserName(context.Context, *UserDaoPost) (*UserDaoInfoResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserInfoByUserName not implemented")
}
func (UnimplementedUserDaoInfoServer) GetUserInfoByUserId(context.Context, *UserDaoPost) (*UserDaoInfoResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserInfoByUserId not implemented")
}
func (UnimplementedUserDaoInfoServer) GetUserInfoByUserNameAndPassword(context.Context, *UserDaoPost) (*UserDaoInfoResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserInfoByUserNameAndPassword not implemented")
}
func (UnimplementedUserDaoInfoServer) mustEmbedUnimplementedUserDaoInfoServer() {}

// UnsafeUserDaoInfoServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserDaoInfoServer will
// result in compilation errors.
type UnsafeUserDaoInfoServer interface {
	mustEmbedUnimplementedUserDaoInfoServer()
}

func RegisterUserDaoInfoServer(s grpc.ServiceRegistrar, srv UserDaoInfoServer) {
	s.RegisterService(&UserDaoInfo_ServiceDesc, srv)
}

func _UserDaoInfo_AddUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserDaoPost)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserDaoInfoServer).AddUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.UserDaoInfo/addUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserDaoInfoServer).AddUser(ctx, req.(*UserDaoPost))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserDaoInfo_GetUserInfoByUserName_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserDaoPost)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserDaoInfoServer).GetUserInfoByUserName(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.UserDaoInfo/getUserInfoByUserName",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserDaoInfoServer).GetUserInfoByUserName(ctx, req.(*UserDaoPost))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserDaoInfo_GetUserInfoByUserId_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserDaoPost)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserDaoInfoServer).GetUserInfoByUserId(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.UserDaoInfo/getUserInfoByUserId",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserDaoInfoServer).GetUserInfoByUserId(ctx, req.(*UserDaoPost))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserDaoInfo_GetUserInfoByUserNameAndPassword_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserDaoPost)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserDaoInfoServer).GetUserInfoByUserNameAndPassword(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.UserDaoInfo/getUserInfoByUserNameAndPassword",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserDaoInfoServer).GetUserInfoByUserNameAndPassword(ctx, req.(*UserDaoPost))
	}
	return interceptor(ctx, in, info, handler)
}

// UserDaoInfo_ServiceDesc is the grpc.ServiceDesc for UserDaoInfo service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UserDaoInfo_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "user.UserDaoInfo",
	HandlerType: (*UserDaoInfoServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "addUser",
			Handler:    _UserDaoInfo_AddUser_Handler,
		},
		{
			MethodName: "getUserInfoByUserName",
			Handler:    _UserDaoInfo_GetUserInfoByUserName_Handler,
		},
		{
			MethodName: "getUserInfoByUserId",
			Handler:    _UserDaoInfo_GetUserInfoByUserId_Handler,
		},
		{
			MethodName: "getUserInfoByUserNameAndPassword",
			Handler:    _UserDaoInfo_GetUserInfoByUserNameAndPassword_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "user_sd.proto",
}
