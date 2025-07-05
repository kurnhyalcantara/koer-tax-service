package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/kurnhyalcantara/koer-tax-service/pkg/utils"
)

type Config struct {
	// Listen address is an array of IP addresses and port combinations.
	// Listen address is an array so that this service can listen to many interfaces at once.
	// You can use this value for example: []string{"192.168.1.12:80", "25.49.25.73:80"} to listen to
	// listen to interfaces with IP address of 192.168.1.12 and 25.49.25.73, both on port 80.
	ListenAddress string `config:"LISTEN_ADDRESS"`

	CorsAllowedHeaders []string `config:"CORS_ALLOWED_HEADERS"`
	CorsAllowedMethods []string `config:"CORS_ALLOWED_METHODS"`
	CorsAllowedOrigins []string `config:"CORS_ALLOWED_ORIGINS"`

	Env      string `config:"ENV"` // Environment variable, e.g. "development", "production", etc.
	AppName  string `config:"APP_NAME"`
	AppKey   string `config:"APP_KEY"`
	MaxRetry int    `config:"MAX_RETRY"`

	SwaggerPath string `config:"SWAGGER_PATH"`

	DbDsn string `config:"DB_DSN"`

	DbHost     string `config:"DB_HOST"`
	DbUser     string `config:"DB_USER"`
	DbPassword string `config:"DB_PASSWORD"`
	DbName     string `config:"DB_NAME"`
	DbPort     string `config:"DB_PORT"`
	DbSslmode  string `config:"DB_SSLMODE"`
	DbTimezone string `config:"DB_TIMEZONE"`
	DbMaxRetry string `config:"DB_MAX_RETRY"`
	DbTimeout  string `config:"DB_TIMEOUT"`

	DbMaxOpenConns int `config:"DB_MAX_OPEN_CONNS"`
	DbMaxIdleConns int `config:"DB_MAX_IDLE_CONNS"`

	JwtSecret   string `config:"JWT_SECRET"`
	JwtDuration string `config:"JWT_DURATION"`

	TaskService     string `config:"TASK_SERVICE"`
	AuthService     string `config:"AUTH_SERVICE"`
	WorkflowService string `config:"WORKFLOW_SERVICE"`
	RoleService     string `config:"ROLE_SERVICE"`
	AccountService  string `config:"ACCOUNT_SERVICE"`

	FluentbitHost string `config:"FLUENTBIT_HOST"`
	FluentbitPort string `config:"FLUENTBIT_PORT"`
	LoggerTag     string `config:"LOGGER_TAG"`
	LoggerOutput  string `config:"LOGGER_OUTPUT"`
	LoggerLevel   string `config:"LOGGER_LEVEL"`

	AmqpHost           string `config:"AMQP_HOST"`
	AmqpUser           string `config:"AMQP_USER"`
	AmqpPassword       string `config:"AMQP_PASSWORD"`
	AmqpReconnectDelay string `config:"AMQP_RECONNECT_DELAY"`

	AmqpRegistrationOnlineQueue string `config:"AMQP_REGISTRATION_ONLINE_QUEUE"`

	// #region Redis
	RedisAddress     string `config:"REDIS_ADDRESS"`
	RedisPassword    string `config:"REDIS_PASSWORD"`
	RedisDB          int    `config:"REDIS_DB"`
	RedisExpiredTime int    `config:"REDIS_EXPIRED_TIME"`
	RedisKey         string `config:"REDIS_KEY"`
	// #endregion

	ResCacheTimeout    time.Duration `config:"RESCACHE_TIMEOUT"`
	ResCacheOnProgress time.Duration `config:"RESCACHE_ON_PROGRESS"`
	ResCacheDoneTime   time.Duration `config:"RESCACHE_DONE_TIME"`

	// #region ESB
	ESBTimeout     string `config:"ESB_TIMEOUT"`
	ESBUsername    string `config:"ESB_USERNAME"`
	ESBPassword    string `config:"ESB_PASSWORD"`
	ESBChannelID   string `config:"ESB_CHANNEL_ID"`
	ESBProductCode string `config:"ESB_PRODUCT_CODE"`
	ESBServiceID   string `config:"ESB_SERVICE_ID"`
	ESBTellerID    string `config:"ESB_TELLER_ID"`
	// #endregion
}

func InitConfig() *Config {

	godotenv.Load(".env")

	maxRetry, err := strconv.Atoi(utils.GetEnv("MAX_RETRY", "3"))
	if err != nil {
		panic(err)
	}

	appName := utils.GetEnv("APP_NAME", "")
	if appName == "" {
		appName = utils.GetEnv("ELASTIC_APM_SERVICE_NAME", "")
	}

	redisExpiredInSecond, err := strconv.Atoi(utils.GetEnv("REDIS_EXPIRED_TIME", "30"))
	if err != nil {
		panic(err)
	}

	redisDB, err := strconv.Atoi(utils.GetEnv("REDIS_DB", "0"))
	if err != nil {
		panic(err)
	}

	dbMaxIdleConns, err := strconv.Atoi(utils.GetEnv("DB_MAX_IDLE_CONNS", "0"))
	if err != nil {
		panic(err)
	}

	dbMaxOpenConns, err := strconv.Atoi(utils.GetEnv("DB_MAX_OPEN_CONNS", "0"))
	if err != nil {
		panic(err)
	}

	rescacheTimeout, err := time.ParseDuration(utils.GetEnv("RESCACHE_TIMEOUT", ""))
	if err != nil {
		panic(err)
	}

	rescacheOnProgress, err := time.ParseDuration(utils.GetEnv("RESCACHE_ON_PROGRESS", ""))
	if err != nil {
		panic(err)
	}

	rescacheDone, err := time.ParseDuration(utils.GetEnv("RESCACHE_DONE_TIME", ""))
	if err != nil {
		panic(err)
	}

	return &Config{
		ListenAddress: fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT")),
		CorsAllowedHeaders: []string{
			"Connection", "User-Agent", "Referer",
			"Accept", "Accept-Language", "Content-Type",
			"Content-Language", "Content-Disposition", "Origin",
			"Content-Length", "Authorization", "ResponseType",
			"X-Requested-With", "X-Forwarded-For", "grpc-metadata-process_id",
		},
		CorsAllowedMethods: []string{"GET", "POST", "PATCH", "DELETE", "PUT"},
		CorsAllowedOrigins: []string{"*"},

		Env:      utils.GetEnv("ENV", ""),
		AppName:  appName,
		AppKey:   utils.GetEnv("APP_KEY", ""),
		MaxRetry: maxRetry,

		DbDsn:      utils.GetEnv("DB_DSN", ""),
		DbHost:     utils.GetEnv("DB_HOST", ""),
		DbUser:     utils.GetEnv("DB_USER", ""),
		DbPassword: utils.GetEnv("DB_PASSWORD", ""),
		DbName:     utils.GetEnv("DB_NAME", ""),
		DbPort:     utils.GetEnv("DB_PORT", ""),
		DbSslmode:  utils.GetEnv("DB_SSLMODE", ""),
		DbTimezone: utils.GetEnv("DB_TIMEZONE", ""),
		DbMaxRetry: utils.GetEnv("DB_MAX_RETRY", ""),
		DbTimeout:  utils.GetEnv("DB_TIMEOUT", ""),

		DbMaxOpenConns: dbMaxOpenConns,
		DbMaxIdleConns: dbMaxIdleConns,

		JwtSecret:   utils.GetEnv("JWT_SECRET", "secret"),
		JwtDuration: utils.GetEnv("JWT_DURATION", "48h"),

		SwaggerPath: utils.GetEnv("SWAGGER_PATH", ""),

		FluentbitHost: utils.GetEnv("FLUENTBIT_HOST", ""),
		FluentbitPort: utils.GetEnv("FLUENTBIT_PORT", ""),
		LoggerTag:     utils.GetEnv("LOGGER_TAG", ""),
		LoggerOutput:  utils.GetEnv("LOGGER_OUTPUT", ""),
		LoggerLevel:   utils.GetEnv("LOGGER_LEVEL", ""),

		AmqpHost:           utils.GetEnv("AMQP_HOST", ""),
		AmqpUser:           utils.GetEnv("AMQP_USER", ""),
		AmqpPassword:       utils.GetEnv("AMQP_PASSWORD", ""),
		AmqpReconnectDelay: utils.GetEnv("AMQP_RECONNECT_DELAY", ""),

		AmqpRegistrationOnlineQueue: utils.GetEnv("AMQP_REGISTRATION_ONLINE_QUEUE", "addons-registration-online-queue"),

		// #region Redis
		RedisAddress:     utils.GetEnv("REDIS_ADDRESS", "1234"),
		RedisPassword:    utils.GetEnv("REDIS_PASSWORD", ""),
		RedisDB:          redisDB,
		RedisExpiredTime: redisExpiredInSecond,
		RedisKey:         utils.GetEnv("REDIS_KEY", "Approval Signature"),
		// #endregion

		ResCacheTimeout:    rescacheTimeout,
		ResCacheOnProgress: rescacheOnProgress,
		ResCacheDoneTime:   rescacheDone,

		ESBTimeout:     utils.GetEnv("ESB_TIMEOUT", "120"), // Default to 30 seconds if not set
		ESBUsername:    utils.GetEnv("ESB_USERNAME", "bricams"),
		ESBPassword:    utils.GetEnv("ESB_PASSWORD", "Bricams4dd0ns"),
		ESBChannelID:   utils.GetEnv("ESB_CHANNEL_ID", "CMSX"),
		ESBProductCode: utils.GetEnv("ESB_PRODUCT_CODE", "91200"),
		ESBServiceID:   utils.GetEnv("ESB_SERVICE_ID", "0002O"),
		ESBTellerID:    utils.GetEnv("ESB_TELLER_ID", "0206051"),

		AuthService:    utils.GetEnv("AUTH_SERVICE", ":9105"),
		AccountService: utils.GetEnv("ACCOUNT_SERVICE", ":9093"),
	}
}

func (c *Config) AsString() string {
	data, _ := json.Marshal(c)
	return string(data)
}
