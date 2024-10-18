package server

import (
	"flag"
	"net"
	"os"
	"slices"

	"passkeeper/internal/entities/config"
)

// Args struct to save Server params
type Args struct {
	Addr     string
	LogLevel string
	DBurl    string
}

// SrvParseArgs parse Server params from command line
func SrvParseArgs() *Args {
	var (
		flagAddrVal     string
		flagLogLevelVal string
		flagDBurl       string

		realAddr     string
		realLogLevel string
		realDB       string

		addrIsSet  bool
		levelIsSet bool

		defaultAddr     = config.DefaultAddr
		defaultLogLevel = config.Level

		addrFlag  = "addr"
		levelFlag = "level"
		dbUrl     = "db"
	)

	inArgs := os.Args
	addrIsSet = slices.Contains(inArgs, "-"+addrFlag)
	levelIsSet = slices.Contains(inArgs, "-"+levelFlag)

	flag.StringVar(&flagAddrVal, addrFlag, defaultAddr, "server listening address")
	flag.StringVar(&flagDBurl, dbUrl, "", "db url, ex: \"postgres://user:pass@host:port/dbname\"")
	flag.StringVar(&flagLogLevelVal, levelFlag, defaultLogLevel, "log level")
	flag.Parse()

	// Read ENV only if flags not set !!!
	envAddr, existsAddr := os.LookupEnv("SRV_ADDR")
	envLogLevel, existsLvl := os.LookupEnv("SRV_LOG_LEVEL")
	envDB, existsDB := os.LookupEnv("SRV_DB")

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

	if existsDB && flagDBurl == "" {
		realDB = envDB
	} else {
		realDB = flagDBurl
	}

	return &Args{
		Addr:     realAddr,
		LogLevel: realLogLevel,
		DBurl:    realDB,
	}

}
