package validator

import "google.golang.org/protobuf/proto"

//go:generate mockery --name=Validator --output=../../tests/mocks --structname=MockValidator

type Validator interface {
	Validate(msg proto.Message) error
}
