package logger

import (
	"log/slog"
	"os"
)

type Logger interface {
	Info(msg string)
	Error(msg string)
}

type slogLogger struct {
	logger *slog.Logger
}

func New(level string) Logger {
	levelByString := map[string]slog.Level{
		"debug": slog.LevelDebug,
		"info":  slog.LevelInfo,
		"warn":  slog.LevelWarn,
		"error": slog.LevelError,
	}
	loggerLevel, ok := levelByString[level]
	if !ok {
		panic("Wrong log level")
	}
	opts := &slog.HandlerOptions{
		Level: loggerLevel,
	}
	return slogLogger{
		slog.New(slog.NewTextHandler(os.Stdout, opts)),
	}
}

func (l slogLogger) Info(msg string) {
	l.logger.Info(msg)
}

func (l slogLogger) Error(msg string) {
	l.logger.Error(msg)
}
