package logger

import (
	"fmt"
	"log/slog"
	"os"
	"time"
)

type Logger struct {
	*slog.Logger

	logFile *os.File
}

func New() (Logger, error) {
	file, err := os.Create(configFileName())
	if err != nil {
		return Logger{}, fmt.Errorf("%w: %s", ErrLoggerConfig, err)
	}

	return Logger{
		Logger: slog.New(slog.NewTextHandler(file, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})),
		logFile: file,
	}, nil
}

func configFileName() string {
	hour, min, seconds := time.Now().Clock()
	year, month, day := time.Now().Date()
	return fmt.Sprintf("one-team-{%d-%d-%d_%d:%d:%d_UTC}.log", year, month, day, hour, min, seconds)
}

func (l Logger) Instance() *slog.Logger {
	return l.Logger
}

func (l Logger) Close() {
	l.logFile.Close()
}
