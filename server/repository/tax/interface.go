package tax

import "context"

//go:generate mockery --name=Tax --output=../../tests/mocks --structname=MockTaxRepo

type Tax interface {
	SaveTaxNumber(ctx context.Context, companyId uint64, taxIdNumber, taxIdName string) error
}
