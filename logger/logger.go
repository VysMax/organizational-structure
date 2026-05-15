package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/VysMax/organizational-structure/config"
)

func Init(cfg *config.Config) (*slog.Logger, error) {
	var writer io.Writer = os.Stdout
	if cfg.Logger.File != "" {
		dir := filepath.Dir(cfg.Logger.File)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to make directory for logs: %w", err)
		}

		file, err := os.OpenFile(cfg.Logger.File, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			return nil, fmt.Errorf("failed to open logs file: %w", err)
		}

		writer = file
	}

	logOpts := slog.NewJSONHandler(writer, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})
	logger := slog.New(logOpts)
	slog.SetDefault(logger)

	return logger, nil
}
