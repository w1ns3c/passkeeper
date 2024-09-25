//go:build linux
// +build linux

package config

var (
	CliLogDir = "/var/log/passkeeper.log" // path to log file on Linux machines
	TmpDir    = "/tmp/"
)
