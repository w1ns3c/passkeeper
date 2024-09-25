package logger

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/stretchr/testify/require"
)

func TestSelectLevel(t *testing.T) {
	tests := []struct {
		name  string
		level string
		want  zerolog.Level
	}{
		{
			name:  "Test 1.1: debug",
			level: "debug",
			want:  zerolog.DebugLevel,
		},
		{
			name:  "Test 1.2: dbg",
			level: "dbg",
			want:  zerolog.DebugLevel,
		},
		{
			name:  "Test 2.1: error",
			level: "error",
			want:  zerolog.ErrorLevel,
		},
		{
			name:  "Test 2.2: err",
			level: "err",
			want:  zerolog.ErrorLevel,
		},
		{
			name:  "Test 3.1: warning",
			level: "warning",
			want:  zerolog.WarnLevel,
		},
		{
			name:  "Test 3.2: wrn",
			level: "wrn",
			want:  zerolog.WarnLevel,
		},
		{
			name:  "Test 3.3: warn",
			level: "Warn",
			want:  zerolog.WarnLevel,
		},
		{
			name:  "Test 4.1: info",
			level: "inf",
			want:  zerolog.InfoLevel,
		},
		{
			name:  "Test 4.2: info",
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

func TestInit(t *testing.T) {

	level := "DBG"
	lvl := SelectLevel(level)

	logger := zerolog.New(os.Stderr).With().
		Timestamp().Logger().Level(lvl)

	logger = logger.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.DateTime,
	})

	tests := []struct {
		name  string
		level string
		want  *zerolog.Logger
	}{
		// TODO: Add test cases.
		{
			name:  "Test 1: success",
			level: level,
			want:  &logger,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Init(tt.level); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Init() = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestInitFile(t *testing.T) {
	type args struct {
		level    string
		filepath string
	}

	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	level := "DBG"
	filepath := "/tmp/123456"

	lvl := SelectLevel(level)

	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	logger := zerolog.New(f).With().
		Timestamp().Logger().Level(lvl)

	logger = logger.Output(zerolog.ConsoleWriter{
		Out:        f,
		TimeFormat: time.DateTime,
	})

	tests := []struct {
		name string
		args args
		want *zerolog.Logger
	}{
		// TODO: Add test cases.
		{
			name: "Test 1: success",
			args: args{
				level:    level,
				filepath: filepath,
			},
			want: &logger,
		},
		{
			name: "Test 2: errror",
			args: args{
				level:    level,
				filepath: "/t111/123",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := InitFile(tt.args.level, tt.args.filepath)
			if (got == nil) == (tt.want == nil) {
				return
			}
			require.Equal(t, tt.want.GetLevel(), got.GetLevel())
		})
	}
}
