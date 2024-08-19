package logger

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"os"
	"strings"
	"time"
)

func Init(level string) *zerolog.Logger {
	//zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	lvl := SelectLevel(level)

	logger := zerolog.New(os.Stderr).With().
		Timestamp().Logger().Level(lvl)

	logger = logger.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.DateTime,
	})

	return &logger
}

func SelectLevel(level string) zerolog.Level {
	level = strings.ToTitle(level)
	switch level {
	// INFO
	case "INF":
		return zerolog.InfoLevel
	case "INFO":
		return zerolog.InfoLevel

	// DEBUG
	case "DBG":
		return zerolog.DebugLevel
	case "DEBUG":
		return zerolog.DebugLevel

	// Warning
	case "WRN":
		return zerolog.WarnLevel
	case "WARNING":
		return zerolog.WarnLevel
	case "WARN":
		return zerolog.WarnLevel

	// Error level
	case "ERR":
		return zerolog.ErrorLevel
	case "EROR":
		return zerolog.ErrorLevel

	default:
		return zerolog.DebugLevel
	}
}
