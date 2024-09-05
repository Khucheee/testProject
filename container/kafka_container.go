package container

import (
	"context"
	"customers_kuber/closer"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"time"
)

func RunKafka() error {

	//конфигурация кафки
	ctx := context.Background()
	kafkaReq := testcontainers.ContainerRequest{
		Name:         "kafka",
		Image:        "confluentinc/confluent-local:7.5.0",
		ExposedPorts: []string{"9092/tcp"},
		Env: map[string]string{
			"KAFKA_ADVERTISED_LISTENERS":           "PLAINTEXT://localhost:29092,PLAINTEXT_HOST://localhost:9092,PLAINTEXT_EXTERNAL://kafka:9095",
			"KAFKA_LISTENER_SECURITY_PROTOCOL_MAP": "PLAINTEXT:PLAINTEXT,PLAINTEXT_HOST:PLAINTEXT,PLAINTEXT_EXTERNAL:PLAINTEXT,CONTROLLER:PLAINTEXT",
			"KAFKA_LISTENERS":                      "PLAINTEXT://localhost:29092,CONTROLLER://localhost:29093,PLAINTEXT_HOST://0.0.0.0:9092,PLAINTEXT_EXTERNAL://kafka:9095"},
		HostConfigModifier: func(hostConfig *container.HostConfig) {
			hostConfig.NetworkMode = "NET"
			hostConfig.PortBindings = nat.PortMap{
				"9092/tcp": []nat.PortBinding{
					{HostIP: "0.0.0.0", HostPort: "9092"},
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

	//передача функции в closer для graceful shutdown
	closer.CloseFunctions = append(closer.CloseFunctions, func() {
		if err = kafkaContainer.Terminate(ctx); err != nil {
			log.Println("failed to terminate kafka container:", err)
			return
		}
		log.Println("kafka container terminated successfully")
	})
	time.Sleep(time.Second * 3)
	return nil
}
