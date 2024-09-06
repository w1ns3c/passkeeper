package server

import (
	"flag"
	"net"
	"os"
	"passkeeper/internal/config"
	"slices"
)

type Args struct {
	Addr     string
	LogLevel string
}

func SrvParseArgs() *Args {
	var (
		flagAddrVal     string
		flagLogLevelVal string

		realAddr     string
		realLogLevel string

		addrIsSet  bool
		levelIsSet bool

		defaultAddr     = config.DefaultAddr
		defaultLogLevel = config.Level

		addrFlag  = "addr"
		levelFlag = "level"
	)

	inArgs := os.Args
	addrIsSet = slices.Contains(inArgs, "-"+addrFlag)
	levelIsSet = slices.Contains(inArgs, "-"+levelFlag)

	flag.StringVar(&flagAddrVal, addrFlag, defaultAddr, "server listening address")
	flag.StringVar(&flagLogLevelVal, levelFlag, defaultLogLevel, "log level")
	flag.Parse()

	// Read ENV only if flags not set !!!
	envAddr, existsAddr := os.LookupEnv("SRV_ADDR")
	envLogLevel, existsLvl := os.LookupEnv("SRV_LOG_LEVEL")

	// validate flag and env values
	// validate net address
	_, err := net.ResolveTCPAddr("tcp", flagAddrVal)
	if (!addrIsSet || err != nil) && existsAddr {
		// check ENV add value
		_, err = net.ResolveTCPAddr("tcp", envAddr)
		if err != nil {
			realAddr = defaultAddr // FLAG && ENV values error
		} else {
			realAddr = envAddr // FLAG value error, ENV is good
		}
	} else {
		realAddr = flagAddrVal // FLAG value is good
	}

	// validate log level
	if !levelIsSet && existsLvl {
		realLogLevel = envLogLevel // FLAG value error, ENV is good
	} else {
		realLogLevel = flagLogLevelVal // FLAG value is good
	}

	return &Args{
		Addr:     realAddr,
		LogLevel: realLogLevel,
	}

}
