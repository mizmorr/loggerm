package logger

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

type Logger struct {
	*zerolog.Logger
}

var (
	logger Logger
	once   sync.Once
)

func Get(logLevel string) *Logger {
	once.Do(func() {
		writer := newConsoleWriter()
		zeroLogger := zerolog.New(writer).With().Logger()
		switch logLevel {
		case "debug":
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		case "info":
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		case "warn", "warning":
			zerolog.SetGlobalLevel(zerolog.WarnLevel)
		case "err", "error":
			zerolog.SetGlobalLevel(zerolog.ErrorLevel)
		case "fatal":
			zerolog.SetGlobalLevel(zerolog.FatalLevel)
		case "panic":
			zerolog.SetGlobalLevel(zerolog.PanicLevel)
		default:
			zerolog.SetGlobalLevel(zerolog.InfoLevel) // log info and above by default
		}
		logger = Logger{&zeroLogger}
	})
	return &logger
}

func newConsoleWriter() *zerolog.ConsoleWriter {
	writer := zerolog.ConsoleWriter{
		Out: os.Stderr,
		// TimeFormat: time.RFC1123,
		FormatLevel: func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("[%s]", i))
		},
		FormatMessage: func(i interface{}) string {
			return fmt.Sprintf("| %s ", i)
		},
		FormatTimestamp: func(i interface{}) string {
			return fmt.Sprintf("%v |", time.Now().Format(time.RFC822))
		},
		PartsExclude: []string{
			zerolog.TimeFieldFormat,
		},
	}
	return &writer
}
