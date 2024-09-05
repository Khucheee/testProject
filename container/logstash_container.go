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

func RunLogstash() error {

	ctx := context.Background()
	logstashReq := testcontainers.ContainerRequest{
		Name:         "logstash",
		Image:        "logstash:8.15.0",
		ExposedPorts: []string{"5044/tcp"},
		Mounts: testcontainers.Mounts(testcontainers.BindMount("C:\\Users\\AKhuchashev\\GolandProjects\\"+
			"Customers_kuber\\config", "/usr/share/logstash/pipeline")),
		Cmd: []string{"logstash", "-f", "/usr/share/logstash/pipeline/logstash.conf"},
		HostConfigModifier: func(hostConfig *container.HostConfig) {
			hostConfig.NetworkMode = "NET"
			hostConfig.PortBindings = nat.PortMap{
				"5044/tcp": []nat.PortBinding{
					{HostIP: "0.0.0.0", HostPort: "5044"},
				}}
		},
		//WaitingFor: wait.ForListeningPort("5044/tcp"),
	}

	//запуск контейнера
	logstashContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: logstashReq,
		Started:          true,
	})
	if err != nil {
		return fmt.Errorf("failed to start logstash container: %s", err)
	}

	//передача функции в closer для graceful shutdown
	closer.CloseFunctions = append(closer.CloseFunctions, func() {
		if err = logstashContainer.Terminate(ctx); err != nil {
			log.Println("failed to terminate logstash container:", err)
			return
		}
		log.Println("logstash container terminated successfully")
	})
	time.Sleep(time.Second * 3)
	return nil
}
