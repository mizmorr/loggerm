package logger

import (
	"context"
)

func GetLoggerFromContext(ctx context.Context) *Logger {
	log, ok := ctx.Value("logger").(*Logger)
	if !ok {
		return Get("info")
	}
	return log
}
