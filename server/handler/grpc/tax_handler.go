package grpc

import (
	"context"
	"net/http"

	pb "github.com/kurnhyalcantara/koer-tax-service/protogen/koer-tax-service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) V3SaveTaxNumber(ctx context.Context, req *pb.V3SaveTaxNumberRequest) (*pb.GeneralBodyResponse, error) {

	currentUser, _, err := h.manager.GetMeFromMD(ctx)
	if err != nil {
		return nil, err
	}

	if currentUser.UserType != "cu" {
		return nil, status.Error(codes.PermissionDenied, "Permission Denied")
	}

	if err := h.validator.Validate(req); err != nil {
		return nil, err
	}

	if err := h.uc.Tax.SaveTaxNumber(ctx, currentUser.CompanyID, req.TaxIdNumber, req.TaxOwnerName); err != nil {
		return nil, err
	}

	return &pb.GeneralBodyResponse{
		Error:   false,
		Code:    http.StatusOK,
		Message: "Success",
	}, nil

}
