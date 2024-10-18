package client

import (
	"os"
	"reflect"
	"testing"
)

func TestCliParseArgs1(t *testing.T) {

	var (
		addrFlag     = "-addr"
		levelFlag    = "-level"
		logFlag      = "-log"
		syncTimeFlag = "-time"
		versionFlag  = "-version"
	)

	tests := []struct {
		name  string
		flags []string
		want  *Args
	}{
		{
			name: "Test 1",
			flags: []string{
				addrFlag, "127.0.0.1:9000",
				levelFlag, "DEBUG",
				logFlag, "/tmp/123",
				syncTimeFlag, "100",
				versionFlag, "myversion",
			},
			want: &Args{
				Addr:        "127.0.0.1:9000",
				LogLevel:    "DEBUG",
				LogFile:     "/tmp/123",
				SyncTime:    100,
				ShowVersion: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = append(os.Args, tt.flags...)
			if got := CliParseArgs(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CliParseArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}
