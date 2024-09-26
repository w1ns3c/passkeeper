//go:build linux
// +build linux

package config

var (
	CliLogDir = "/var/log/" // path to log file on Linux machines
	TmpDir    = "/tmp/"
)
