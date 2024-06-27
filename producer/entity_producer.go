package producer

import (
	"Customers/closer"
	"Customers/model"
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	kafka "github.com/segmentio/kafka-go"
	"log"
	"time"
)

var entityProducerInstance *entityProducer

type EntityProducer interface {
	ProduceEntityToKafka(entity model.Entity)
	CloseEntityProducer() func()
}

type entityProducer struct {
	writer *kafka.Writer
}

func GetEntityProducer() EntityProducer {
	if entityProducerInstance != nil {
		return entityProducerInstance
	}

	//создаем продюсера для отправки сообщений в kafka
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "json_topic",
	})

	//инициализируем инстанс продюсера, передаем функцию клозеру для graceful shutdown
	entityProducerInstance = &entityProducer{writer}
	closer.CloseFunctions = append(closer.CloseFunctions, entityProducerInstance.CloseEntityProducer())
	return entityProducerInstance
}

func (producer *entityProducer) ProduceEntityToKafka(entity model.Entity) {

	//парсим полученную структуру в json
	message, err := json.Marshal(entity)
	if err != nil {
		log.Println("failed to marshal JSON while produce into kafka:", err)
	}

	//отправляем сообщение в топик
	err = producer.writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(entity.Id),
			Value: message,
		},
	)
	if err != nil {
		log.Println("failed to produce message into kafka:", err)
	}
}

func (producer *entityProducer) CloseEntityProducer() func() {
	return func() {
		if err := producer.writer.Close(); err != nil {
			log.Println("producer closing failed:", err)
			return
		}
		log.Println("entityProduccer closed successfully")
	}
}

func CreateTopic(topic string, numPartitions int32, replicationFactor int16) error {

	//создаем админа кластера
	time.Sleep(time.Second * 10)
	admin, err := sarama.NewClusterAdmin([]string{"localhost:9092"}, nil)
	if err != nil {
		log.Println("failed to create cluster admin to create topic in kafka:", err)
	}

	//после создания топика закрываем админа
	defer func() {
		if err := admin.Close(); err != nil {
			log.Println("Failed to close admin object after creating topic in kafka:", err)
		}
	}()

	//Собираем структуру, чтобы скормить её в функцию создания топика
	topicDetail := sarama.TopicDetail{
		NumPartitions:     numPartitions,
		ReplicationFactor: replicationFactor,
	}

	//создаем топик

	err = admin.CreateTopic(topic, &topicDetail, false)
	if err != nil && err != sarama.ErrTopicAlreadyExists {
		log.Println("Failed to create topic in kafka:", err)
		return err
	}
	return nil
}
