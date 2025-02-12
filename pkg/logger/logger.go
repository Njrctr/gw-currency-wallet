package logger

import (
	"log/slog"
	"os"
)

func InitLogger(logLevel string) *slog.Logger {

	switch logLevel {
	case "debug":
		slog.SetLogLoggerLevel(slog.LevelDebug)
	case "info":
		slog.SetLogLoggerLevel(slog.LevelInfo)
	default:
		slog.SetLogLoggerLevel(slog.LevelInfo)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	return logger
}
