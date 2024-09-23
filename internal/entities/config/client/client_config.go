package client

import (
	"flag"
	"net"
	"os"
	"slices"
	"strconv"

	"passkeeper/internal/entities/config"
)

// Args struct to save Client params
type Args struct {
	Addr     string
	LogLevel string
	LogFile  string
	SyncTime int // inseconds
}

// CliParseArgs parse Client params from command line
func CliParseArgs() *Args {
	var (
		flagAddrVal     string
		flagLogLevelVal string
		flagLogFile     string
		flagSyncTime    int

		realAddr        string
		realLogLevel    string
		realLogFile     string
		realSyncTimeVal int

		defaultAddr     = config.DefaultAddr
		defaultLogLevel = config.Level
		defaultLogPath  = config.CliLogFilePath
		defaultSyncTime = int(config.SyncDefault.Seconds())

		addrFlag     = "addr"
		levelFlag    = "level"
		logFlag      = "log"
		syncTimeFlag = "time"
	)

	inArgs := os.Args
	addrIsSet := slices.Contains(inArgs, "-"+addrFlag)
	levelIsSet := slices.Contains(inArgs, "-"+levelFlag)
	logIsSet := slices.Contains(inArgs, "-"+logFlag)
	syncIsSet := slices.Contains(inArgs, "-"+syncTimeFlag)

	flag.StringVar(&flagAddrVal, addrFlag, defaultAddr, "client connect address")
	flag.StringVar(&flagLogLevelVal, levelFlag, defaultLogLevel, "log level")
	flag.StringVar(&flagLogFile, logFlag, defaultLogPath, "log file path")
	flag.IntVar(&flagSyncTime, syncTimeFlag, defaultSyncTime, "time to sync credentials in seconds")
	flag.Parse()

	// Read ENV only if flags not set !!!
	envAddr, existsAddr := os.LookupEnv("CLI_ADDR")
	envLogLevel, existsLvl := os.LookupEnv("CLI_LOG_LEVEL")
	envLogFile, existLF := os.LookupEnv("CLI_LOG_FILE")
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
	realLogLevel = flagLogLevelVal // possible default FLAG value | FLAG good value
	if !levelIsSet && existsLvl {
		realLogLevel = envLogLevel // FLAG value error, ENV is good
	}

	// validate sync time
	envSTint, err := strconv.Atoi(envST)
	if !syncIsSet && existST && err == nil {
		realSyncTimeVal = envSTint // FLAG value error, ENV is good
	} else {
		realSyncTimeVal = flagSyncTime // FLAG value is good
	}

	realLogFile = flagLogFile // possible default FLAG value | FLAG good value
	if !logIsSet && existLF {
		realLogFile = envLogFile // FLAG value error, ENV is good
	} else {
		realLogFile = flagLogFile
	}

	return &Args{
		Addr:     realAddr,
		LogLevel: realLogLevel,
		LogFile:  realLogFile,
		SyncTime: realSyncTimeVal,
	}

}
