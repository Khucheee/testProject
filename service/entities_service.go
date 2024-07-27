package service

import (
	"customers_kuber/listener"
	"customers_kuber/model"
	"customers_kuber/producer"
	"customers_kuber/repository"
	"github.com/google/uuid"
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
	entityServiceInstance = &entityService{entityRepository, entityProducer}
	return entityServiceInstance
}

func (service *entityService) SaveEntity(test model.Test) {
	entity := model.Entity{Id: uuid.NewString(), Test: test}
	service.producer.ProduceEntityToKafka(entity)
}

func (service *entityService) GetAllEntities() []model.Entity {
	entities := service.repository.GetEntities()
	return entities
}

func (service *entityService) GetOneEntity() model.Entity {
	return service.repository.GetOneEntity()
}
