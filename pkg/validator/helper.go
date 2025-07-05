package validator

import (
	"errors"
	"fmt"
	"strings"

	"buf.build/go/protovalidate"
)

// WrapValidationError converts a ValidationError into a human-readable error message.
func WrapValidationError(err error) error {
	var vErr *protovalidate.ValidationError
	if !errors.As(err, &vErr) {
		return err
	}

	var userMessages []string
	for _, violation := range vErr.Violations {
		msg := violation.Proto.Message
		userMessages = append(userMessages, *msg)
	}

	return fmt.Errorf(strings.Join(userMessages, ", "))
}
