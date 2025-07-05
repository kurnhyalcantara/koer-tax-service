package tax

import (
	"context"
)

//go:generate mockery --name=UseCase --output=../../tests/mocks --structname=MockTaxUseCase

type UseCase interface {
	SaveTaxNumber(ctx context.Context, companyId uint64, taxIdNumber, taxIdName string) error
}
