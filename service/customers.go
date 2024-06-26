package service

import (
	"Customers/model"
	"Customers/producer"
	"Customers/repository"
	"github.com/google/uuid"
)

func SaveCustomer(test model.Test) {
	//здесь надо вызывать продюсера и передавать ему entity
	entity := model.Entity{Id: uuid.NewString(), Test: test}
	producer.ProduceToKafka(entity)
}

func GetAllCustomers() []model.Entity {
	//здесь надо вызвать метод в базе, который вернет все данные
	entities := repository.Db.GetCustomers()
	return entities
}
