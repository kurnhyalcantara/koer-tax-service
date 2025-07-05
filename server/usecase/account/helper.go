package account

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ValidateAccountStatus(accountStatus, accountType, productCode string) error {
	accountStatusDesc := GetAccountStatusDesc(accountStatus)
	isSA := accountType == "SA"
	isCA := accountType == "CA"

	// Validate Savings Account (SA)
	if isSA && accountStatus != "1" {
		if IsInvalidAccountStatus(accountStatus) {
			return status.Errorf(codes.Internal, "Invalid Debit Account. Account is %s", accountStatusDesc)
		}
	}

	// Validate Current Account (CA) with productCode "G3"
	if isCA && productCode == "G3" && accountStatus != "6" {
		if IsInvalidAccountStatus(accountStatus) {
			return status.Errorf(codes.Internal, "Invalid Sender Account. Account is %s", accountStatusDesc)
		}
	}

	// Validate Current Account (CA) with other productCodes
	if isCA && accountStatus != "1" && productCode != "G3" {
		if IsInvalidAccountStatus(accountStatus) {
			return status.Errorf(codes.Internal, "Invalid Sender Account. Account is %s", accountStatusDesc)
		}
	}

	return nil
}

// Helper function to map accountStatus to a human-readable description
func GetAccountStatusDesc(accountStatus string) string {
	switch accountStatus {
	case "2":
		return "closed"
	case "3":
		return "matured"
	case "4":
		return "new today"
	case "5":
		return "zero actual"
	case "6", "7":
		return "frozen"
	case "8":
		return "Charged Off"
	case "9":
		return "dormant"
	default:
		return "unknown"
	}
}

// Helper function to check if the accountStatus is invalid
func IsInvalidAccountStatus(accountStatus string) bool {
	invalidStatuses := []string{"2", "3", "4", "5", "6", "7", "8", "9"}
	for _, status := range invalidStatuses {
		if accountStatus == status {
			return true
		}
	}
	return false
}
