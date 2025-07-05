package tax

import (
	"context"
	"errors"

	"github.com/kurnhyalcantara/koer-tax-service/pkg/log"
	"github.com/kurnhyalcantara/koer-tax-service/server/repository/tax"
)

type Core struct {
	Logger log.LoggerCore
	Repo   tax.Tax
}

// SaveTaxNumber implements UseCase.
func (c *Core) SaveTaxNumber(ctx context.Context, companyId uint64, taxIdNumber, taxIdName string) error {
	if companyId == 0 {
		return errors.New("companyId is required")
	}
	return c.Repo.SaveTaxNumber(ctx, companyId, taxIdNumber, taxIdName)
}

func NewTaxUseCase(logger log.LoggerCore, repo tax.Tax) UseCase {
	return &Core{
		Logger: logger,
		Repo:   repo,
	}
}
