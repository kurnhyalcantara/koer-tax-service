package tax

import (
	"context"
	"fmt"

	"github.com/kurnhyalcantara/koer-tax-service/pkg/log"
	"github.com/kurnhyalcantara/koer-tax-service/server/infrastructure/db"
)

type CoreDB struct {
	db     db.DbCore
	logger log.LoggerCore
}

// SaveTaxNumber implements Tax.
func (c *CoreDB) SaveTaxNumber(ctx context.Context, companyId uint64, taxIdNumber string, taxIdName string) error {

	c.logger.WithFunctionName("SaveTaxNumber")

	var query string = `
		INSERT INTO tax_saved_payer (company_id, tax_id_number, tax_id_name)
		VALUES ($1, $2, $3)
		ON CONFLICT (company_id, tax_id_number) DO UPDATE
			SET tax_id_name = EXCLUDED.tax_id_name,
    updated_at = now()
	`

	ctxTimeout, ctxCancel := context.WithTimeout(ctx, c.db.GetTimeout())
	defer ctxCancel()

	_, err := c.db.GetDb().ExecContext(ctxTimeout, query, companyId, taxIdNumber, taxIdName)
	if err != nil {
		c.logger.Error(log.LogPayload{Message: "Error exec query: ", Metadata: err})
		return fmt.Errorf("Unexpected error")
	}

	return nil

}

func NewTaxRepository(db db.DbCore, logger log.LoggerCore) Tax {
	return &CoreDB{
		db:     db,
		logger: logger,
	}
}
