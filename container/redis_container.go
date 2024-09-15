package container

import (
	"context"
	"customers_kuber/closer"
	"customers_kuber/logger"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"log/slog"
	"time"
)

func RunRedis() error {

	//конфигурация редиса, прокидывание портов
	ctx := context.Background()
	redisReq := testcontainers.ContainerRequest{
		Name:         "redis",
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
		return fmt.Errorf("failed to start redis container: %s", err)
	}

	//передача функции в closer для graceful shutdown
	closer.CloseFunctions = append(closer.CloseFunctions, func() {
		if err = redisContainer.Terminate(ctx); err != nil {
			ctx = logger.WithLogError(ctx, err)
			slog.ErrorContext(ctx, "failed to terminate redis container")
			return
		}
		slog.Info("redis container terminated successfully")
	})
	time.Sleep(time.Second * 3)
	return nil
}
