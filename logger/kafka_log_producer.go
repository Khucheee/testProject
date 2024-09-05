package logger

import (
	"context"
	"customers_kuber/closer"
	"customers_kuber/config"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
)

var logProducerInstance *logProducer

type defaultLog struct {
	Log string
}

func CreateLogTopic() error {
	kafkaAddress := fmt.Sprintf("%s:%s", config.KafkaHost, config.KafkaPort)

	kafkaConnect, err := kafka.Dial("tcp", kafkaAddress)
	if err != nil {
		return fmt.Errorf("failed to create kafka connection while creating topic: %s", err)
	}
	topicConfig := kafka.TopicConfig{Topic: config.KafkaLogTopic, NumPartitions: 1, ReplicationFactor: 1}
	if err = kafkaConnect.CreateTopics(topicConfig); err != nil {
		return fmt.Errorf("something going wrong while creating kafka log topic: %s", err)
	}
	return nil
}

type LogProducer interface {
	ProduceLogToKafka(log string)
	CloseLogProducer() func()
}

type logProducer struct {
	writer *kafka.Writer
}

func GetLogProducer() (LogProducer, error) {

	//если сущность уже есть, то возвращаем её
	if logProducerInstance != nil {
		return logProducerInstance, nil
	}
	//определяю адрес кафки
	kafkaAddress := fmt.Sprintf("%s:%s", config.KafkaHost, config.KafkaPort)

	//создаем продюсера для отправки сообщений в kafka
	//надо подумать как проверить работоспособность райтера

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{kafkaAddress}, ///здесь нужно помнить про адрес
		Topic:   config.KafkaLogTopic,
	})

	//инициализируем инстанс продюсера, передаем функцию клозеру для graceful shutdown
	logProducerInstance = &logProducer{writer}
	closer.CloseFunctions = append(closer.CloseFunctions, logProducerInstance.CloseLogProducer())
	return logProducerInstance, nil
}

func (producer *logProducer) ProduceLogToKafka(logs string) {

	message, err := json.Marshal(defaultLog{logs})
	if err != nil {
		log.Printf("failed to marshal JSON while produce into kafka: %s", err)
		return
	}
	_ = producer.writer.WriteMessages(context.Background(), kafka.Message{Value: message})

}

func (producer *logProducer) CloseLogProducer() func() {
	return func() {
		if err := producer.writer.Close(); err != nil {
			log.Println("log producer closing failed:", err)
			return
		}
		log.Println("logProducer closed successfully")
	}
}
