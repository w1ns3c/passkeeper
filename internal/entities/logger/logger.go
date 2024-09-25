package logger

import (
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

// Init setup logger output
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

// InitFile setup logger output to file
func InitFile(level, filepath string) *zerolog.Logger {
	//zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	lvl := SelectLevel(level)

	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil
	}

	logger := zerolog.New(f).With().
		Timestamp().Logger().Level(lvl)

	logger = logger.Output(zerolog.ConsoleWriter{
		Out:        f,
		TimeFormat: time.DateTime,
	})

	return &logger
}

// SelectLevel choose level of logging, return specific zerolog.Level
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
	case "ERROR":
		return zerolog.ErrorLevel

	default:
		return zerolog.DebugLevel
	}
}
