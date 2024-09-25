// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: initia/gov/v1/query.proto

package govv1

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
	Query_Params_FullMethodName             = "/initia.gov.v1.Query/Params"
	Query_EmergencyProposals_FullMethodName = "/initia.gov.v1.Query/EmergencyProposals"
	Query_Proposal_FullMethodName           = "/initia.gov.v1.Query/Proposal"
	Query_Proposals_FullMethodName          = "/initia.gov.v1.Query/Proposals"
	Query_TallyResult_FullMethodName        = "/initia.gov.v1.Query/TallyResult"
)

// QueryClient is the client API for Query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type QueryClient interface {
	// Params queries params of the gov module.
	Params(ctx context.Context, in *QueryParamsRequest, opts ...grpc.CallOption) (*QueryParamsResponse, error)
	// EmergencyProposals queries emergency proposals.
	EmergencyProposals(ctx context.Context, in *QueryEmergencyProposalsRequest, opts ...grpc.CallOption) (*QueryEmergencyProposalsResponse, error)
	// Proposal queries proposal details based on ProposalID.
	Proposal(ctx context.Context, in *QueryProposalRequest, opts ...grpc.CallOption) (*QueryProposalResponse, error)
	// Proposals queries all proposals based on given status.
	Proposals(ctx context.Context, in *QueryProposalsRequest, opts ...grpc.CallOption) (*QueryProposalsResponse, error)
	// TallyResult queries the tally of a proposal vote.
	TallyResult(ctx context.Context, in *QueryTallyResultRequest, opts ...grpc.CallOption) (*QueryTallyResultResponse, error)
}

type queryClient struct {
	cc grpc.ClientConnInterface
}

func NewQueryClient(cc grpc.ClientConnInterface) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) Params(ctx context.Context, in *QueryParamsRequest, opts ...grpc.CallOption) (*QueryParamsResponse, error) {
	out := new(QueryParamsResponse)
	err := c.cc.Invoke(ctx, Query_Params_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) EmergencyProposals(ctx context.Context, in *QueryEmergencyProposalsRequest, opts ...grpc.CallOption) (*QueryEmergencyProposalsResponse, error) {
	out := new(QueryEmergencyProposalsResponse)
	err := c.cc.Invoke(ctx, Query_EmergencyProposals_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) Proposal(ctx context.Context, in *QueryProposalRequest, opts ...grpc.CallOption) (*QueryProposalResponse, error) {
	out := new(QueryProposalResponse)
	err := c.cc.Invoke(ctx, Query_Proposal_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) Proposals(ctx context.Context, in *QueryProposalsRequest, opts ...grpc.CallOption) (*QueryProposalsResponse, error) {
	out := new(QueryProposalsResponse)
	err := c.cc.Invoke(ctx, Query_Proposals_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) TallyResult(ctx context.Context, in *QueryTallyResultRequest, opts ...grpc.CallOption) (*QueryTallyResultResponse, error) {
	out := new(QueryTallyResultResponse)
	err := c.cc.Invoke(ctx, Query_TallyResult_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
// All implementations must embed UnimplementedQueryServer
// for forward compatibility
type QueryServer interface {
	// Params queries params of the gov module.
	Params(context.Context, *QueryParamsRequest) (*QueryParamsResponse, error)
	// EmergencyProposals queries emergency proposals.
	EmergencyProposals(context.Context, *QueryEmergencyProposalsRequest) (*QueryEmergencyProposalsResponse, error)
	// Proposal queries proposal details based on ProposalID.
	Proposal(context.Context, *QueryProposalRequest) (*QueryProposalResponse, error)
	// Proposals queries all proposals based on given status.
	Proposals(context.Context, *QueryProposalsRequest) (*QueryProposalsResponse, error)
	// TallyResult queries the tally of a proposal vote.
	TallyResult(context.Context, *QueryTallyResultRequest) (*QueryTallyResultResponse, error)
	mustEmbedUnimplementedQueryServer()
}

// UnimplementedQueryServer must be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (UnimplementedQueryServer) Params(context.Context, *QueryParamsRequest) (*QueryParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Params not implemented")
}
func (UnimplementedQueryServer) EmergencyProposals(context.Context, *QueryEmergencyProposalsRequest) (*QueryEmergencyProposalsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EmergencyProposals not implemented")
}
func (UnimplementedQueryServer) Proposal(context.Context, *QueryProposalRequest) (*QueryProposalResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Proposal not implemented")
}
func (UnimplementedQueryServer) Proposals(context.Context, *QueryProposalsRequest) (*QueryProposalsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Proposals not implemented")
}
func (UnimplementedQueryServer) TallyResult(context.Context, *QueryTallyResultRequest) (*QueryTallyResultResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TallyResult not implemented")
}
func (UnimplementedQueryServer) mustEmbedUnimplementedQueryServer() {}

// UnsafeQueryServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to QueryServer will
// result in compilation errors.
type UnsafeQueryServer interface {
	mustEmbedUnimplementedQueryServer()
}

func RegisterQueryServer(s grpc.ServiceRegistrar, srv QueryServer) {
	s.RegisterService(&Query_ServiceDesc, srv)
}

func _Query_Params_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryParamsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Params(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_Params_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Params(ctx, req.(*QueryParamsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_EmergencyProposals_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryEmergencyProposalsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).EmergencyProposals(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_EmergencyProposals_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).EmergencyProposals(ctx, req.(*QueryEmergencyProposalsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Proposal_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryProposalRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Proposal(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_Proposal_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Proposal(ctx, req.(*QueryProposalRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Proposals_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryProposalsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Proposals(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_Proposals_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Proposals(ctx, req.(*QueryProposalsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_TallyResult_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryTallyResultRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).TallyResult(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_TallyResult_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).TallyResult(ctx, req.(*QueryTallyResultRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Query_ServiceDesc is the grpc.ServiceDesc for Query service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Query_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "initia.gov.v1.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Params",
			Handler:    _Query_Params_Handler,
		},
		{
			MethodName: "EmergencyProposals",
			Handler:    _Query_EmergencyProposals_Handler,
		},
		{
			MethodName: "Proposal",
			Handler:    _Query_Proposal_Handler,
		},
		{
			MethodName: "Proposals",
			Handler:    _Query_Proposals_Handler,
		},
		{
			MethodName: "TallyResult",
			Handler:    _Query_TallyResult_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "initia/gov/v1/query.proto",
}
