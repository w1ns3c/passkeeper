package cli

import "testing"

func TestFilterEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "Valid simple email",
			email:   "email@email.com",
			wantErr: false,
		},
		{
			name:    "Valid hard email",
			email:   "email44.11-email.me@email99.corp.com",
			wantErr: false,
		},
		{
			name:    "Invalid email (dot end username)",
			email:   "email.@email.com",
			wantErr: true,
		},
		{
			name:    "Invalid email (short email code)",
			email:   "email@email.c",
			wantErr: true,
		},
		{
			name:    "Invalid email (double dots username)",
			email:   "email..email@email.com",
			wantErr: true,
		},
		{
			name:    "Invalid email (double dots email)",
			email:   "email@email..c",
			wantErr: true,
		},
		{
			name:    "Invalid email (double @)",
			email:   "email@alex@email.com",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := FilterEmail(tt.email); (err != nil) != tt.wantErr {
				t.Errorf("FilterEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
