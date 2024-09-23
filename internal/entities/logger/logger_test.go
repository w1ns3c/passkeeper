package logger

import (
	"testing"

	"github.com/rs/zerolog"
)

func TestSelectLevel(t *testing.T) {
	tests := []struct {
		name  string
		level string
		want  zerolog.Level
	}{
		{
			name:  "Test 1: debug",
			level: "debug",
			want:  zerolog.DebugLevel,
		},
		{
			name:  "Test 2: error",
			level: "error",
			want:  zerolog.ErrorLevel,
		},
		{
			name:  "Test 3: warning",
			level: "warning",
			want:  zerolog.WarnLevel,
		},
		{
			name:  "Test 4: info",
			level: "info",
			want:  zerolog.InfoLevel,
		},
		{
			name:  "Test 5: unknown",
			level: "unknown",
			want:  zerolog.DebugLevel,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SelectLevel(tt.level); got != tt.want {
				t.Errorf("SelectLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}
