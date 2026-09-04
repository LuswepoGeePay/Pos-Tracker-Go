package utils

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
)

var logger *slog.Logger

func InitLogger(logFilePath string) error {
	level := slog.LevelInfo
	switch strings.ToLower(strings.TrimSpace(os.Getenv("LOG_LEVEL"))) {
	case "debug":
		level = slog.LevelDebug
	case "warn", "warning":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}

	handlerOptions := &slog.HandlerOptions{
		Level: level,
	}

	writers := []io.Writer{os.Stdout}

	if logFilePath == "" {
		logFilePath = os.Getenv("LOG_FILE")
	}
	if logFilePath == "" {
		logFilePath = "pos_master.log"
	}

	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, handlerOptions))
		slog.SetDefault(logger)
		return fmt.Errorf("failed to open log file %s: %w (logging to stdout only)", logFilePath, err)
	}

	w := io.MultiWriter(writers[0], file)
	logger = slog.New(slog.NewJSONHandler(w, handlerOptions))
	slog.SetDefault(logger)

	return nil
}

func Log(level slog.Level, msg string, args ...interface{}) {
	ctx := context.Background()

	if logger != nil {
		logger.Log(ctx, level, msg, args...)
		return
	}
	slog.Log(ctx, level, msg, args...)
}

func Info(msg string, args ...any) {
	Log(slog.LevelInfo, msg, args...)
}

func Warn(msg string, args ...any) {
	Log(slog.LevelWarn, msg, args...)
}

func Error(msg string, args ...any) {
	Log(slog.LevelError, msg, args...)
}
