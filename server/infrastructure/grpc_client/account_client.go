package grpcclient

import (
	"context"

	account_pb "github.com/kurnhyalcantara/koer-tax-service/protogen/account-service"
	"google.golang.org/grpc"
)

type AccountGrpcClient struct {
	client AccountService
}

// ListAccountRPC implements AccountService.
func (a *AccountGrpcClient) ListAccountRPC(ctx context.Context, req *account_pb.ListAccountRequest, opts ...grpc.CallOption) (*account_pb.ListAccountResponse, error) {
	return a.client.ListAccountRPC(ctx, req, opts...)
}

// ValidateAccount implements AccountService.
func (a *AccountGrpcClient) ValidateAccount(ctx context.Context, req *account_pb.ValidateAccountRequest, opts ...grpc.CallOption) (*account_pb.ValidateAccountResponse, error) {
	return a.client.ValidateAccount(ctx, req, opts...)
}

// ValidateAccountRPC implements AccountService.
func (a *AccountGrpcClient) ValidateAccountRPC(ctx context.Context, req *account_pb.ValidateAccountRequest, opts ...grpc.CallOption) (*account_pb.ValidateAccountResponse, error) {
	return a.client.ValidateAccountRPC(ctx, req, opts...)
}

func NewAccountGrpcClient(conn *grpc.ClientConn) AccountService {
	return &AccountGrpcClient{
		client: account_pb.NewApiServiceClient(conn),
	}
}
