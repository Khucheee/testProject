package repository

import (
	"customers_kuber/closer"
	"customers_kuber/config"
	"customers_kuber/model"
	"fmt"
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

	//устанавливаемм адрес базы
	dbConfig := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		config.PostgresHost, config.PostgresUser, config.PostgresPassword, config.PostgresDatabaseName, config.PostgresPort)

	//открываем соединение
	dbConnect, err := gorm.Open(postgres.Open(dbConfig), &gorm.Config{})
	if err != nil {
		log.Println("failed to create database connection:", err)
	}

	//делаем миграцию
	err = dbConnect.AutoMigrate(&model.Entity{})
	if err != nil {
		log.Println("migration failed:", err)
	}

	//инициализируем инстанс репозитория
	entityRepositoryInstance = &entityRepository{dbConnect}

	//передаем функцию закрытия в клозер для graceful shutdown
	closer.CloseFunctions = append(closer.CloseFunctions, entityRepositoryInstance.CloseEntityRepository())
	return entityRepositoryInstance
}

func (repository *entityRepository) CloseEntityRepository() func() {
	dbInterface, err := repository.db.DB()
	if err != nil {
		log.Println("failed to get DB interface while closing connection:", err)
	}
	return func() {
		if err := dbInterface.Close(); err != nil {
			log.Println("failed while closing DB connection:", err)
			return
		}
		log.Println("entityRepository closed successfully")
	}
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
