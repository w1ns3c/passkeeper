package client

import (
	"flag"
	"net"
	"os"
	"passkeeper/internal/config"
	"slices"
	"strconv"
)

type Args struct {
	Addr     string
	LogLevel string
	SyncTime int // inseconds
}

func CliParseArgs() *Args {
	var (
		flagAddrVal     string
		flagLogLevelVal string
		flagSyncTime    int

		realAddr        string
		realLogLevel    string
		realSyncTimeVal int

		addrIsSet  bool
		levelIsSet bool
		syncIsSet  bool

		defaultAddr     = config.DefaultAddr
		defaultLogLevel = config.Level
		defaultSyncTime = int(config.SyncDefault.Seconds())

		addrFlag     = "addr"
		levelFlag    = "level"
		syncTimeFlag = "time"
	)

	inArgs := os.Args
	addrIsSet = slices.Contains(inArgs, "-"+addrFlag)
	levelIsSet = slices.Contains(inArgs, "-"+levelFlag)
	syncIsSet = slices.Contains(inArgs, "-"+syncTimeFlag)

	flag.StringVar(&flagAddrVal, addrFlag, defaultAddr, "client connect address")
	flag.StringVar(&flagLogLevelVal, levelFlag, defaultLogLevel, "log level")
	flag.IntVar(&flagSyncTime, syncTimeFlag, defaultSyncTime, "time to sync credentials in seconds")
	flag.Parse()

	// Read ENV only if flags not set !!!
	envAddr, existsAddr := os.LookupEnv("CLI_ADDR")
	envLogLevel, existsLvl := os.LookupEnv("CLI_LOG_LEVEL")
	envST, existST := os.LookupEnv("CLI_SYNC_TIME")

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

	// validate
	if !levelIsSet && existsLvl {
		realLogLevel = envLogLevel // FLAG value error, ENV is good
	} else {
		realLogLevel = flagLogLevelVal // FLAG value is good
	}

	// validate sync time
	envSTint, err := strconv.Atoi(envST)
	if !syncIsSet && existST && err == nil {
		realSyncTimeVal = envSTint // FLAG value error, ENV is good
	} else {
		realSyncTimeVal = flagSyncTime // FLAG value is good
	}

	return &Args{
		Addr:     realAddr,
		LogLevel: realLogLevel,
		SyncTime: realSyncTimeVal,
	}

}
