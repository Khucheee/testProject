package main

import (
	"Customers/closer"
	"Customers/container"
	"Customers/controller"
)

func main() {
	//запускаем сервис
	container.RunPostgres()
	container.RunKafka()
	go closer.CtrlC()
	controller.GetEntityController().Route()
}
