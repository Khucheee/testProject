package service

import (
	"customers_kuber/cache"
	"customers_kuber/listener"
	"customers_kuber/model"
	"customers_kuber/producer"
	"customers_kuber/repository"
	"github.com/google/uuid"
	"log"
)

var entityServiceInstance *entityService

type EntityService interface {
	SaveEntity(test model.Test)
	GetAllEntities() []model.Entity
	GetOneEntity() model.Entity
}

type entityService struct {
	repository repository.EntityRepository
	producer   producer.EntityProducer
	cache      cache.EntityCache
}

func GetEntityService() EntityService {
	if entityServiceInstance != nil {
		return entityServiceInstance
	}
	producer.CreateTopic()

	//запускаем прослушивание кафки
	entityListener := listener.GetEntityListener()
	go entityListener.StartListening()

	//получаем коннект к базе
	entityRepository := repository.GetEntityRepository()

	//полчаем продюсера
	entityProducer := producer.GetEntityProducer()

	//получаем доступ к кэшу
	entityCache := cache.GetEntityCache()

	entityServiceInstance = &entityService{entityRepository, entityProducer, entityCache}
	return entityServiceInstance
}

func (service *entityService) SaveEntity(test model.Test) {
	//собираем структуру entity
	entity := model.Entity{Id: uuid.NewString(), Test: test}

	//отдаем данные продюсеру
	service.producer.ProduceEntityToKafka(entity)

	//удаляем кэш
	service.cache.ClearCache()
}

func (service *entityService) GetAllEntities() []model.Entity {
	//забираем кэш
	if entities := service.cache.GetCache(); entities != nil {
		log.Println("getting data from cache")
		return entities
	}

	//если кэш пустой, то идём в базу
	log.Println("cache is empty")
	entities := service.repository.GetEntities()

	//обновляем кэш
	service.cache.UpdateCache(entities)
	return entities
}

func (service *entityService) GetOneEntity() model.Entity {
	return service.repository.GetOneEntity()
}
