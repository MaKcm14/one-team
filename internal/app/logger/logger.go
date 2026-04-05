package logger

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/MaKcm14/one-team/internal/config"
)

type Logger struct {
	*slog.Logger

	logFile *os.File
}

func New(cfg config.LoggerConfig) (Logger, error) {
	file, err := os.Create(configFileName())
	if err != nil {
		return Logger{}, fmt.Errorf("%w: %s", ErrLoggerConfig, err)
	}

	lvl, err := getLoggerLvl(cfg.Mode)
	if err != nil {
		return Logger{}, err
	}

	return Logger{
		Logger: slog.New(slog.NewTextHandler(file, &slog.HandlerOptions{
			Level: lvl,
		})),
		logFile: file,
	}, nil
}

func configFileName() string {
	hour, min, seconds := time.Now().UTC().Clock()
	year, month, day := time.Now().UTC().Date()
	return fmt.Sprintf("./logs/one-team___%d-%d-%d_%d-%d-%d_UTC___.log", year, month, day, hour, min, seconds)
}

func (l Logger) Instance() *slog.Logger {
	return l.Logger
}

func (l Logger) Close() {
	l.logFile.Close()
}

func getLoggerLvl(lvlName string) (slog.Level, error) {
	switch strings.ToLower(lvlName) {
	case "debug":
		return slog.LevelDebug, nil

	case "info":
		return slog.LevelInfo, nil

	case "warn":
		return slog.LevelWarn, nil

	case "error":
		return slog.LevelError, nil
	}
	return 0, fmt.Errorf("unknown logger level")
}
