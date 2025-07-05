package tax

import "testing"

func TestValidateAccountStatus(t *testing.T) {
	type args struct {
		accountStatus string
		accountType   string
		productCode   string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateAccountStatus(tt.args.accountStatus, tt.args.accountType, tt.args.productCode); (err != nil) != tt.wantErr {
				t.Errorf("ValidateAccountStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
