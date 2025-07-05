package databasemock

import (
	"database/sql"
	"errors"

	databasewrapper "github.com/kurnhyalcantara/koer-tax-service/server/infrastructure/db/wrapper"
)

type DatabaseMock struct {
	DbPq *sql.DB
	databasewrapper.DatabaseInterface
	databasewrapper.DatabaseConnectionInterface
}

func (dm *DatabaseMock) Open(driverName, dataSourceName string) (*sql.DB, error) {
	if driverName == "" {
		return nil, errors.New("failed connect to database postgresql")
	}

	return &sql.DB{}, nil
}

func (dm *DatabaseMock) Ping() error {
	if dm.DbPq == nil {
		return errors.New("ping failed, connection timeout or disconnect")
	}
	return nil
}

func (dm *DatabaseMock) Close() error {
	return nil
}

func (dm *DatabaseMock) SetMaxIdleConns(n int) {}
func (dm *DatabaseMock) SetMaxOpenConns(n int) {}

func (dm *DatabaseMock) Begin() (*sql.Tx, error) {
	return &sql.Tx{}, nil
}
