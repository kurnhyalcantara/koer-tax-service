package validator

import (
	"buf.build/go/protovalidate"
	"google.golang.org/protobuf/proto"
)

type ProtoValidator struct {
	validator protovalidate.Validator
}

func NewProtoValidator() (*ProtoValidator, error) {
	v, err := protovalidate.New()
	if err != nil {
		return nil, err
	}
	return &ProtoValidator{validator: v}, nil
}

func (p *ProtoValidator) Validate(msg proto.Message) error {
	err := p.validator.Validate(msg)
	if err != nil {
		return WrapValidationError(err)
	}
	return nil
}
