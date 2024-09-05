package main

import (
	"customers_kuber/closer"
	"customers_kuber/config"
	"customers_kuber/container"
	"customers_kuber/controller"
	"log"
	"sync"
)

func main() {

	//установка конфига
	config.SetConfig()

	//проверка окружения, если запуск локальный, используем testcontainers
	if config.Kuber == "" {
		if err := container.RunRedis(); err != nil {
			log.Printf("failed to start application: %s", err)
			return
		}
		if err := container.RunPostgres(); err != nil {
			log.Printf("failed to start application: %s", err)
			return
		}
		if err := container.RunKafka(); err != nil {
			log.Printf("failed to start application: %s", err)
			return
		}
		if err := container.RunLogstash(); err != nil {
			log.Printf("failed to start application: %s", err)
			return
		}
	}

	//logger.InitLogging()
	//отлавливаем сигнал ctrl+c для graceful shutdown
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		closer.CtrlC()
		wg.Done()
	}()

	//запускаем сервис
	entityController, err := controller.GetEntityController()
	if err != nil {
		log.Printf("failed to start application: %s", err)
		return
	}

	//запускаем сервер
	entityController.Route()

	//ждем завершения graceful shutdown
	wg.Wait()
	log.Println("Выходим из программы")
}
