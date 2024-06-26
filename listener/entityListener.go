package listener

import (
	"Customers/model"
	"Customers/repository"
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

func StartListeningKafka() {
	fmt.Println("Запущено прослушивание топика в кафке")
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    "json_topic",
		GroupID:  "json_consumer_group",
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	// Чтение сообщений из Kafka
	for {
		time.Sleep(time.Second * 1)
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Fatalf("failed to read message: %v", err)
		}
		var entity model.Entity
		err = json.Unmarshal(msg.Value, &entity)
		if err != nil {
			log.Printf("failed to unmarshal message: %v", err)
			continue
		}
		fmt.Println("Message from kafka recieved:", entity)
		repository.Db.SaveCustomer(entity)
	}
	fmt.Println("Прослушивание отключено")
}
