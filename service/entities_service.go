package service

import (
	"context"
	"customers_kuber/cache"
	"customers_kuber/listener"
	"customers_kuber/model"
	"customers_kuber/producer"
	"customers_kuber/repository"
	"fmt"
	"github.com/google/uuid"
	"log"
	"log/slog"
)

var entityServiceInstance *entityService

type EntityService interface {
	SaveEntity(context.Context, model.Test) error
	GetAllEntities(context.Context, string) ([]model.Entity, error)
	UpdateEntity(context.Context, model.Entity) error
	DeleteEntity(context.Context, uuid.UUID) error
}

type entityService struct {
	repository repository.EntityRepository
	producer   producer.EntityProducer
	cache      cache.EntityCache
}

func GetEntityService() (EntityService, error) {

	if entityServiceInstance != nil {
		return entityServiceInstance, nil
	}

	if err := producer.CreateTopic(); err != nil {
		return entityServiceInstance, fmt.Errorf("failed to create service: %s", err)
	}

	//создаем листенера
	entityListener, err := listener.GetEntityListener()
	if err != nil {
		return entityServiceInstance, fmt.Errorf("failed to create service: %s", err)
	}

	//запускаем прослушивание кафки
	go entityListener.StartListening()

	//получаем коннект к базе
	entityRepository, err := repository.GetEntityRepository()
	if err != nil {
		return entityServiceInstance, fmt.Errorf("failed to create service: %s", err)
	}

	//полчаем продюсера
	entityProducer, err := producer.GetEntityProducer()
	if err != nil {
		return entityServiceInstance, fmt.Errorf("failed to create service: %s", err)
	}

	//получаем доступ к кэшу
	entityCache, err := cache.GetEntityCache()
	if err != nil {
		slog.Warn("failed to get cache in service")
	}

	entityServiceInstance = &entityService{entityRepository, entityProducer, entityCache}
	return entityServiceInstance, nil
}

func (service *entityService) SaveEntity(ctx context.Context, test model.Test) error {

	//собираем структуру entity
	entity := model.Entity{Id: uuid.New(), Test: test}

	//отдаем данные продюсеру
	if err := service.producer.ProduceEntityToKafka(entity); err != nil {
		return fmt.Errorf("failed to uptdate entity in service: %s", err)
	}
	return nil
}

func (service *entityService) GetAllEntities(ctx context.Context, pathForCache string) ([]model.Entity, error) {

	//прокидываем путь для формирования ключа по которому будем обращаться в кэш
	service.cache.SetPath(pathForCache)

	//обращаемся к кэшу
	if entities := service.cache.GetCache(); entities != nil {
		slog.Info("getting data from cache!", "data", entities)
		return entities, nil
	}

	//если кэш пустой, то идём в базу
	log.Println("cache is empty")
	entities, err := service.repository.GetEntities()
	if err != nil {
		log.Printf("failed to get entities in service: %s", err)
		return entities, err
	}

	//обновляем кэш
	service.cache.UpdateCache(entities)
	return entities, nil
}

func (service *entityService) UpdateEntity(ctx context.Context, entity model.Entity) error {

	//обновляем данные в репозитории
	if err := service.repository.UpdateEntity(entity); err != nil {
		log.Printf("failed to update entities in service: %s", err)
		return err
	}

	//если данные обновились, то чистим кэш
	service.cache.ClearCache()
	return nil
}

func (service *entityService) DeleteEntity(ctx context.Context, id uuid.UUID) error {
	if err := service.repository.DeleteEntity(id); err != nil {
		log.Printf("failed to delete entities in service: %s", err)
		return err
	}

	service.cache.ClearCache()
	return nil
	//service.cache.DeleteEntity(entityForCache)
}
