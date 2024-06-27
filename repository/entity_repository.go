package repository

import (
	"Customers/model"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var entityRepositoryInstance *entityRepository

type EntityRepository interface {
	SaveEntity(e model.Entity)
	GetEntities() []model.Entity
	GetOneEntity() model.Entity
}

type entityRepository struct {
	db *gorm.DB
}

func GetEntityRepository() EntityRepository {
	if entityRepositoryInstance != nil {
		return entityRepositoryInstance
	}

	//открываем соединение
	dbConfig := "host=localhost user=admin password=root dbname=mydatabase port=5432 sslmode=disable"
	dbConnect, err := gorm.Open(postgres.Open(dbConfig), &gorm.Config{})
	if err != nil {
		log.Println("failed to create database connection:", err)
	}

	//делаем миграцию
	err = dbConnect.AutoMigrate(&model.Entity{})
	if err != nil {
		log.Println("migration failed:", err)
	}

	entityRepositoryInstance = &entityRepository{dbConnect}
	return entityRepositoryInstance
}

func (repository *entityRepository) SaveEntity(e model.Entity) {

	//сохраняю Entity в базу
	result := repository.db.Create(&e)
	if result.Error != nil {
		log.Println(result.Error)
	}
}

func (repository *entityRepository) GetEntities() []model.Entity {

	//запрашиваем в базе все записи Entities
	var entities []model.Entity
	if result := repository.db.Find(&entities); result.Error != nil {
		log.Println(result.Error)
		return entities
	}
	return entities
}

func (repository *entityRepository) GetOneEntity() model.Entity {
	var entity model.Entity
	if result := repository.db.First(&entity); result.Error != nil {
		log.Println(result.Error)
		return entity
	}
	return entity
}
