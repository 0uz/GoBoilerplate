package user

import (
	"reflect"
	"testing"

	"github.com/ouz/goauthboilerplate/internal/domain/user"
)

func TestNewPassword(t *testing.T) {
	type args struct {
		plaintext string
	}
	tests := []struct {
		name    string
		args    args
		want    *user.Password
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := user.NewPassword(tt.args.plaintext)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPassword_Verify(t *testing.T) {
	type args struct {
		plaintext string
	}
	tests := []struct {
		name string
		p    *user.Password
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.Verify(tt.args.plaintext); got != tt.want {
				t.Errorf("Password.Verify() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPassword_Hashed(t *testing.T) {
	tests := []struct {
		name string
		p    *user.Password
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.Hashed(); got != tt.want {
				t.Errorf("Password.Hashed() = %v, want %v", got, tt.want)
			}
		})
	}
}
