package logger

import (
	"fmt"
	"io"
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

var path string = ""

func Get(logPath, logLevel string) *Logger {
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	if err != nil {
		return nil
	}
	path = logPath

	defer file.Close()
	once.Do(func() {
		writer := newConsoleWriter()
		multiWriter := io.MultiWriter(writer, file)
		zeroLogger := zerolog.New(multiWriter).With().Logger()
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
	levelColors := map[zerolog.Level]string{
		zerolog.InfoLevel:  "\033[34m", // Синий
		zerolog.WarnLevel:  "\033[33m", // Жёлтый
		zerolog.ErrorLevel: "\033[31m", // Красный
		zerolog.DebugLevel: "\033[32m", // Зелёный
	}
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
	writer.FormatLevel = func(i interface{}) string {
		if l, ok := i.(string); ok {
			colorLog, _ := zerolog.ParseLevel(l)
			color := levelColors[colorLog]
			return color + l + "\033[0m"
		}
		return i.(string)
	}

	return &writer
}
