package producer

import (
	"context"
	"customers_kuber/closer"
	"customers_kuber/config"
	"customers_kuber/model"
	"encoding/json"
	"fmt"
	kafka "github.com/segmentio/kafka-go"
	"log"
)

var entityProducerInstance *entityProducer

type EntityProducer interface {
	ProduceEntityToKafka(entity model.Entity) error
	CloseEntityProducer() func()
}

type entityProducer struct {
	writer *kafka.Writer
}

func GetEntityProducer() (EntityProducer, error) {

	//если сущность уже есть, то возвращаем её
	if entityProducerInstance != nil {
		return entityProducerInstance, nil
	}

	//определяю адрес кафки
	kafkaAddress := fmt.Sprintf("%s:%s", config.KafkaHost, config.KafkaPort)

	//создаем продюсера для отправки сообщений в kafka
	//надо подумать как проверить работоспособность райтера
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{kafkaAddress}, ///здесь нужно помнить про адрес
		Topic:   config.KafkaTopic,
	})

	//инициализируем инстанс продюсера, передаем функцию клозеру для graceful shutdown
	entityProducerInstance = &entityProducer{writer}
	closer.CloseFunctions = append(closer.CloseFunctions, entityProducerInstance.CloseEntityProducer())
	return entityProducerInstance, nil
}

func (producer *entityProducer) ProduceEntityToKafka(entity model.Entity) error {

	//парсим полученную структуру в json
	message, err := json.Marshal(entity)
	if err != nil {
		log.Printf("failed to marshal JSON while produce into kafka: %s", err)
		return err
	}

	//отправляем сообщение в топик
	idForKafka, err := entity.Id.MarshalBinary()
	if err != nil {
		log.Printf("failed to convert uuid into bytes in producer: %s", err)
		return err
	}

	err = producer.writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   idForKafka,
			Value: message,
		},
	)
	if err != nil {
		log.Printf("failed to produce message into kafka: %s", err)
		return err
	}
	return nil

}

func (producer *entityProducer) CloseEntityProducer() func() {
	return func() {
		if err := producer.writer.Close(); err != nil {
			log.Println("producer closing failed:", err)
			return
		}
		log.Println("entityProducer closed successfully")
	}
}

func CreateTopic() error {
	kafkaAddress := fmt.Sprintf("%s:%s", config.KafkaHost, config.KafkaPort)

	kafkaConnect, err := kafka.Dial("tcp", kafkaAddress)
	if err != nil {
		return fmt.Errorf("failed to create kafka connection while creating topic: %s", err)
	}
	topicConfig := kafka.TopicConfig{Topic: config.KafkaTopic, NumPartitions: 1, ReplicationFactor: 1}
	if err = kafkaConnect.CreateTopics(topicConfig); err != nil {
		return fmt.Errorf("something going wrong while creating kafka topic: %s", err)
	}
	return nil
}
