package db

import "time"

type Config struct {
	Host         string
	Port         string
	User         string
	Password     string
	DatabaseName string
	SslMode      string
	TimeZone     string
	MaxRetry     int
	Timeout      time.Duration
}
