package container

import (
	"context"
	"customers_kuber/closer"
	"customers_kuber/config"
	"customers_kuber/logger"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"log/slog"
	"time"
)

func RunKafka() error {

	//конфигурация кафки
	ctx := context.Background()
	kafkaReq := testcontainers.ContainerRequest{
		Name:         "kafka",
		Image:        "confluentinc/confluent-local:7.5.0",
		ExposedPorts: []string{config.KafkaPort + "/tcp"},
		Env: map[string]string{
			"KAFKA_ADVERTISED_LISTENERS":           "PLAINTEXT://localhost:29092,PLAINTEXT_HOST://localhost:9092,PLAINTEXT_EXTERNAL://kafka:9095",
			"KAFKA_LISTENER_SECURITY_PROTOCOL_MAP": "PLAINTEXT:PLAINTEXT,PLAINTEXT_HOST:PLAINTEXT,PLAINTEXT_EXTERNAL:PLAINTEXT,CONTROLLER:PLAINTEXT",
			"KAFKA_LISTENERS":                      "PLAINTEXT://localhost:29092,CONTROLLER://localhost:29093,PLAINTEXT_HOST://0.0.0.0:9092,PLAINTEXT_EXTERNAL://kafka:9095"},
		HostConfigModifier: func(hostConfig *container.HostConfig) {
			hostConfig.NetworkMode = "NET"
			hostConfig.PortBindings = nat.PortMap{
				"9092/tcp": []nat.PortBinding{
					{HostIP: "0.0.0.0", HostPort: config.KafkaPort},
				}}
		},

		WaitingFor: wait.ForListeningPort("9092/tcp"),
	}

	//запуск контейнера
	kafkaContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: kafkaReq,
		Started:          true,
	})
	if err != nil {
		return fmt.Errorf("failed to start kafka container: %s", err)
	}
	config.KafkaEnabled = true
	//передача функции в closer для graceful shutdown
	closer.CloseFunctions = append(closer.CloseFunctions, func() {
		if err = kafkaContainer.Terminate(ctx); err != nil {
			ctx = logger.WithLogError(ctx, err)
			slog.ErrorContext(ctx, "failed to terminate kafka container")
			return
		}
		config.KafkaEnabled = false
		slog.Info("kafka container terminated successfully")
	})
	time.Sleep(time.Second * 3)
	return nil
}
