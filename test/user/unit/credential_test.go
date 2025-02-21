package user

import (
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
		{
			name: "Valid password credential",
			args: args{
				credentialType: user.CredentialTypePassword,
				secret:         "validpassword123",
			},
			wantErr: false,
		},
		{
			name: "Invalid credential type",
			args: args{
				credentialType: "INVALID_TYPE",
				secret:         "validpassword123",
			},
			wantErr: true,
		},
		{
			name: "Empty password",
			args: args{
				credentialType: user.CredentialTypePassword,
				secret:         "",
			},
			wantErr: true,
		},
		{
			name: "Too short password",
			args: args{
				credentialType: user.CredentialTypePassword,
				secret:         "short",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := user.NewCredential(tt.args.credentialType, tt.args.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCredential() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == nil {
					t.Error("NewCredential() = nil, want non-nil")
					return
				}
				if got.CredentialType != tt.args.credentialType {
					t.Errorf("NewCredential().CredentialType = %v, want %v", got.CredentialType, tt.args.credentialType)
				}
				if got.Hash == "" {
					t.Error("NewCredential().Hash is empty")
				}
				if got.Hash == tt.args.secret {
					t.Error("NewCredential().Hash equals plaintext secret")
				}
				if got.CreatedAt.IsZero() {
					t.Error("NewCredential().CreatedAt is zero")
				}
				if got.UpdatedAt.IsZero() {
					t.Error("NewCredential().UpdatedAt is zero")
				}
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
		{
			name: "Valid password",
			args: args{
				password: "mypassword123",
			},
			want: true,
		},
		{
			name: "Invalid password",
			args: args{
				password: "wrongpassword",
			},
			want: false,
		},
		{
			name: "Empty password",
			args: args{
				password: "",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a credential for testing
			var err error
			if tt.want {
				tt.c, err = user.NewCredential(user.CredentialTypePassword, tt.args.password)
			} else {
				tt.c, err = user.NewCredential(user.CredentialTypePassword, "mypassword123")
			}
			if err != nil {
				t.Fatalf("Failed to create test credential: %v", err)
			}

			if got := tt.c.IsPasswordValid(tt.args.password); got != tt.want {
				t.Errorf("Credential.IsPasswordValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
