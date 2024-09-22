package entities

import "testing"

func TestHideDBpass(t *testing.T) {
	tests := []struct {
		name  string
		dbURL string
		want  string
	}{
		// TODO: Add test cases.
		{
			name:  "Test 1: default url",
			dbURL: "postgres://username:123456@127.0.0.1:5432/dbname",
			want:  "postgres://username:******@127.0.0.1:5432/dbname",
		},
		{
			name:  "Test 2: url with quotes",
			dbURL: "postgres://username:\"123456\"@127.0.0.1:5432/dbname",
			want:  "postgres://username:******@127.0.0.1:5432/dbname",
		},
		{
			name:  "Test 3: pass with quotes, : and @",
			dbURL: "postgres://username:\"12:@3456\"@127.0.0.1:5432/dbname",
			want:  "postgres://username:******@127.0.0.1:5432/dbname",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HideDBpass(tt.dbURL); got != tt.want {
				t.Errorf("HideDBpass() = %v, want %v", got, tt.want)
			}
		})
	}
}
