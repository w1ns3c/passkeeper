package server

import (
	"os"
	"reflect"
	"testing"
)

func TestSrvParseArgs(t *testing.T) {
	var (
		addrFlag  = "-addr"
		levelFlag = "-level"
		dbUrl     = "-db"
	)

	tests := []struct {
		name  string
		flags []string
		want  *Args
	}{
		{
			name: "Test 1: setup all values by flags",
			flags: []string{
				addrFlag, "127.0.0.1:9000",
				levelFlag, "DEBUG",
				dbUrl, "postgres://UseR:Passwd@127.0.0.1:11111/myDB",
			},
			want: &Args{
				Addr:     "127.0.0.1:9000",
				LogLevel: "DEBUG",
				DBurl:    "postgres://UseR:Passwd@127.0.0.1:11111/myDB",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = append(os.Args, tt.flags...)
			if got := SrvParseArgs(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SrvParseArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}
