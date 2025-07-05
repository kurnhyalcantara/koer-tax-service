package tax

import (
	"context"
	"fmt"
	"testing"

	"github.com/kurnhyalcantara/koer-tax-service/pkg/log"
	"github.com/kurnhyalcantara/koer-tax-service/server/repository/tax"
	"github.com/kurnhyalcantara/koer-tax-service/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCore_SaveTaxNumber(t *testing.T) {
	type fields struct {
		Logger log.LoggerCore
		Repo   tax.Tax
	}
	type args struct {
		ctx         context.Context
		companyId   uint64
		taxIdNumber string
		taxIdName   string
	}

	type assertErr func(t *testing.T, err error)

	tests := []struct {
		name      string
		fields    fields
		args      args
		wantErr   bool
		assertErr assertErr
	}{
		// TODO: Add test cases.
		{
			name: "companyId is nil",
			args: args{
				companyId: 0,
			},
			wantErr: true,
			assertErr: func(t *testing.T, err error) {
				assert.ErrorContains(t, err, "companyId is required")
			},
		},
		{
			name: "repo return err",
			args: args{
				ctx:         context.TODO(),
				companyId:   123,
				taxIdNumber: "123",
				taxIdName:   "test",
			},
			fields: fields{
				Repo: func() tax.Tax {
					mockTaxRepo := new(mocks.MockTaxRepo)
					mockTaxRepo.On("SaveTaxNumber", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
						Return(fmt.Errorf("Unexpected error"))
					return mockTaxRepo
				}(),
			},
			wantErr: true,
			assertErr: func(t *testing.T, err error) {
				assert.EqualError(t, err, "Unexpected error")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Core{
				Logger: tt.fields.Logger,
				Repo:   tt.fields.Repo,
			}
			err := c.SaveTaxNumber(tt.args.ctx, tt.args.companyId, tt.args.taxIdNumber, tt.args.taxIdName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Core.SaveTaxNumber() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && tt.assertErr != nil {
				tt.assertErr(t, err)
			}
		})
	}
}
