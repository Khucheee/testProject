package producer

import (
	"Customers/model"
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	kafka "github.com/segmentio/kafka-go"
	"log"
)

func ProduceToKafka(entity model.Entity) {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "json_topic",
	})

	defer func() {
		if err := writer.Close(); err != nil {
			log.Fatalf("failed to close writer: %v", err)
		}
	}()
	message, err := json.Marshal(entity)
	if err != nil {
		log.Fatalf("failed to marshal entity: %v", err)
	}
	// Отправляем сообщение в топик
	err = writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(entity.Id),
			Value: message,
		},
	)
	if err != nil {
		log.Fatalf("failed to write message: %v", err)
	}
	log.Printf("sent message: %+v\n", entity)
}

func CreateTopic(topic string, numPartitions int32, replicationFactor int16) error {
	//создаем админа кластера
	admin, err := sarama.NewClusterAdmin([]string{"localhost:9092"}, nil)
	if err != nil {
		fmt.Println("Не получилось создать админа кластера для создания топика", err)
	}
	defer func() {
		if err := admin.Close(); err != nil {
			fmt.Println("Не получилось закрыть админку для создания топика", err)
		}
	}()
	//Собираем структуру, чтобы скормить её в функцию создания топика
	topicDetail := sarama.TopicDetail{
		NumPartitions:     numPartitions,
		ReplicationFactor: replicationFactor,
	}
	//Создаем топик
	err = admin.CreateTopic(topic, &topicDetail, false)
	if err != nil && err != sarama.ErrTopicAlreadyExists {
		fmt.Println("Не получилось создать топик в кафке", err)
		return err
	}
	return nil
}
