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
	for i, violation := range vErr.Violations {
		field := formatField(*violation.Proto.GetField().GetElements()[i].FieldName)
		rule := violation.RuleDescriptor.Name()
		ruleVal := violation.RuleValue

		switch rule {
		case "required":
			return fmt.Errorf("%s is required", field)
		case "email":
			return fmt.Errorf("%s must be a valid email", field)
		case "min_len":
			return fmt.Errorf("%s must be at least %d characters", field, ruleVal.Int())
		case "max_len":
			return fmt.Errorf("%s must be at most %d characters", field, ruleVal.Int())
		default:
			if msg := violation.Proto.GetMessage(); msg != "" {
				return errors.New(err.Error())
			}
			return fmt.Errorf("%s is invalid", field)
		}
	}

	return errors.New(strings.Join(userMessages, ", "))
}

func formatField(field string) string {
	// Ubah snake_case ke Title Case
	parts := strings.Split(field, "_")
	for i := range parts {
		parts[i] = strings.Title(parts[i])
	}
	return strings.Join(parts, " ")
}
