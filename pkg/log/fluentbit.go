package log

import (
	"log"

	"github.com/fluent/fluent-logger-golang/fluent"
)

type FluentBitCore interface {
	Write(p []byte) (n int, err error)
	Sync() error
}

type FluentBitHook struct {
	fluentbit *fluent.Fluent
}

func NewFluentBitHook(flntCfg fluent.Config) FluentBitCore {
	fluentBit, err := fluent.New(flntCfg)
	if err != nil {
		log.Fatalf("failed init fluent bit: %v", err)
	}

	return &FluentBitHook{
		fluentbit: fluentBit,
	}
}

// Sync implements FluentBitCore.
func (f *FluentBitHook) Sync() error {
	return f.fluentbit.Close()
}

// Write implements FluentBitCore.
func (f *FluentBitHook) Write(p []byte) (n int, err error) {
	logEntry := map[string]string{"message": string(p)}
	err = f.fluentbit.Post("zap", logEntry)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}
