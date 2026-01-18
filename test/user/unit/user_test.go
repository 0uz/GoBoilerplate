package user

import (
	"testing"

	"github.com/google/uuid"
	vo "github.com/ouz/goboilerplate/internal/domain/shared"
	"github.com/ouz/goboilerplate/internal/domain/user"
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
		{
			name: "Valid user",
			args: args{
				username: "testuser",
				password: "validpass123",
				email:    vo.Email{Address: "test@example.com"},
			},
			wantErr: false,
		},
		{
			name: "Empty username",
			args: args{
				username: "",
				password: "validpass123",
				email:    vo.Email{Address: "test@example.com"},
			},
			wantErr: true,
		},
		{
			name: "Too short username",
			args: args{
				username: "ab",
				password: "validpass123",
				email:    vo.Email{Address: "test@example.com"},
			},
			wantErr: true,
		},
		{
			name: "Invalid password",
			args: args{
				username: "testuser",
				password: "short",
				email:    vo.Email{Address: "test@example.com"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := user.NewUser(tt.args.username, tt.args.password, tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == nil {
					t.Error("NewUser() = nil, want non-nil")
					return
				}
				if got.Username != tt.args.username {
					t.Errorf("NewUser().Username = %v, want %v", got.Username, tt.args.username)
				}
				if got.Email != tt.args.email.Address {
					t.Errorf("NewUser().Email = %v, want %v", got.Email, tt.args.email.Address)
				}
				if len(got.Roles) != 1 || got.Roles[0].Name != user.UserRoleUser {
					t.Error("NewUser() should have exactly one USER role")
				}
				if len(got.Credentials) != 1 {
					t.Error("NewUser() should have exactly one credential")
				}
				if len(got.Confirmations) != 1 {
					t.Error("NewUser() should have exactly one confirmation")
				}
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
		{
			name:    "Create anonymous user",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := user.NewAnonymousUser()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAnonymousUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == nil {
					t.Error("NewAnonymousUser() = nil, want non-nil")
					return
				}
				if !got.Anonymous {
					t.Error("NewAnonymousUser().Anonymous = false, want true")
				}
				if !got.Enabled {
					t.Error("NewAnonymousUser().Enabled = false, want true")
				}
				if !got.Verified {
					t.Error("NewAnonymousUser().Verified = false, want true")
				}
				if len(got.Roles) != 1 || got.Roles[0].Name != user.UserRoleAnonymous {
					t.Error("NewAnonymousUser() should have exactly one ANONYMOUS role")
				}
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
		{
			name: "Add valid credential",
			u: &user.User{
				ID: uuid.New().String(),
			},
			args: args{
				password: "validpass123",
			},
			wantErr: false,
		},
		{
			name: "Add invalid credential",
			u: &user.User{
				ID: uuid.New().String(),
			},
			args: args{
				password: "short",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initialCredCount := len(tt.u.Credentials)
			err := tt.u.AddCredential(tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("User.AddCredential() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && len(tt.u.Credentials) != initialCredCount+1 {
				t.Error("User.AddCredential() did not add credential")
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
		{
			name: "Add confirmation",
			u: &user.User{
				ID: uuid.New().String(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initialConfCount := len(tt.u.Confirmations)
			err := tt.u.AddConfirmation()
			if (err != nil) != tt.wantErr {
				t.Errorf("User.AddConfirmation() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && len(tt.u.Confirmations) != initialConfCount+1 {
				t.Error("User.AddConfirmation() did not add confirmation")
			}
		})
	}
}

func TestUser_IsPasswordValid(t *testing.T) {
	validPassword := "validpass123"
	u := &user.User{
		ID: uuid.New().String(),
	}
	_ = u.AddCredential(validPassword)

	type args struct {
		password string
	}
	tests := []struct {
		name string
		u    *user.User
		args args
		want bool
	}{
		{
			name: "Valid password",
			u:    u,
			args: args{
				password: validPassword,
			},
			want: true,
		},
		{
			name: "Invalid password",
			u:    u,
			args: args{
				password: "wrongpass123",
			},
			want: false,
		},
		{
			name: "Empty password",
			u:    u,
			args: args{
				password: "",
			},
			want: false,
		},
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
	u := &user.User{
		ID: uuid.New().String(),
	}
	role, _ := user.NewUserRole(u.ID, user.UserRoleUser)
	u.Roles = []user.UserRole{*role}

	type args struct {
		role user.UserRoleName
	}
	tests := []struct {
		name string
		u    *user.User
		args args
		want bool
	}{
		{
			name: "Has role",
			u:    u,
			args: args{
				role: user.UserRoleUser,
			},
			want: true,
		},
		{
			name: "Does not have role",
			u:    u,
			args: args{
				role: user.UserRoleAnonymous,
			},
			want: false,
		},
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
		{
			name: "Confirm user",
			u: &user.User{
				ID:       uuid.New().String(),
				Enabled:  false,
				Verified: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.u.Confirm()
			if !tt.u.Enabled {
				t.Error("User.Confirm() did not enable user")
			}
			if !tt.u.Verified {
				t.Error("User.Confirm() did not verify user")
			}
		})
	}
}
