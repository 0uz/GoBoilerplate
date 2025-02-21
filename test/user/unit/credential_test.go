package user

import (
	"reflect"
	"testing"

	"github.com/ouz/goauthboilerplate/internal/domain/user"
)

func TestNewCredential(t *testing.T) {
	type args struct {
		credentialType user.CredentialType
		secret         string
	}
	tests := []struct {
		name    string
		args    args
		want    *user.Credential
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := user.NewCredential(tt.args.credentialType, tt.args.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCredential() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCredential() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCredential_IsPasswordValid(t *testing.T) {
	type args struct {
		password string
	}
	tests := []struct {
		name string
		c    *user.Credential
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.IsPasswordValid(tt.args.password); got != tt.want {
				t.Errorf("Credential.IsPasswordValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
