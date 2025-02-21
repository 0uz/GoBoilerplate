package user

import (
	"reflect"
	"testing"

	"github.com/ouz/goauthboilerplate/internal/domain/user"		
)

func TestNewUserRole(t *testing.T) {
	type args struct {
		userID string
		name   user.UserRoleName
	}
	tests := []struct {
		name    string
		args    args
		want    *user.UserRole
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := user.NewUserRole(tt.args.userID, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewUserRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserRole() = %v, want %v", got, tt.want)
			}
		})
	}
}
