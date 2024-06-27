package main

import (
	"Customers/controller"
	"Customers/producer"
	"context"
	"fmt"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"os/exec"
	"time"
)

func main() {
	//запускаем сервис
	controller.GetEntityController().Route()
}

func RunContainers() {
	cmd := exec.Command("docker", "compose", "up", "-d")
	cmd.Run()
	//тут лучше вызвать команду для получения статуса контейнеров, но пока тут просто статичесское ожидание
	time.Sleep(time.Second * 5)
	producer.CreateTopic("json_topic", 1, 1)

}

func RunContainer() {
	ctx := context.Background()
	_, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:12.19-alpine3.20"),
		testcontainers.WithHostPortAccess(5432),
		postgres.WithDatabase("mydatabase"),
		postgres.WithPassword("root"),
		postgres.WithUsername("admin"),
	)
	if err != nil {
		fmt.Println("Контекнер с базой не поднялся:", err)
	}
}
