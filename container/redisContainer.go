package container

import (
	"context"
	"customers_kuber/closer"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"log"
	"time"
)

func RunRedis() {
	ctx := context.Background()
	redisReq := testcontainers.ContainerRequest{
		Image:        "redis:latest",
		ExposedPorts: []string{"6379/tcp", "6379/tcp"},
		HostConfigModifier: func(hostConfig *container.HostConfig) {
			hostConfig.PortBindings = nat.PortMap{
				"6379/tcp": []nat.PortBinding{
					{HostIP: "0.0.0.0", HostPort: "6379"},
				}}
		}}

	//запуск контейнера
	redisContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: redisReq,
		Started:          true,
	})
	if err != nil {
		log.Println("failed to start redis container:", err)
	}

	//передача функции в closer для graceful shutdown
	closer.CloseFunctions = append(closer.CloseFunctions, func() {
		if err = redisContainer.Terminate(ctx); err != nil {
			log.Println("failed to terminate redis container:", err)
			return
		}
		log.Println("redis container terminated successfully")
	})
	time.Sleep(time.Second * 3)
}
