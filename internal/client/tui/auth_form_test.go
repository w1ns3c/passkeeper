package tui

import (
	"testing"
)

func TestFilterResource(t *testing.T) {

	tests := []struct {
		name  string
		input string
		want  string
	}{
		// TODO: Add test cases.
		{
			name:  "Filter \"https://\"",
			input: "https://somesitewithpath.org.com/path/like/this?param1=value",
			want:  "somesitewithpath.org.com",
		},
		{
			name:  "Filter \"http://\"",
			input: "http://somesitewithpath.org.com/path/like/this?param1=value",
			want:  "somesitewithpath.org.com",
		},
		{
			name:  "Filter uri path only",
			input: "somesitewithpath.org.com/path/like/this?param1=value",
			want:  "somesitewithpath.org.com",
		},
		{
			name:  "Filter mistake (https:/)",
			input: "https:/somesitewithpath.org.com/path/like/this?param1=value",
			want:  "somesitewithpath.org.com",
		},
		{
			name:  "Filter mistake (http:/)",
			input: "http:/somesitewithpath.org.com/path/like/this?param1=value",
			want:  "somesitewithpath.org.com",
		},
		{
			name:  "Filter mistake (https:/)",
			input: "https:/somesitewithpath.org.com/path/like/this?param1=value",
			want:  "somesitewithpath.org.com",
		},
		{
			name:  "Filter mistake other proto (ftp:/)",
			input: "ftp:/somesitewithpath.org.com/path/like/this?param1=value",
			want:  "somesitewithpath.org.com",
		},
		{
			name:  "Filter other proto",
			input: "anyprotohere://somesitewithpath.org.com/path/like/this?param1=value",
			want:  "somesitewithpath.org.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FilterResource(tt.input); got != tt.want {
				t.Errorf("FilterResource() = %v, want %v", got, tt.want)
			}
		})
	}
}
