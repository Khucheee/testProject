package container

import (
	"customers_kuber/config"
	"log"
)

//На счет логирования
//Филосовия slog говорит, что логирование не должно завершать программу
//Поэтому в slog нет функции, аналогичной Fatal(), пока решил оставить так

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
