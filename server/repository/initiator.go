package repository

import (
	"github.com/kurnhyalcantara/koer-tax-service/pkg/log"
	"github.com/kurnhyalcantara/koer-tax-service/server/infrastructure/db"
	"github.com/kurnhyalcantara/koer-tax-service/server/repository/tax"
)

type Repositories struct {
	TaxRepo tax.Tax
}

type Dependencies struct {
	db     db.DbCore
	logger log.LoggerCore
}

func InitRepositories(db db.DbCore, logger log.LoggerCore) *Repositories {
	return &Repositories{
		TaxRepo: tax.NewTaxRepository(db, logger),
	}
}
