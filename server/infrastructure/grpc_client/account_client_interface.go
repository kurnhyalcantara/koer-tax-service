package grpcclient

import (
	"context"

	account_pb "github.com/kurnhyalcantara/koer-tax-service/protogen/account-service"
	"google.golang.org/grpc"
)

//go:generate mockery --name=AccountService --output=../../tests/mocks --structname=MockAccountService

type AccountService interface {
	// Add methods specific to account service here
	ListAccountRPC(ctx context.Context, req *account_pb.ListAccountRequest, opts ...grpc.CallOption) (*account_pb.ListAccountResponse, error)
	ValidateAccountRPC(ctx context.Context, req *account_pb.ValidateAccountRequest, opts ...grpc.CallOption) (*account_pb.ValidateAccountResponse, error)
	ValidateAccount(ctx context.Context, req *account_pb.ValidateAccountRequest, opts ...grpc.CallOption) (*account_pb.ValidateAccountResponse, error)
}
