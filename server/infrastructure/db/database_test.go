package db

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	databasemock "github.com/kurnhyalcantara/koer-tax-service/server/infrastructure/db/mock"

	"github.com/stretchr/testify/suite"
)

type DatabaseTestSuite struct {
	suite.Suite
	ctx context.Context
}

func (s *DatabaseTestSuite) SetupTest() {
	s.ctx = context.Background()
}

func TestInitDatabase(t *testing.T) {
	suite.Run(t, new(DatabaseTestSuite))
}

func (s *DatabaseTestSuite) TestDatabase_InitConnectionDB() {
	type expectation struct {
		out *DbSql
		err error
	}

	tests := map[string]struct {
		driverName string
		config     Config
		expected   expectation
	}{
		"Success": {
			driverName: "mysql",
			config:     Config{},
			expected: expectation{
				out: &DbSql{
					driverName: "mysql",
				},
				err: nil,
			},
		},
	}

	for scenario, tt := range tests {
		s.T().Run(scenario, func(t *testing.T) {
			dbSql := InitConnectionDB(tt.driverName, tt.config, nil)

			if tt.expected.out.driverName != dbSql.GetDriverName() {
				t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.out, dbSql)
			}
		})
	}
}

func (s *DatabaseTestSuite) TestDatabase_GetDbSql() {
	type expectation struct {
		sqlDb *sql.DB
	}

	tests := map[string]struct {
		dbSql    *DbSql
		expected expectation
	}{
		"Success": {
			dbSql: &DbSql{
				SqlDb: &sql.DB{},
			},
			expected: expectation{
				sqlDb: &sql.DB{},
			},
		},
		"Failed": {
			dbSql: &DbSql{},
			expected: expectation{
				sqlDb: nil,
			},
		},
	}

	for scenario, tt := range tests {
		s.T().Run(scenario, func(t *testing.T) {
			dbSql := tt.dbSql.SqlDb

			if dbSql != nil {
				if tt.expected.sqlDb == nil {
					t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.sqlDb, dbSql)
				}
			} else {
				if tt.expected.sqlDb != nil {
					t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.sqlDb, dbSql)
				}
			}
		})
	}
}

func (s *DatabaseTestSuite) TestDatabase_GetDbTx() {
	type expectation struct {
		sqlTx *sql.Tx
	}

	tests := map[string]struct {
		dbSql    *DbSql
		expected expectation
	}{
		"Success": {
			dbSql: &DbSql{
				SqlTx: &sql.Tx{},
			},
			expected: expectation{
				sqlTx: &sql.Tx{},
			},
		},
		"Failed": {
			dbSql: &DbSql{},
			expected: expectation{
				sqlTx: nil,
			},
		},
	}

	for scenario, tt := range tests {
		s.T().Run(scenario, func(t *testing.T) {
			dbSql := tt.dbSql.SqlTx

			if dbSql != nil {
				if tt.expected.sqlTx == nil {
					t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.sqlTx, dbSql)
				}
			} else {
				if tt.expected.sqlTx != nil {
					t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.sqlTx, dbSql)
				}
			}
		})
	}
}

func (s *DatabaseTestSuite) TestDatabase_GetTimeout() {
	type expectation struct {
		sqlTimeout time.Duration
	}

	tests := map[string]struct {
		dbSql    *DbSql
		expected expectation
	}{
		"Success": {
			dbSql: &DbSql{
				Config: Config{Timeout: 30},
			},
			expected: expectation{
				sqlTimeout: 30,
			},
		},
	}

	for scenario, tt := range tests {
		s.T().Run(scenario, func(t *testing.T) {
			timeout := tt.dbSql.GetTimeout()

			if timeout != tt.expected.sqlTimeout {
				t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.sqlTimeout, timeout)
			}
		})
	}
}

func (s *DatabaseTestSuite) TestDatabase_Connect() {
	type expectation struct {
		err error
	}

	tests := map[string]struct {
		dbSql    *DbSql
		config   Config
		expected expectation
	}{
		"Success": {
			dbSql: &DbSql{
				driverName: "mysql",
				Config:     Config{},
				Dbw:        &databasemock.DatabaseMock{},
			},
			expected: expectation{
				err: nil,
			},
		},
		"Failed": {
			dbSql: &DbSql{
				driverName: "",
				Config:     Config{},
				Dbw:        &databasemock.DatabaseMock{},
			},
			expected: expectation{
				err: errors.New("failed connect to database postgresql"),
			},
		},
	}

	for scenario, tt := range tests {
		s.T().Run(scenario, func(t *testing.T) {
			err := tt.dbSql.Connect()

			if tt.expected.err != nil {
				if tt.expected.err.Error() != err.Error() {
					t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.err, err)
				}
			}
		})
	}
}

func (s *DatabaseTestSuite) TestDatabase_GetConnection() {
	type expectation struct {
		err error
	}

	tests := map[string]struct {
		dbSql    *DbSql
		expected expectation
	}{
		"Success": {
			dbSql: &DbSql{
				driverName: "mysql",
				Conn: &databasemock.DatabaseMock{
					DbPq: &sql.DB{},
				},
			},
			expected: expectation{
				err: nil,
			},
		},
		"Failed": {
			dbSql: &DbSql{
				driverName: "mysql",
				Dbw:        &databasemock.DatabaseMock{},
				Conn:       &databasemock.DatabaseMock{},
			},
			expected: expectation{
				err: nil,
			},
		},
	}

	for scenario, tt := range tests {
		s.T().Run(scenario, func(t *testing.T) {
			tt.dbSql.AddCounter()
			err := tt.dbSql.CheckConnection()

			if tt.expected.err != err {
				t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.err, err)
			}
		})
	}
}

func (s *DatabaseTestSuite) TestDatabase_TryConnect() {
	type expectation struct {
		err error
	}

	tests := map[string]struct {
		dbSql    *DbSql
		expected expectation
	}{
		"Success": {
			dbSql: &DbSql{
				driverName: "mysql",
				Config: Config{
					MaxRetry: 1,
				},
				Dbw: &databasemock.DatabaseMock{},
				Conn: &databasemock.DatabaseMock{
					DbPq: &sql.DB{},
				},
			},
			expected: expectation{
				err: nil,
			},
		},
		"Failed": {
			dbSql: &DbSql{
				driverName: "",
				Config: Config{
					MaxRetry: 1,
				},
				Dbw: &databasemock.DatabaseMock{},
				Conn: &databasemock.DatabaseMock{
					DbPq: &sql.DB{},
				},
			},
			expected: expectation{
				err: errors.New("failed connect to database postgresql"),
			},
		},
	}

	for scenario, tt := range tests {
		s.T().Run(scenario, func(t *testing.T) {
			err := tt.dbSql.TryConnect()

			if err != nil {
				if tt.expected.err.Error() != err.Error() {
					t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.err, err)
				}
			}
		})
	}
}

func (s *DatabaseTestSuite) TestDatabase_ClosePmConnection() {
	type expectation struct {
		err error
	}

	tests := map[string]struct {
		dbSql    *DbSql
		expected expectation
	}{
		"Success": {
			dbSql: &DbSql{
				Conn: &databasemock.DatabaseMock{},
			},
			expected: expectation{
				err: nil,
			},
		},
		"Failed": {
			dbSql: &DbSql{},
			expected: expectation{
				err: errors.New("database connection already closed"),
			},
		},
	}

	for scenario, tt := range tests {
		s.T().Run(scenario, func(t *testing.T) {
			err := tt.dbSql.ClosePmConnection()

			if err != nil {
				if tt.expected.err.Error() != err.Error() {
					t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.err, err)
				}
			}
		})
	}
}

func (s *DatabaseTestSuite) TestDatabase_SetMaxConns() {
	type expectation struct {
		err error
	}

	tests := map[string]struct {
		dbSql    *DbSql
		expected expectation
	}{
		"Success": {
			dbSql: &DbSql{
				Conn: &databasemock.DatabaseMock{},
			},
			expected: expectation{
				err: nil,
			},
		},
	}

	for scenario, tt := range tests {
		s.T().Run(scenario, func(t *testing.T) {
			tt.dbSql.SetMaxIdleConns(0)
			tt.dbSql.SetMaxOpenConns(100)
		})
	}
}

func (s *DatabaseTestSuite) TestDatabase_Begin() {
	type expectation struct {
		err error
	}

	tests := map[string]struct {
		dbSql    *DbSql
		expected expectation
	}{
		"Success": {
			dbSql: &DbSql{
				Conn: &databasemock.DatabaseMock{},
			},
			expected: expectation{
				err: nil,
			},
		},
	}

	for scenario, tt := range tests {
		s.T().Run(scenario, func(t *testing.T) {
			err := tt.dbSql.StartTransaction()

			if err != nil {
				if tt.expected.err.Error() != err.Error() {
					t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.err, err)
				}
			}
		})
	}
}
