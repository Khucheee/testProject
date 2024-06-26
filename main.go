package main

import (
	"Customers/controller"
	"Customers/listener"
	"Customers/producer"
	"Customers/repository"
	"os/exec"
	"time"
)

func main() {
	RunContainers()
	controller := controller.NewController()
	repository.NewPostgres()
	go listener.StartListeningKafka()
	controller.Route()

}

func RunContainers() {
	cmd := exec.Command("docker", "compose", "up", "-d")
	cmd.Run()
	//тут лучше вызвать команду для получения статуса контейнеров, но пока тут просто статичесское ожидание
	time.Sleep(time.Second * 5)
	producer.CreateTopic("json_topic", 1, 1)

}
