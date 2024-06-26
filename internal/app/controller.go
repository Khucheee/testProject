package app

import (
	storage2 "Customers/internal/storage"
	"fmt"
)

var controller *Controller

type Controller struct {
	producer *Producer
	consumer *Consumer
	storage  *storage2.Storage
}

func NewController() *Controller {
	if controller == nil {
		p := NewProducer()
		c := NewConsumer()
		s := storage2.NewStorage()
		//s.cash = c.db.Restore()
		controller = &Controller{p, c, s}
		return controller
	}
	fmt.Println("Контроллер уже существует")
	return controller
}
