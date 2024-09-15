package listener

import (
	"context"
	"customers_kuber/cache"
	"customers_kuber/closer"
	"customers_kuber/config"
	"customers_kuber/logger"
	"customers_kuber/model"
	"customers_kuber/repository"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"log/slog"
	"time"
)

var entityListenerInstance *entityListener

type EntityListener interface {
	StartListening()
}

type entityListener struct {
	reader     *kafka.Reader
	repository repository.EntityRepository
	stopSignal bool
	cache      cache.EntityCache
}

func GetEntityListener() (EntityListener, error) {

	//если сущность уже существует, возвращаем её
	if entityListenerInstance != nil {
		return entityListenerInstance, nil
	}

	ctx := context.Background()

	//определяю адрес кафки
	kafkaAddress := fmt.Sprintf("%s:%s", config.KafkaHost, config.KafkaPort)

	//Создаю объект для чтения сообщений из kafka
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafkaAddress}, // здесь адрес подключения к кафке
		Topic:   config.KafkaTopic,
	})

	cacheConnect, err := cache.GetEntityCache()
	if err != nil {
		log.Printf("failed to get cache in listener: %s", err)
	}

	//инициализируем инсстанс листенера
	entityListenerInstance := &entityListener{reader: reader, stopSignal: false, cache: cacheConnect}
	entityRepository, err := repository.GetEntityRepository()
	entityListenerInstance.repository = entityRepository
	if err != nil {
		return entityListenerInstance, fmt.Errorf("failed to initislize repository in listener: %s", err)
	}

	//передаем функцию закрытия в клозер для graceful shutdown
	closer.CloseFunctions = append(closer.CloseFunctions, func() {
		entityListenerInstance.stopSignal = true
		if err := entityListenerInstance.reader.Close(); err != nil {
			ctx = logger.WithLogError(ctx, err)
			slog.ErrorContext(ctx, "failed to close listener")
			return
		}
		log.Println("entityListener closed successfully")
	})
	return entityListenerInstance, nil
}

func (listener *entityListener) StartListening() {

	//вызываем раз в секунду ReadMessage чтобы забрать сообщение из топика
	for {
		if listener.stopSignal == true {
			break
		}
		time.Sleep(time.Second * 1)
		ctx := context.Background()
		msg, err := listener.reader.ReadMessage(ctx)
		if err != nil {
			if err.Error() == "fetching message: EOF" {
				continue
			}
			log.Println("failed to read message:", err)
			continue
		}

		//парсю полученную json в структуру Entity
		var entity model.Entity
		err = json.Unmarshal(msg.Value, &entity)
		if err != nil {
			log.Println("failed to deserialize message from kafka:", err)
			continue
		}
		slog.InfoContext(logger.WithLogValues(ctx, entity), "listener received message from kafka")

		//сохраняю Entity в базу
		for i := 1; i < 3; i++ {
			if err = listener.repository.SaveEntity(entity); err == nil {
				break
			}
			log.Printf("Failed to save entity in service: %s", err)
			log.Printf("Retry to save entity: %d", i)
		}
		listener.cache.ClearCache(ctx)
	}
}
