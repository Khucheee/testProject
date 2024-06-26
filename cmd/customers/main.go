package main

import (
	"Customers/internal/app"
	"os/exec"
	"time"
)

func main() {
	RunKafka()
	controller := app.NewController()
	controller.Route()
}

func RunKafka() {
	cmd := exec.Command("docker", "compose", "up", "-d")
	cmd.Run()
	time.Sleep(time.Second * 5)
	app.CreateTopic("json_topic", 1, 1)

}
