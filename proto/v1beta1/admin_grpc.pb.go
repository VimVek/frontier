// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: raystack/frontier/v1beta1/admin.proto

package frontierv1beta1

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
	AdminService_ListAllUsers_FullMethodName         = "/raystack.frontier.v1beta1.AdminService/ListAllUsers"
	AdminService_ListGroups_FullMethodName           = "/raystack.frontier.v1beta1.AdminService/ListGroups"
	AdminService_ListAllOrganizations_FullMethodName = "/raystack.frontier.v1beta1.AdminService/ListAllOrganizations"
	AdminService_ListProjects_FullMethodName         = "/raystack.frontier.v1beta1.AdminService/ListProjects"
	AdminService_ListRelations_FullMethodName        = "/raystack.frontier.v1beta1.AdminService/ListRelations"
	AdminService_ListResources_FullMethodName        = "/raystack.frontier.v1beta1.AdminService/ListResources"
	AdminService_ListPolicies_FullMethodName         = "/raystack.frontier.v1beta1.AdminService/ListPolicies"
	AdminService_CreateRole_FullMethodName           = "/raystack.frontier.v1beta1.AdminService/CreateRole"
	AdminService_UpdateRole_FullMethodName           = "/raystack.frontier.v1beta1.AdminService/UpdateRole"
	AdminService_DeleteRole_FullMethodName           = "/raystack.frontier.v1beta1.AdminService/DeleteRole"
	AdminService_CreatePermission_FullMethodName     = "/raystack.frontier.v1beta1.AdminService/CreatePermission"
	AdminService_UpdatePermission_FullMethodName     = "/raystack.frontier.v1beta1.AdminService/UpdatePermission"
	AdminService_DeletePermission_FullMethodName     = "/raystack.frontier.v1beta1.AdminService/DeletePermission"
	AdminService_ListPreferences_FullMethodName      = "/raystack.frontier.v1beta1.AdminService/ListPreferences"
	AdminService_CreatePreferences_FullMethodName    = "/raystack.frontier.v1beta1.AdminService/CreatePreferences"
)

// AdminServiceClient is the client API for AdminService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AdminServiceClient interface {
	// Users
	ListAllUsers(ctx context.Context, in *ListAllUsersRequest, opts ...grpc.CallOption) (*ListAllUsersResponse, error)
	// Group
	ListGroups(ctx context.Context, in *ListGroupsRequest, opts ...grpc.CallOption) (*ListGroupsResponse, error)
	// Organizations
	ListAllOrganizations(ctx context.Context, in *ListAllOrganizationsRequest, opts ...grpc.CallOption) (*ListAllOrganizationsResponse, error)
	// Projects
	ListProjects(ctx context.Context, in *ListProjectsRequest, opts ...grpc.CallOption) (*ListProjectsResponse, error)
	// Relations
	ListRelations(ctx context.Context, in *ListRelationsRequest, opts ...grpc.CallOption) (*ListRelationsResponse, error)
	// Resources
	ListResources(ctx context.Context, in *ListResourcesRequest, opts ...grpc.CallOption) (*ListResourcesResponse, error)
	// Policies
	ListPolicies(ctx context.Context, in *ListPoliciesRequest, opts ...grpc.CallOption) (*ListPoliciesResponse, error)
	// Roles
	CreateRole(ctx context.Context, in *CreateRoleRequest, opts ...grpc.CallOption) (*CreateRoleResponse, error)
	UpdateRole(ctx context.Context, in *UpdateRoleRequest, opts ...grpc.CallOption) (*UpdateRoleResponse, error)
	DeleteRole(ctx context.Context, in *DeleteRoleRequest, opts ...grpc.CallOption) (*DeleteRoleResponse, error)
	// Permissions
	CreatePermission(ctx context.Context, in *CreatePermissionRequest, opts ...grpc.CallOption) (*CreatePermissionResponse, error)
	UpdatePermission(ctx context.Context, in *UpdatePermissionRequest, opts ...grpc.CallOption) (*UpdatePermissionResponse, error)
	DeletePermission(ctx context.Context, in *DeletePermissionRequest, opts ...grpc.CallOption) (*DeletePermissionResponse, error)
	// Preferences
	ListPreferences(ctx context.Context, in *ListPreferencesRequest, opts ...grpc.CallOption) (*ListPreferencesResponse, error)
	CreatePreferences(ctx context.Context, in *CreatePreferencesRequest, opts ...grpc.CallOption) (*CreatePreferencesResponse, error)
}

type adminServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAdminServiceClient(cc grpc.ClientConnInterface) AdminServiceClient {
	return &adminServiceClient{cc}
}

func (c *adminServiceClient) ListAllUsers(ctx context.Context, in *ListAllUsersRequest, opts ...grpc.CallOption) (*ListAllUsersResponse, error) {
	out := new(ListAllUsersResponse)
	err := c.cc.Invoke(ctx, AdminService_ListAllUsers_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminServiceClient) ListGroups(ctx context.Context, in *ListGroupsRequest, opts ...grpc.CallOption) (*ListGroupsResponse, error) {
	out := new(ListGroupsResponse)
	err := c.cc.Invoke(ctx, AdminService_ListGroups_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminServiceClient) ListAllOrganizations(ctx context.Context, in *ListAllOrganizationsRequest, opts ...grpc.CallOption) (*ListAllOrganizationsResponse, error) {
	out := new(ListAllOrganizationsResponse)
	err := c.cc.Invoke(ctx, AdminService_ListAllOrganizations_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminServiceClient) ListProjects(ctx context.Context, in *ListProjectsRequest, opts ...grpc.CallOption) (*ListProjectsResponse, error) {
	out := new(ListProjectsResponse)
	err := c.cc.Invoke(ctx, AdminService_ListProjects_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminServiceClient) ListRelations(ctx context.Context, in *ListRelationsRequest, opts ...grpc.CallOption) (*ListRelationsResponse, error) {
	out := new(ListRelationsResponse)
	err := c.cc.Invoke(ctx, AdminService_ListRelations_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminServiceClient) ListResources(ctx context.Context, in *ListResourcesRequest, opts ...grpc.CallOption) (*ListResourcesResponse, error) {
	out := new(ListResourcesResponse)
	err := c.cc.Invoke(ctx, AdminService_ListResources_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminServiceClient) ListPolicies(ctx context.Context, in *ListPoliciesRequest, opts ...grpc.CallOption) (*ListPoliciesResponse, error) {
	out := new(ListPoliciesResponse)
	err := c.cc.Invoke(ctx, AdminService_ListPolicies_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminServiceClient) CreateRole(ctx context.Context, in *CreateRoleRequest, opts ...grpc.CallOption) (*CreateRoleResponse, error) {
	out := new(CreateRoleResponse)
	err := c.cc.Invoke(ctx, AdminService_CreateRole_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminServiceClient) UpdateRole(ctx context.Context, in *UpdateRoleRequest, opts ...grpc.CallOption) (*UpdateRoleResponse, error) {
	out := new(UpdateRoleResponse)
	err := c.cc.Invoke(ctx, AdminService_UpdateRole_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminServiceClient) DeleteRole(ctx context.Context, in *DeleteRoleRequest, opts ...grpc.CallOption) (*DeleteRoleResponse, error) {
	out := new(DeleteRoleResponse)
	err := c.cc.Invoke(ctx, AdminService_DeleteRole_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminServiceClient) CreatePermission(ctx context.Context, in *CreatePermissionRequest, opts ...grpc.CallOption) (*CreatePermissionResponse, error) {
	out := new(CreatePermissionResponse)
	err := c.cc.Invoke(ctx, AdminService_CreatePermission_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminServiceClient) UpdatePermission(ctx context.Context, in *UpdatePermissionRequest, opts ...grpc.CallOption) (*UpdatePermissionResponse, error) {
	out := new(UpdatePermissionResponse)
	err := c.cc.Invoke(ctx, AdminService_UpdatePermission_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminServiceClient) DeletePermission(ctx context.Context, in *DeletePermissionRequest, opts ...grpc.CallOption) (*DeletePermissionResponse, error) {
	out := new(DeletePermissionResponse)
	err := c.cc.Invoke(ctx, AdminService_DeletePermission_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminServiceClient) ListPreferences(ctx context.Context, in *ListPreferencesRequest, opts ...grpc.CallOption) (*ListPreferencesResponse, error) {
	out := new(ListPreferencesResponse)
	err := c.cc.Invoke(ctx, AdminService_ListPreferences_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminServiceClient) CreatePreferences(ctx context.Context, in *CreatePreferencesRequest, opts ...grpc.CallOption) (*CreatePreferencesResponse, error) {
	out := new(CreatePreferencesResponse)
	err := c.cc.Invoke(ctx, AdminService_CreatePreferences_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AdminServiceServer is the server API for AdminService service.
// All implementations must embed UnimplementedAdminServiceServer
// for forward compatibility
type AdminServiceServer interface {
	// Users
	ListAllUsers(context.Context, *ListAllUsersRequest) (*ListAllUsersResponse, error)
	// Group
	ListGroups(context.Context, *ListGroupsRequest) (*ListGroupsResponse, error)
	// Organizations
	ListAllOrganizations(context.Context, *ListAllOrganizationsRequest) (*ListAllOrganizationsResponse, error)
	// Projects
	ListProjects(context.Context, *ListProjectsRequest) (*ListProjectsResponse, error)
	// Relations
	ListRelations(context.Context, *ListRelationsRequest) (*ListRelationsResponse, error)
	// Resources
	ListResources(context.Context, *ListResourcesRequest) (*ListResourcesResponse, error)
	// Policies
	ListPolicies(context.Context, *ListPoliciesRequest) (*ListPoliciesResponse, error)
	// Roles
	CreateRole(context.Context, *CreateRoleRequest) (*CreateRoleResponse, error)
	UpdateRole(context.Context, *UpdateRoleRequest) (*UpdateRoleResponse, error)
	DeleteRole(context.Context, *DeleteRoleRequest) (*DeleteRoleResponse, error)
	// Permissions
	CreatePermission(context.Context, *CreatePermissionRequest) (*CreatePermissionResponse, error)
	UpdatePermission(context.Context, *UpdatePermissionRequest) (*UpdatePermissionResponse, error)
	DeletePermission(context.Context, *DeletePermissionRequest) (*DeletePermissionResponse, error)
	// Preferences
	ListPreferences(context.Context, *ListPreferencesRequest) (*ListPreferencesResponse, error)
	CreatePreferences(context.Context, *CreatePreferencesRequest) (*CreatePreferencesResponse, error)
	mustEmbedUnimplementedAdminServiceServer()
}

// UnimplementedAdminServiceServer must be embedded to have forward compatible implementations.
type UnimplementedAdminServiceServer struct {
}

func (UnimplementedAdminServiceServer) ListAllUsers(context.Context, *ListAllUsersRequest) (*ListAllUsersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListAllUsers not implemented")
}
func (UnimplementedAdminServiceServer) ListGroups(context.Context, *ListGroupsRequest) (*ListGroupsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListGroups not implemented")
}
func (UnimplementedAdminServiceServer) ListAllOrganizations(context.Context, *ListAllOrganizationsRequest) (*ListAllOrganizationsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListAllOrganizations not implemented")
}
func (UnimplementedAdminServiceServer) ListProjects(context.Context, *ListProjectsRequest) (*ListProjectsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListProjects not implemented")
}
func (UnimplementedAdminServiceServer) ListRelations(context.Context, *ListRelationsRequest) (*ListRelationsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListRelations not implemented")
}
func (UnimplementedAdminServiceServer) ListResources(context.Context, *ListResourcesRequest) (*ListResourcesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListResources not implemented")
}
func (UnimplementedAdminServiceServer) ListPolicies(context.Context, *ListPoliciesRequest) (*ListPoliciesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListPolicies not implemented")
}
func (UnimplementedAdminServiceServer) CreateRole(context.Context, *CreateRoleRequest) (*CreateRoleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateRole not implemented")
}
func (UnimplementedAdminServiceServer) UpdateRole(context.Context, *UpdateRoleRequest) (*UpdateRoleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateRole not implemented")
}
func (UnimplementedAdminServiceServer) DeleteRole(context.Context, *DeleteRoleRequest) (*DeleteRoleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteRole not implemented")
}
func (UnimplementedAdminServiceServer) CreatePermission(context.Context, *CreatePermissionRequest) (*CreatePermissionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreatePermission not implemented")
}
func (UnimplementedAdminServiceServer) UpdatePermission(context.Context, *UpdatePermissionRequest) (*UpdatePermissionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdatePermission not implemented")
}
func (UnimplementedAdminServiceServer) DeletePermission(context.Context, *DeletePermissionRequest) (*DeletePermissionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeletePermission not implemented")
}
func (UnimplementedAdminServiceServer) ListPreferences(context.Context, *ListPreferencesRequest) (*ListPreferencesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListPreferences not implemented")
}
func (UnimplementedAdminServiceServer) CreatePreferences(context.Context, *CreatePreferencesRequest) (*CreatePreferencesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreatePreferences not implemented")
}
func (UnimplementedAdminServiceServer) mustEmbedUnimplementedAdminServiceServer() {}

// UnsafeAdminServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AdminServiceServer will
// result in compilation errors.
type UnsafeAdminServiceServer interface {
	mustEmbedUnimplementedAdminServiceServer()
}

func RegisterAdminServiceServer(s grpc.ServiceRegistrar, srv AdminServiceServer) {
	s.RegisterService(&AdminService_ServiceDesc, srv)
}

func _AdminService_ListAllUsers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListAllUsersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServiceServer).ListAllUsers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AdminService_ListAllUsers_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServiceServer).ListAllUsers(ctx, req.(*ListAllUsersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdminService_ListGroups_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListGroupsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServiceServer).ListGroups(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AdminService_ListGroups_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServiceServer).ListGroups(ctx, req.(*ListGroupsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdminService_ListAllOrganizations_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListAllOrganizationsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServiceServer).ListAllOrganizations(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AdminService_ListAllOrganizations_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServiceServer).ListAllOrganizations(ctx, req.(*ListAllOrganizationsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdminService_ListProjects_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListProjectsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServiceServer).ListProjects(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AdminService_ListProjects_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServiceServer).ListProjects(ctx, req.(*ListProjectsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdminService_ListRelations_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListRelationsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServiceServer).ListRelations(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AdminService_ListRelations_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServiceServer).ListRelations(ctx, req.(*ListRelationsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdminService_ListResources_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListResourcesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServiceServer).ListResources(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AdminService_ListResources_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServiceServer).ListResources(ctx, req.(*ListResourcesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdminService_ListPolicies_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListPoliciesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServiceServer).ListPolicies(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AdminService_ListPolicies_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServiceServer).ListPolicies(ctx, req.(*ListPoliciesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdminService_CreateRole_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRoleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServiceServer).CreateRole(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AdminService_CreateRole_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServiceServer).CreateRole(ctx, req.(*CreateRoleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdminService_UpdateRole_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateRoleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServiceServer).UpdateRole(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AdminService_UpdateRole_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServiceServer).UpdateRole(ctx, req.(*UpdateRoleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdminService_DeleteRole_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRoleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServiceServer).DeleteRole(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AdminService_DeleteRole_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServiceServer).DeleteRole(ctx, req.(*DeleteRoleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdminService_CreatePermission_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreatePermissionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServiceServer).CreatePermission(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AdminService_CreatePermission_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServiceServer).CreatePermission(ctx, req.(*CreatePermissionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdminService_UpdatePermission_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdatePermissionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServiceServer).UpdatePermission(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AdminService_UpdatePermission_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServiceServer).UpdatePermission(ctx, req.(*UpdatePermissionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdminService_DeletePermission_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeletePermissionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServiceServer).DeletePermission(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AdminService_DeletePermission_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServiceServer).DeletePermission(ctx, req.(*DeletePermissionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdminService_ListPreferences_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListPreferencesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServiceServer).ListPreferences(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AdminService_ListPreferences_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServiceServer).ListPreferences(ctx, req.(*ListPreferencesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdminService_CreatePreferences_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreatePreferencesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServiceServer).CreatePreferences(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AdminService_CreatePreferences_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServiceServer).CreatePreferences(ctx, req.(*CreatePreferencesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AdminService_ServiceDesc is the grpc.ServiceDesc for AdminService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AdminService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "raystack.frontier.v1beta1.AdminService",
	HandlerType: (*AdminServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListAllUsers",
			Handler:    _AdminService_ListAllUsers_Handler,
		},
		{
			MethodName: "ListGroups",
			Handler:    _AdminService_ListGroups_Handler,
		},
		{
			MethodName: "ListAllOrganizations",
			Handler:    _AdminService_ListAllOrganizations_Handler,
		},
		{
			MethodName: "ListProjects",
			Handler:    _AdminService_ListProjects_Handler,
		},
		{
			MethodName: "ListRelations",
			Handler:    _AdminService_ListRelations_Handler,
		},
		{
			MethodName: "ListResources",
			Handler:    _AdminService_ListResources_Handler,
		},
		{
			MethodName: "ListPolicies",
			Handler:    _AdminService_ListPolicies_Handler,
		},
		{
			MethodName: "CreateRole",
			Handler:    _AdminService_CreateRole_Handler,
		},
		{
			MethodName: "UpdateRole",
			Handler:    _AdminService_UpdateRole_Handler,
		},
		{
			MethodName: "DeleteRole",
			Handler:    _AdminService_DeleteRole_Handler,
		},
		{
			MethodName: "CreatePermission",
			Handler:    _AdminService_CreatePermission_Handler,
		},
		{
			MethodName: "UpdatePermission",
			Handler:    _AdminService_UpdatePermission_Handler,
		},
		{
			MethodName: "DeletePermission",
			Handler:    _AdminService_DeletePermission_Handler,
		},
		{
			MethodName: "ListPreferences",
			Handler:    _AdminService_ListPreferences_Handler,
		},
		{
			MethodName: "CreatePreferences",
			Handler:    _AdminService_CreatePreferences_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "raystack/frontier/v1beta1/admin.proto",
}
