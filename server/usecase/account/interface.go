package account

import "context"

//go:generate mockery --name=UseCase --output=../../tests/mocks --structname=MockAccountUseCase

type UseCase interface {
	V3ValidateAccount(ctx context.Context, accountNumber string) error
}
