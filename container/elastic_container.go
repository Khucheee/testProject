package container

import (
	"context"
	"customers_kuber/closer"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"log"
	"time"
)

func RunElastic() error {

	ctx := context.Background()
	elasticReq := testcontainers.ContainerRequest{
		Name:         "elasticsearch",
		Image:        "elasticsearch:8.15.0",
		ExposedPorts: []string{"9200/tcp"},
		Env: map[string]string{
			"discovery.type":         "single-node",
			"xpack.security.enabled": "false", // Отключаем xpack security для упрощения тестов
		},
		HostConfigModifier: func(hostConfig *container.HostConfig) {
			hostConfig.NetworkMode = "NET"
			hostConfig.PortBindings = nat.PortMap{
				"9200/tcp": []nat.PortBinding{
					{HostIP: "0.0.0.0", HostPort: "9200"},
				}}
		},
	}

	//запуск контейнера
	elasticContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: elasticReq,
		Started:          true,
	})
	if err != nil {
		return fmt.Errorf("failed to start elastic container: %s", err)
	}

	//передача функции в closer для graceful shutdown
	closer.CloseFunctions = append(closer.CloseFunctions, func() {
		if err = elasticContainer.Terminate(ctx); err != nil {
			log.Println("failed to terminate elastic container:", err)
			return
		}
		log.Println("elastic container terminated successfully")
	})
	time.Sleep(time.Second * 3)
	return nil
}
