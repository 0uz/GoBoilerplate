package shared

import (
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
		{
			name: "Valid email",
			args: args{
				address: "test@example.com",
			},
			want: shared.Email{
				Address: "test@example.com",
			},
			wantErr: false,
		},
		{
			name: "Empty email",
			args: args{
				address: "",
			},
			wantErr: true,
		},
		{
			name: "Invalid email - no @",
			args: args{
				address: "testexample.com",
			},
			wantErr: true,
		},
		{
			name: "Invalid email - no domain",
			args: args{
				address: "test@",
			},
			wantErr: true,
		},
		{
			name: "Invalid email - no local part",
			args: args{
				address: "@example.com",
			},
			wantErr: true,
		},
		{
			name: "Invalid email - special chars",
			args: args{
				address: "test!@example.com",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := shared.NewEmail(tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.Address != tt.want.Address {
				t.Errorf("NewEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}
