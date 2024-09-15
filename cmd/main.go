package main

import (
	"context"
	"customers_kuber/closer"
	"customers_kuber/config"
	"customers_kuber/container"
	"customers_kuber/controller"
	"customers_kuber/logger"
)

func main() {

	wg := closer.InitGracefulShutdown()

	config.SetConfig()

	container.CreateContainers()

	logger.InitLogging()

	controller.GetEntityController().Route(context.Background())

	wg.Wait()

}
