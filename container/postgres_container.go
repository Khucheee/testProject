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

func RunPostgres() error {

	//конфигурация базы данных
	ctx := context.Background()
	postgreSQLReq := testcontainers.ContainerRequest{
		Name:         "postgres",
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp", "5435/tcp"},
		Env: map[string]string{
			"POSTGRES_DB":       "postgres",
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "root",
		},

		HostConfigModifier: func(hostConfig *container.HostConfig) {
			hostConfig.PortBindings = nat.PortMap{
				"5432/tcp": []nat.PortBinding{
					{HostIP: "0.0.0.0", HostPort: "5432"},
				}}
		}}

	//запуск контейнера
	pgContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: postgreSQLReq,
		Started:          true,
	})
	if err != nil {
		return fmt.Errorf("failed to start postgres container: %s", err)
	}

	//передача функции в closer для graceful shutdown
	closer.CloseFunctions = append(closer.CloseFunctions, func() {
		if err = pgContainer.Terminate(ctx); err != nil {
			ctx = logger.WithLogError(ctx, err)
			slog.ErrorContext(ctx, "failed to terminate postgres container")
			return
		}
		slog.Info("postgres container terminated successfully")
	})
	time.Sleep(time.Second * 3)
	return nil
}
