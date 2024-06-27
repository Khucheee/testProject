package listener

import (
	"Customers/model"
	"Customers/repository"
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

var entityListenerInstance *entityListener

type EntityListener interface {
	StartListening()
}

type entityListener struct {
	reader     *kafka.Reader
	repository repository.EntityRepository
}

func GetEntityListener() EntityListener {
	if entityListenerInstance != nil {
		return entityListenerInstance
	}

	//Создаю объект для чтения сообщений из kafka
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    "json_topic",
		GroupID:  "json_consumer_group",
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
	entityRepository := repository.GetEntityRepository()
	return &entityListener{reader, entityRepository}

}

func (listener *entityListener) StartListening() {

	//вызываем раз в секунду ReadMessage чтобы забрать сообщение из топика
	for {
		time.Sleep(time.Second * 1)
		msg, err := listener.reader.ReadMessage(context.Background())
		if err != nil {
			log.Println("failed to read message:", err)
		}

		//парсю полученную json в структуру Entity
		var entity model.Entity
		err = json.Unmarshal(msg.Value, &entity)
		if err != nil {
			log.Println("failed to deserialize message from kafka:", err)
			continue
		}
		log.Println("message from kafka received:", entity)

		//сохраняю Entity в базу
		listener.repository.SaveEntity(entity)
	}
}
