package app

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"log"
	"time"
)

var producer *Producer

type Producer struct {
	p       sarama.SyncProducer
	topic   string
	brokers []string
}

func NewProducer() *Producer {
	if producer != nil {
		return producer
	}
	//создаю продюсера
	producer = &Producer{topic: "json_topic", brokers: []string{"localhost:9092"}}
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Consumer.Return.Errors = true
	config.Producer.Timeout = 5 * time.Second
	p, err := sarama.NewSyncProducer(producer.brokers, config)
	if err != nil {
		fmt.Println("Не получилось создать продюсера ", err)
	}
	producer.p = p
	//defer producer.Close()
	return producer
}

func (p *Producer) SendMessage(message Entity) {
	//собираю json
	msg, err := json.Marshal(message)
	if err != nil {
		log.Println("Не распарсился JSON в продюсере", err)
		return
	}
	//собираю сообщение в кафку
	producerMsg := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.ByteEncoder(msg),
	}
	//закидываю сообщение в кафку
	partition, offset, err := p.p.SendMessage(producerMsg)
	if err != nil {
		fmt.Println("Failed to send message:", err)
		return
	}
	fmt.Println("Message sent to partition %d at offset\n", partition, offset)
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
