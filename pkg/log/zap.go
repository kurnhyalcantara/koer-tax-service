package log

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fluent/fluent-logger-golang/fluent"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//go:generate mockery --name=LoggerCore --output=../../tests/mocks --structname=MockLoggerCore

type LoggerCore interface {
	WithFunctionName(string) LoggerCore
	WithProcessId(string) LoggerCore
	WithStartTime(*time.Time) LoggerCore
	WithTaskId(uint64) LoggerCore
	EndTime() LoggerCore

	Info(payload LogPayload)
	Warn(payload LogPayload)
	Error(payload LogPayload)
	Debug(payload LogPayload)
	Fatal(payload LogPayload)
	QueueMessageInfo(payload LogPayload)
	QueueMessageError(payload LogPayload)
}

type LoggerConfig struct {
	Env           string
	ProductName   string
	ServiceName   string
	LogLevel      string
	LogOutput     string
	FluentbitHost string
	FluentbitPort string
}

type Logger struct {
	zapLog       *zap.Logger
	productName  string
	serviceName  string
	hostName     string
	functionName string
	processId    string
	startTime    *time.Time
	endTime      *time.Time
	taskId       uint64
}

func (l *Logger) WithFunctionName(fn string) LoggerCore {
	// clone the logger (shallow copy)
	newLogger := *l
	newLogger.functionName = fn
	return &newLogger
}

// WithProcessId implements LoggerCore.
func (l *Logger) WithProcessId(pid string) LoggerCore {
	newLogger := *l           // shallow copy struct Logger
	newLogger.processId = pid // set processId baru
	return &newLogger         // return pointer ke copy
}

// WithStartTime implements LoggerCore.
func (l *Logger) WithStartTime(startTime *time.Time) LoggerCore {
	newLogger := *l
	newLogger.startTime = startTime
	return &newLogger
}

// WithTaskId implements LoggerCore.
func (l *Logger) WithTaskId(taskId uint64) LoggerCore {
	newLogger := *l
	newLogger.taskId = taskId
	return &newLogger
}

func New(loggerConfig *LoggerConfig) LoggerCore {
	logLevel, levelErr := zap.ParseAtomicLevel(loggerConfig.LogLevel)
	if levelErr != nil {
		panic(levelErr)
	}

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.MessageKey = "message"

	zapEncoder := zapcore.NewJSONEncoder(encoderCfg)

	hostname, _ := os.Hostname()

	var core zapcore.Core

	if strings.EqualFold(loggerConfig.LogOutput, "elastic") {

		fluentPort, err := strconv.Atoi(loggerConfig.FluentbitPort)
		if err != nil {
			fluentPort = 24421
		}

		fluentBithook := NewFluentBitHook(fluent.Config{
			FluentHost: loggerConfig.FluentbitHost,
			FluentPort: fluentPort,
		})

		if logLevel.Level() == zap.DebugLevel {
			core = zapcore.NewTee(
				zapcore.NewCore(zapEncoder, fluentBithook, logLevel),
				zapcore.NewCore(zapEncoder, os.Stdout, logLevel),
			)
		} else {
			core = zapcore.NewTee(
				zapcore.NewCore(zapEncoder, fluentBithook, logLevel),
			)
		}
	} else {
		core = zapcore.NewTee(
			zapcore.NewCore(zapEncoder, os.Stdout, logLevel),
		)
	}

	zapLogger := zap.New(core)

	return &Logger{
		zapLogger,
		loggerConfig.ProductName,
		loggerConfig.ServiceName,
		hostname,
		"",
		"",
		nil,
		nil,
		0,
	}
}
