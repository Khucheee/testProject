package container

import (
	"customers_kuber/config"
	"log/slog"
	"os"
)

func CreateContainers() {
	if config.Kuber == "" {
		if err := RunRedis(); err != nil {
			slog.Error("failed to start application", "error", err)
			os.Exit(1)
		}
		if err := RunPostgres(); err != nil {
			slog.Error("failed to start application", "error", err)
			os.Exit(1)
		}
		if err := RunKafka(); err != nil {
			slog.Error("failed to start application", "error", err)
			os.Exit(1)
		}
		if err := RunElastic(); err != nil {
			slog.Error("failed to start application", "error", err)
			os.Exit(1)
		}
		if err := RunLogstash(); err != nil {
			slog.Error("failed to start application", "error", err)
			os.Exit(1)
		}
		if err := RunKibana(); err != nil {
			slog.Error("failed to start application", "error", err)
			os.Exit(1)
		}
	}
}
