package user

import (
	"testing"

	"github.com/ouz/goboilerplate/internal/domain/user"
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
		{
			name: "Valid password",
			args: args{
				plaintext: "validpassword123",
			},
			wantErr: false,
		},
		{
			name: "Empty password",
			args: args{
				plaintext: "",
			},
			wantErr: true,
		},
		{
			name: "Too short password",
			args: args{
				plaintext: "short",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := user.NewPassword(tt.args.plaintext)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("NewPassword() = nil, want non-nil")
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
		{
			name: "Correct password",
			args: args{
				plaintext: "mypassword123",
			},
			want: true,
		},
		{
			name: "Incorrect password",
			args: args{
				plaintext: "wrongpassword",
			},
			want: false,
		},
		{
			name: "Empty password",
			args: args{
				plaintext: "",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.want {
				var err error
				tt.p, err = user.NewPassword(tt.args.plaintext)
				if err != nil {
					t.Fatalf("Failed to create test password: %v", err)
				}
			} else {
				var err error
				tt.p, err = user.NewPassword("mypassword123")
				if err != nil {
					t.Fatalf("Failed to create test password: %v", err)
				}
			}

			if got := tt.p.Verify(tt.args.plaintext); got != tt.want {
				t.Errorf("Password.Verify() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPassword_Hashed(t *testing.T) {
	tests := []struct {
		name     string
		password string
	}{
		{
			name:     "Regular password",
			password: "mypassword123",
		},
		{
			name:     "Long password",
			password: "verylongpasswordwithmanycharacters123!@#",
		},
		{
			name:     "Password with special chars",
			password: "pass!@#$%^&*()",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := user.NewPassword(tt.password)
			if err != nil {
				t.Fatalf("Failed to create test password: %v", err)
			}
			got := p.Hashed()
			if got == "" {
				t.Error("Password.Hashed() = empty string, want non-empty hash")
			}
			if got == tt.password {
				t.Error("Password.Hashed() returned plaintext password")
			}
		})
	}
}
