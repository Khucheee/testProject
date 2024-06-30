package service

import (
	"Customers/listener"
	"Customers/model"
	"Customers/producer"
	"Customers/repository"
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
	producer.CreateTopic("json_topic", 1, 1)
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
