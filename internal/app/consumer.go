package app

import (
	"Customers/internal/PostgreDB"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"log"
)

var consumer *Consumer

type Consumer struct {
	c       sarama.Consumer
	topic   string
	brokers []string
	db      *PostgreDB.Postgres
}

func NewConsumer() *Consumer {
	if consumer != nil {
		return consumer
	}
	//создаем консьюмера и подключаем к нему базу
	consumer = &Consumer{topic: "json_topic", brokers: []string{"localhost:9092"}, db: PostgreDB.NewPostgres()}
	config := sarama.NewConfig()
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Return.Errors = true
	c, err := sarama.NewConsumer(consumer.brokers, config)
	if err != nil {
		fmt.Println("Не получилось создать консьюмера:", err)
	}
	consumer.c = c
	return consumer
}
func (c *Consumer) Read() {
	//начинаем слушать топик
	partitionConsumer, err := c.c.ConsumePartition(c.topic, 0, sarama.OffsetOldest)
	if err != nil {
		fmt.Println("Не получилось создать консьюмера:", err)
	}
	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			fmt.Println("Не получилось перестать слушать топик", err)
		}
	}()
	var message Entity
	// Читаем сообщение, пока только одно
	msg := <-partitionConsumer.Messages()
	err = json.Unmarshal(msg.Value, &message)
	if err != nil {
		log.Println("Не получилось распарсить JSON из кафки", err)
	} else {
		fmt.Println("Получено сообщение", message)
	}
	//сохраняем в базу
	c.SaveCustomer(message)
}
func (c *Consumer) SaveCustomer(e Entity) {
	fmt.Println("Обращаемся к базе из консьюмера")
	c.db.SaveCustomer(e)
}

//Теперь пишем в базу

//return message
/*select {
case msg := <-partitionConsumer.Messages():
	fmt.Println("вот тут мы получили сообщение из кафки")
	err := json.Unmarshal(msg.Value, &message)
	if err != nil {
		log.Println("Failed to unmarshal message:", err)
	} else {
		log.Printf("Received message: %+v\n", message)
	}
	return message
case err := <-partitionConsumer.Errors():
	fmt.Println("Error while consuming partition:", err)
}
return message*/
