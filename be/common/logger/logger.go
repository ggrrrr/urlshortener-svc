package logger

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/caarlos0/env/v11"
)

type config struct {
	Level  string `env:"LOG_LEVEL"`
	Format string `env:"LOG_FORMAT" envDefault:"text"`
}

func Configure() {
	cfg := config{}

	err := env.Parse(&cfg)
	if err != nil {
		fmt.Println(err)
		cfg.Level = "DEBUG"
		cfg.Format = "TEXT"
	}

	opts := &slog.HandlerOptions{
		Level: parseLogLevel(cfg.Level),
	}

	var logger *slog.Logger

	switch strings.ToUpper(cfg.Format) {
	case "JSON":
		logger = configureJSON(opts)
	default:
		logger = configureText(opts)
	}

	slog.SetDefault(logger)
}

func configureText(opts *slog.HandlerOptions) *slog.Logger {
	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))
	return logger
}

func configureJSON(opts *slog.HandlerOptions) *slog.Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
	return logger
}

func parseLogLevel(from string) slog.Level {
	switch strings.ToUpper(from) {
	case "INFO":
		return slog.LevelInfo
	default:
		return slog.LevelDebug
	}
}
