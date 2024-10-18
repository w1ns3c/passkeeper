//go:build windows
// +build windows

package config

var (
	// path to log file on Windows machines
	TmpDir    = "C:\\Windows\\Temp\\"
	CliLogDir = TmpDir
)
