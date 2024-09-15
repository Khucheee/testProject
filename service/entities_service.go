package service

import (
	"context"
	"customers_kuber/cache"
	"customers_kuber/listener"
	"customers_kuber/logger"
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

// todo спросить в Давида, надо ли хранить эти сущности в сервисе, ведь можно обращаться к ним через "Get(сущность)"
// я за то, чтобы оставить так, ведь сильно меньше обработок ошибок при создании сервиса получается

// entityService принимает данные от контроллера и отвечает за бизнес-логику приложения
type entityService struct {
	repository repository.EntityRepository
	producer   producer.EntityProducer
	cache      cache.EntityCache
}

// GetEntityService создает сервис, внутри создает коннекты к базе,
// также создает топик кафки и запускает лисенера,
// реализовано через синглтон
func GetEntityService() (EntityService, error) {

	//проверяем наличие EntityService
	if entityServiceInstance != nil {
		return entityServiceInstance, nil
	}

	//создаем топик в кафке, если не получится, приложение не запустится
	if err := producer.CreateTopic(); err != nil {
		return entityServiceInstance, err
	}

	//создаем листенера и запускаем его, если не получится, приложение не запустится
	entityListener, err := listener.GetEntityListener()
	if err != nil {
		return entityServiceInstance, err
	}
	go entityListener.StartListening()

	//получаем коннект к базе, если не получится, приложение не запустится
	entityRepository, err := repository.GetEntityRepository()
	if err != nil {
		return entityServiceInstance, err
	}

	//полчаем продюсера, если не получится, приложение не запустится
	entityProducer, err := producer.GetEntityProducer()
	if err != nil {
		return entityServiceInstance, err
	}

	//получаем доступ к кэшу, если не получится, приложение будет работать
	entityCache, err := cache.GetEntityCache()
	if err != nil {
		slog.Error("failed to get cache in service")
	}

	//собираем структуру и возвращаем её
	entityServiceInstance = &entityService{entityRepository, entityProducer, entityCache}

	return entityServiceInstance, nil
}

// SaveEntity отвечает за сохраниние сущности в приложении, передает entity продюсеру kafka
func (service *entityService) SaveEntity(ctx context.Context, test model.Test) error {

	//собираем структуру entity и отдаем данные продюсеру
	entity := model.Entity{Id: uuid.New(), Test: test}
	if err := service.producer.ProduceEntityToKafka(entity); err != nil {
		return fmt.Errorf("failed to uptdate entity in service: %s", err)
	}
	slog.InfoContext(logger.WithLogValues(ctx, entity), "new entity successfully sent to kafka")

	return nil
}

// GetAllEntities возвращает все ранее созданные entities,
// сначала проверяет их в кэше, если от туда вернется nil,
// пойдет забирать значения из базы, затем обновит кэш
func (service *entityService) GetAllEntities(ctx context.Context, pathForCache string) ([]model.Entity, error) {

	//прокидываем путь для формирования ключа по которому будем обращаться в кэш
	service.cache.SetPath(pathForCache)

	//обращаемся к кэшу, если кэш есть, возвращаем данные
	if entities := service.cache.GetCache(ctx); entities != nil {
		slog.InfoContext(logger.WithLogValues(ctx, entities), "service received not nil cache, returning values from cache")
		return entities, nil
	}
	slog.Info("service received nil cache, getting data from DB")

	//если кэш пустой, то идём в базу
	entities, err := service.repository.GetEntities()
	if err != nil {
		return entities, err
	}
	slog.InfoContext(logger.WithLogValues(ctx, entities), "service successfully received data from DB, start updating cache")

	//обновляем кэш
	//todo возможно ошибки кэша можно обрабатывать в сервисе, или же сделать кастомные ошибки и обрабатывать их в контроллере
	service.cache.UpdateCache(ctx, entities)
	slog.Info("cache updated successfully")

	return entities, nil
}

func (service *entityService) UpdateEntity(ctx context.Context, entity model.Entity) error {

	//обновляем данные в репозитории
	if err := service.repository.UpdateEntity(entity); err != nil {
		return err
	}
	slog.InfoContext(logger.WithLogValues(ctx, entity), "service successfully updated entity in repository")

	//если данные обновились, то чистим кэш
	service.cache.ClearCache(ctx)
	slog.Info("service cleaned cache successfully")

	return nil
}

func (service *entityService) DeleteEntity(ctx context.Context, id uuid.UUID) error {
	if err := service.repository.DeleteEntity(id); err != nil {
		log.Printf("failed to delete entities in service: %s", err)
		return err
	}
	slog.InfoContext(logger.WithLogValues(ctx, id), "service successfully deleted entity in repository")

	service.cache.ClearCache(ctx)
	slog.Info("service cleaned cache successfully")
	return nil
	//service.cache.DeleteEntity(entityForCache)
}
