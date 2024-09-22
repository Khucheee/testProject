package producer

import (
	"context"
	"customers_kuber/closer"
	"customers_kuber/config"
	"customers_kuber/logger"
	"customers_kuber/model"
	"encoding/json"
	"fmt"
	kafka "github.com/segmentio/kafka-go"
	"log"
	"log/slog"
)

var entityProducerInstance *entityProducer

type EntityProducer interface {
	ProduceEntityToKafka(entity model.Entity) error
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
	closer.CloseFunctions = append(closer.CloseFunctions, func() {
		if err := writer.Close(); err != nil {
			slog.ErrorContext(logger.WithLogError(context.Background(), err), "failed to close entityProducer")
			log.Println("producer closing failed:", err)
			return
		}
		slog.Info("entityProducer closed successfully")
	})
	return entityProducerInstance, nil
}

func (producer *entityProducer) ProduceEntityToKafka(entity model.Entity) error {

	//парсим полученную структуру в json
	message, err := json.Marshal(entity)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON while produce into kafka: %s", err)
	}

	//отправляем сообщение в топик
	idForKafka, err := entity.Id.MarshalBinary()
	if err != nil {
		return fmt.Errorf("failed to convert uuid into bytes in producer: %s", err)
	}
	kafkaMessage := kafka.Message{Key: idForKafka, Value: message}

	if err = producer.writer.WriteMessages(context.Background(), kafkaMessage); err != nil {
		return fmt.Errorf("failed to produce message into kafka: %s", err)
	}
	slog.Info("entity successfully produced to kafka")
	return nil
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
