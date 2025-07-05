package grpc

import (
	"context"
	"log"
	"reflect"
	"testing"

	"github.com/kurnhyalcantara/koer-tax-service/pkg/validator"
	pb "github.com/kurnhyalcantara/koer-tax-service/protogen/koer-tax-service"
	jwtmanagerDom "github.com/kurnhyalcantara/koer-tax-service/server/domain/security/jwt_manager"
	jwtmanager "github.com/kurnhyalcantara/koer-tax-service/server/infrastructure/security/jwt_manager"
	"github.com/kurnhyalcantara/koer-tax-service/server/usecase"
	"github.com/kurnhyalcantara/koer-tax-service/tests/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestHandler_V3SaveTaxNumber(t *testing.T) {
	type fields struct {
		validator            *validator.ProtoValidator
		uc                   *usecase.UseCases
		manager              jwtmanager.JwtManagerCore
		KoerTaxServiceServer pb.KoerTaxServiceServer
	}
	type args struct {
		ctx context.Context
		req *pb.V3SaveTaxNumberRequest
	}

	type assertErr func(t *testing.T, err error)

	tests := []struct {
		name      string
		fields    fields
		args      args
		want      *pb.GeneralBodyResponse
		wantErr   bool
		assertErr assertErr
	}{
		// TODO: Add test cases.
		{
			name: "failed get user data",
			fields: fields{
				manager: func() jwtmanager.JwtManagerCore {
					mockJwt := new(mocks.MockJwtManagerCore)
					mockJwt.On("GetMeFromMD", mock.Anything).
						Return(nil, nil, status.Error(codes.Unauthenticated, "Unauthorized"))
					return mockJwt
				}(),
			},
			args: args{
				ctx: context.TODO(),
				req: &pb.V3SaveTaxNumberRequest{},
			},
			want:    nil,
			wantErr: true,
			assertErr: func(t *testing.T, err error) {
				require.ErrorContains(t, err, "Unauthorized")
			},
		},
		{
			name: "unauthorized user type",
			fields: fields{
				manager: func() jwtmanager.JwtManagerCore {
					mockJwt := new(mocks.MockJwtManagerCore)
					mockJwt.On("GetMeFromMD", mock.Anything).
						Return(&jwtmanagerDom.UserData{
							UserType: "ba",
						}, nil, nil)
					return mockJwt
				}(),
			},
			args: args{
				ctx: context.TODO(),
				req: &pb.V3SaveTaxNumberRequest{},
			},
			want:    nil,
			wantErr: true,
			assertErr: func(t *testing.T, err error) {
				require.ErrorContains(t, err, "Permission Denied")
			},
		},
		{
			name: "request is nil",
			fields: fields{
				manager: func() jwtmanager.JwtManagerCore {
					mockJwt := new(mocks.MockJwtManagerCore)
					mockJwt.On("GetMeFromMD", mock.Anything).
						Return(&jwtmanagerDom.UserData{
							UserType: "cu",
						}, nil, nil)
					return mockJwt
				}(),
				validator: func() *validator.ProtoValidator {
					protoValidator, err := validator.NewProtoValidator()
					if err != nil {
						log.Fatal(err)
					}
					return protoValidator
				}(),
			},
			args: args{
				ctx: context.TODO(),
				req: &pb.V3SaveTaxNumberRequest{},
			},
			want:    nil,
			wantErr: true,
			assertErr: func(t *testing.T, err error) {
				require.Equal(t, err, "taxIdnumber is required")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				validator:            tt.fields.validator,
				uc:                   tt.fields.uc,
				manager:              tt.fields.manager,
				KoerTaxServiceServer: tt.fields.KoerTaxServiceServer,
			}
			got, err := h.V3SaveTaxNumber(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Handler.V3SaveTaxNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.assertErr != nil {
				tt.assertErr(t, err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Handler.V3SaveTaxNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}
