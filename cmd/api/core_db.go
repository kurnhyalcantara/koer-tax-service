package main

import (
	"log"
	"strconv"
	"time"

	"github.com/kurnhyalcantara/koer-tax-service/server/infrastructure/db"
	databasewrapper "github.com/kurnhyalcantara/koer-tax-service/server/infrastructure/db/wrapper"
)

func StartDBConnection() db.DbCore {

	var err error
	maxRetry, convErr := strconv.Atoi(appConfig.DbMaxRetry)
	if convErr != nil {
		maxRetry = 3
	}

	dbTimeout, convErr := strconv.Atoi(appConfig.DbTimeout)
	if convErr != nil {
		dbTimeout = 120
	}

	dbMain := db.InitConnectionDB("postgres", db.Config{
		Host:         appConfig.DbHost,
		Port:         appConfig.DbPort,
		User:         appConfig.DbUser,
		Password:     appConfig.DbPassword,
		DatabaseName: appConfig.DbName,
		SslMode:      appConfig.DbSslmode,
		TimeZone:     appConfig.DbTimezone,
		MaxRetry:     maxRetry,
		Timeout:      time.Duration(dbTimeout) * time.Second,
	}, &databasewrapper.DatabaseWrapper{})

	err = dbMain.Connect()

	if err != nil {
		log.Fatalf("failed init db connect: %v", err)
	}

	dbMain.SetMaxIdleConns(appConfig.DbMaxIdleConns)
	dbMain.SetMaxOpenConns(appConfig.DbMaxOpenConns)

	return dbMain
}

func CloseDBConnection(dbMain db.DbCore) error {
	if err := dbMain.ClosePmConnection(); err != nil {
		log.Fatalf("failed close db connection: %v", err)
	}
	return nil
}
