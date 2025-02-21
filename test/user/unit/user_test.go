package user

import (
	"reflect"
	"testing"

	"github.com/ouz/goauthboilerplate/internal/domain/user"
	vo "github.com/ouz/goauthboilerplate/internal/domain/shared"
)

func TestNewUser(t *testing.T) {
	type args struct {
		username string
		password string
		email    vo.Email
	}
	tests := []struct {
		name    string
		args    args
		want    *user.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := user.NewUser(tt.args.username, tt.args.password, tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewAnonymousUser(t *testing.T) {
	tests := []struct {
		name    string
		want    *user.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := user.NewAnonymousUser()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAnonymousUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAnonymousUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUser_AddCredential(t *testing.T) {
	type args struct {
		password string
	}
	tests := []struct {
		name    string
		u       *user.User
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.u.AddCredential(tt.args.password); (err != nil) != tt.wantErr {
				t.Errorf("User.AddCredential() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUser_AddConfirmation(t *testing.T) {
	tests := []struct {
		name    string
		u       *user.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.u.AddConfirmation(); (err != nil) != tt.wantErr {
				t.Errorf("User.AddConfirmation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUser_IsPasswordValid(t *testing.T) {
	type args struct {
		password string
	}
	tests := []struct {
		name string
		u    *user.User
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.u.IsPasswordValid(tt.args.password); got != tt.want {
				t.Errorf("User.IsPasswordValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUser_HasRole(t *testing.T) {
	type args struct {
		role user.UserRoleName
	}
	tests := []struct {
		name string
		u    *user.User
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.u.HasRole(tt.args.role); got != tt.want {
				t.Errorf("User.HasRole() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUser_Confirm(t *testing.T) {
	tests := []struct {
		name string
		u    *user.User
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.u.Confirm()
		})
	}
}
