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
		container.RunPostgres()
		container.RunKafka()
	}

	//отлавливаем сигнал ctrl+c для graceful shutdown
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		closer.CtrlC()
		wg.Done()
	}()

	//запускаем сервис
	controller.GetEntityController().Route()

	//ждем завершения graceful shutdown
	wg.Wait()
	log.Println("Выходим из программы")
}
