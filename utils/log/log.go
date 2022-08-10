package log

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"multi-bot/config"
	"os"
	"runtime"
	"strconv"
	"time"
)

const (
	traceId = "traceid"
)

func init() {
	if config.EnvIsDev() {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	if config.UseConsoleWrite() {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	zerolog.TimeFieldFormat = time.RFC3339Nano
}

func Panic() *zerolog.Event {
	_, file, line, ok := runtime.Caller(1)
	e := log.Panic()
	if ok {
		e = e.Str("line", file+":"+strconv.Itoa(line))
	}
	return e
}

func Error() *zerolog.Event {
	_, file, line, ok := runtime.Caller(1)
	e := log.Error().Stack()
	if ok {
		e = e.Str("line", file+":"+strconv.Itoa(line))
	}
	return e
}

func Debug() *zerolog.Event {
	_, file, line, ok := runtime.Caller(1)
	e := log.Debug()
	if ok {
		e = e.Str("line", file+":"+strconv.Itoa(line))
	}
	return e
}

func Warn() *zerolog.Event {
	_, file, line, ok := runtime.Caller(1)
	e := log.Warn()
	if ok {
		e = e.Str("line", file+":"+strconv.Itoa(line))
	}
	return e
}

func Info() *zerolog.Event {
	_, file, line, ok := runtime.Caller(1)
	e := log.Info()
	if ok {
		e = e.Str("line", file+":"+strconv.Itoa(line))
	}
	return e
}
