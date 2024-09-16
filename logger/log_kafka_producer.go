package logger

import (
	"context"
	"customers_kuber/closer"
	"customers_kuber/config"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log/slog"
	"sync"
)

var logProducerInstance *logProducer
var mu sync.Mutex

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
	ProduceLogToKafka(log []byte)
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
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{kafkaAddress}, ///здесь нужно помнить про адрес
		Topic:   config.KafkaLogTopic,
	})

	//инициализируем инстанс продюсера, передаем функцию клозеру для graceful shutdown
	logProducerInstance = &logProducer{writer}
	closer.CloseFunctions = append(closer.CloseFunctions, logProducerInstance.CloseLogProducer())
	return logProducerInstance, nil
}

func (producer *logProducer) ProduceLogToKafka(logs []byte) {
	mu.Lock()
	defer mu.Unlock()
	err := producer.writer.WriteMessages(context.Background(), kafka.Message{Value: logs})
	if err != nil && config.KafkaEnabled {
		fmt.Printf("\nFAILED TO PRODUCE LOG TO KAFKA: %s\n", err)
	}
}

func (producer *logProducer) CloseLogProducer() func() {
	return func() {
		if err := producer.writer.Close(); err != nil {
			slog.ErrorContext(WithLogError(context.Background(), err), "log producer closing failed")
			return
		}
		slog.Info("logProducer closed successfully")
	}
}
