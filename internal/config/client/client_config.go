package client

import (
	"flag"
	"net"
	"net/netip"
	"os"
	"slices"
)

type Args struct {
	Addr     string
	LogLevel string
}

func CliParseArgs() *Args {
	var (
		flagAddrVal     string
		flagLogLevelVal string
		realAddr        string
		realLogLevel    string
		addrIsSet       bool
		levelIsSet      bool

		defaultAddr     = "localhost:8000"
		defaultLogLevel = "info"

		addrFlag  = "addr"
		levelFlag = "level"
	)

	inArgs := os.Args
	addrIsSet = slices.Contains(inArgs, "-"+addrFlag)
	levelIsSet = slices.Contains(inArgs, "-"+levelFlag)

	flag.StringVar(&flagAddrVal, addrFlag, defaultAddr, "client connect address")
	flag.StringVar(&flagLogLevelVal, levelFlag, defaultLogLevel, "log level")
	flag.Parse()

	// Read ENV only if flags not set !!!
	envAddr, existsAddr := os.LookupEnv("CLI_ADDR")
	envLogLevel, existsLvl := os.LookupEnv("CLI_LOG_LEVEL")

	// validate flag and env values
	// validate net address
	_, err := net.ResolveTCPAddr("tcp", flagAddrVal)
	if !addrIsSet || err != nil {
		if existsAddr {
			// check ENV add value
			_, err = netip.ParseAddrPort(envAddr)
			if err != nil {
				realAddr = defaultAddr // FLAG && ENV values error
			} else {
				realAddr = envAddr // FLAG value error, ENV is good
			}
		}
	} else {
		realAddr = flagAddrVal // FLAG value is good
	}

	// validate
	if !levelIsSet {
		if existsLvl {
			realAddr = envLogLevel // FLAG value error, ENV is good
		}
	} else {
		realAddr = flagAddrVal // FLAG value is good
	}

	return &Args{
		Addr:     realAddr,
		LogLevel: realLogLevel,
	}

}
