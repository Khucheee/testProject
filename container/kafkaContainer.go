package container

import (
	"Customers/closer"
	"context"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"time"
)

func RunKafka() {
	ctx := context.Background()

	//конфигурация кафки
	kafkaReq := testcontainers.ContainerRequest{
		Image:        "confluentinc/confluent-local:7.5.0",
		ExposedPorts: []string{"9092/tcp", "9293/tcp"},
		HostConfigModifier: func(hostConfig *container.HostConfig) {
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
		log.Println("failed to start Kafka container:", err)
	}

	//передача функции в closer для graceful shutdown
	closer.CloseFunctions = append(closer.CloseFunctions, func() {
		if err = kafkaContainer.Terminate(ctx); err != nil {
			log.Println("failed to terminate postgres container:", err)
			return
		}
		log.Println("kafka container terminated successfully")
	})
	time.Sleep(time.Second * 3)
}
