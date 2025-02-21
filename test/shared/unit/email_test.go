package shared

import (
	"reflect"
	"testing"

	"github.com/ouz/goauthboilerplate/internal/domain/shared"
)

func TestNewEmail(t *testing.T) {
	type args struct {
		address string
	}
	tests := []struct {
		name    string
		args    args
		want    shared.Email
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := shared.NewEmail(tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}
