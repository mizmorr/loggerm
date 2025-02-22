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
	path = logPath

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	if err != nil {
		return nil
	}
	once.Do(func() {
		writerToStd := newWriter(os.Stdout, false)

		writerToFile := newWriter(file, true)

		multiWriter := io.MultiWriter(writerToStd, writerToFile)
		zeroLogger := zerolog.New(multiWriter).With().Caller().Logger()
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
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		}
		logger = Logger{&zeroLogger}
	})
	return &logger
}

func newWriter(placeToWrite *os.File, isFile bool) *zerolog.ConsoleWriter {
	levelColors := map[zerolog.Level]string{
		zerolog.TraceLevel: "\033[37m", // Белый
		zerolog.DebugLevel: "\033[36m", // Голубой
		zerolog.InfoLevel:  "\033[34m", // Синий
		zerolog.WarnLevel:  "\033[33m", // Жёлтый
		zerolog.ErrorLevel: "\033[31m", // Красный
		zerolog.FatalLevel: "\033[35m", // Пурпурный
		zerolog.PanicLevel: "\033[41m", // Красный фон
	}

	writer := zerolog.ConsoleWriter{
		Out: placeToWrite,

		FormatLevel: func(i interface{}) string {
			str := strings.ToUpper(fmt.Sprintf("[%s]", i))
			if !isFile {
				if l, ok := i.(string); ok {
					level, _ := zerolog.ParseLevel(l)
					color := levelColors[level]
					return color + str + "\033[0m"
				}
			}
			return str
		},
		FormatMessage: func(i interface{}) string {
			return fmt.Sprintf("%s ", i)
		},
		FormatTimestamp: func(i interface{}) string {
			return fmt.Sprintf("%v |", time.Now().Format(time.RFC822))
		},
		PartsExclude: []string{zerolog.TimeFieldFormat},
	}

	if isFile {
		writer.NoColor = true
		writer.FormatCaller = func(i interface{}) string {
			return fmt.Sprintf("| %s |", i.(string)) // Кастомизация caller
		}
	}
	return &writer
}
