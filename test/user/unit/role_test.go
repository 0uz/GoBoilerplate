package user

import (
	"testing"

	"github.com/google/uuid"
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
		{
			name: "Valid regular user role",
			args: args{
				userID: uuid.New().String(),
				name:   user.UserRoleUser,
			},
			wantErr: false,
		},
		{
			name: "Valid anonymous user role",
			args: args{
				userID: uuid.New().String(),
				name:   user.UserRoleAnonymous,
			},
			wantErr: false,
		},
		{
			name: "Invalid role name",
			args: args{
				userID: uuid.New().String(),
				name:   "INVALID_ROLE",
			},
			wantErr: true,
		},
		{
			name: "Empty user ID",
			args: args{
				userID: "",
				name:   user.UserRoleUser,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := user.NewUserRole(tt.args.userID, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewUserRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == nil {
					t.Error("NewUserRole() = nil, want non-nil")
					return
				}
				if got.UserID != tt.args.userID {
					t.Errorf("NewUserRole().UserID = %v, want %v", got.UserID, tt.args.userID)
				}
				if got.Name != tt.args.name {
					t.Errorf("NewUserRole().Name = %v, want %v", got.Name, tt.args.name)
				}
				if got.CreatedAt.IsZero() {
					t.Error("NewUserRole().CreatedAt is zero")
				}
			}
		})
	}
}
