package container

import (
	"customers_kuber/config"
	"log"
)

func CreateContainers() {
	if config.Kuber == "" {
		if err := RunRedis(); err != nil {
			log.Fatalf("failed to start application: %s", err)
		}
		if err := RunPostgres(); err != nil {
			log.Fatalf("failed to start application: %s", err)
		}
		if err := RunKafka(); err != nil {
			log.Fatalf("failed to start application: %s", err)
		}
		if err := RunElastic(); err != nil {
			log.Fatalf("failed to start application: %s", err)
		}
		if err := RunLogstash(); err != nil {
			log.Fatalf("failed to start application: %s", err)
		}
		if err := RunKibana(); err != nil {
			log.Fatalf("failed to start application: %s", err)
		}
	}
}
