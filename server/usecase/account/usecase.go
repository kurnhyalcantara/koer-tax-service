package account

import (
	"context"

	"github.com/kurnhyalcantara/koer-tax-service/pkg/constants"
	"github.com/kurnhyalcantara/koer-tax-service/pkg/log"
	account_pb "github.com/kurnhyalcantara/koer-tax-service/protogen/account-service"
	grpcclient "github.com/kurnhyalcantara/koer-tax-service/server/infrastructure/grpc_client"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Core struct {
	AccountService grpcclient.AccountService
	Logger         log.LoggerCore
}

// V3ValidateAccount implements UseCase.
func (c *Core) V3ValidateAccount(ctx context.Context, accountNumber string) error {
	c.Logger.
		WithFunctionName("V3ValidateAccount")

	accountData, err := c.AccountService.ListAccountRPC(ctx, &account_pb.ListAccountRequest{
		Account: &account_pb.Account{
			AccountNumber: accountNumber,
		},
		Limit:     1,
		Page:      1,
		ProductID: constants.ProductID,
	})
	if err != nil {
		c.Logger.Error(log.LogPayload{Message: "Failed list account rpc", Metadata: err})
		return status.Error(codes.Internal, "Unexpected Error")
	}

	if len(accountData.Data) <= 0 {
		return status.Error(codes.InvalidArgument, "Account not found")
	}

	if accountData.Data[0].AccountType == constants.AccountTypeVA {
		_, err = c.AccountService.ValidateAccount(ctx, &account_pb.ValidateAccountRequest{
			AccountNo: accountNumber,
			Type:      "VA",
		})
		if err != nil {
			c.Logger.Error(log.LogPayload{Message: "Failed validate account (va)", Metadata: err})
			return status.Error(codes.Internal, "Unexpected Error")
		}
	} else {
		validateRes, err := c.AccountService.ValidateAccountRPC(ctx, &account_pb.ValidateAccountRequest{
			AccountNo: accountNumber,
		})
		if err != nil {
			c.Logger.Error(log.LogPayload{Message: "Failed validate account", Metadata: err})
			return status.Error(codes.Internal, "Unexpected Error")
		}

		if err := ValidateAccountStatus(
			validateRes.Data.Status,
			validateRes.Data.ProductCode,
			validateRes.Data.AcctType); err != nil {
			return err
		}
	}

	return nil
}

func NewAccountUseCase(client grpcclient.AccountService, logger log.LoggerCore) UseCase {
	return &Core{
		AccountService: client,
		Logger:         logger,
	}
}
