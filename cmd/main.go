package main

import (
	"Customers/closer"
	"Customers/controller"
	"context"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"log"
	"os/exec"
	"time"
)

func main() {
	//запускаем сервис
	RunPostgres()
	RunKafka()
	//producer.CreateTopic("json_topic", 1, 1)
	go closer.CtrlC()
	controller.GetEntityController().Route()

}

func RunContainers() {
	cmd := exec.Command("docker", "compose", "up", "-d")
	cmd.Run()
	//тут лучше вызвать команду для получения статуса контейнеров, но пока тут просто статичесское ожидание
	time.Sleep(time.Second * 5)

}

func RunPostgres() {
	ctx := context.Background()
	postgreSQLReq := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_DB":       "mydatabase",
			"POSTGRES_USER":     "admin",
			"POSTGRES_PASSWORD": "root",
		},
		HostConfigModifier: func(hostConfig *container.HostConfig) {
			hostConfig.PortBindings = nat.PortMap{
				"5432/tcp": []nat.PortBinding{
					{HostIP: "0.0.0.0", HostPort: "5432"},
				}}
		}}

	pgC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: postgreSQLReq,
		Started:          true,
	})
	if err != nil {
		log.Fatalf("Failed to start PostgreSQL container: %s", err)
	}

	pgPort, err := pgC.MappedPort(ctx, "5432")
	if err != nil {
		log.Fatalf("Failed to get PostgreSQL mapped port: %s", err, pgPort)
	}

}

func RunKafka() {

}
