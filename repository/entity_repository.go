package repository

import (
	"customers_kuber/closer"
	"customers_kuber/config"
	"customers_kuber/model"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var entityRepositoryInstance *entityRepository

type EntityRepository interface {
	SaveEntity(e model.Entity) error
	GetEntities() ([]model.Entity, error)
	UpdateEntity(update model.Entity) error
	DeleteEntity(uuid2 uuid.UUID) error
}

type entityRepository struct {
	db *gorm.DB
}

func GetEntityRepository() (EntityRepository, error) {
	if entityRepositoryInstance != nil {
		return entityRepositoryInstance, nil
	}

	//устанавливаемм адрес базы
	dbConfig := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		config.PostgresHost, config.PostgresUser, config.PostgresPassword, config.PostgresDatabaseName, config.PostgresPort)

	//открываем соединение
	dbConnect, err := gorm.Open(postgres.Open(dbConfig), &gorm.Config{
		//Logger: logger.Default.LogMode(logger.Info),
	})

	//инициализируем инстанс репозитория
	entityRepositoryInstance = &entityRepository{dbConnect}
	if err != nil {
		return entityRepositoryInstance, fmt.Errorf("failed to create database connection: %s", err)
	}

	//делаем миграцию
	err = dbConnect.AutoMigrate(&model.Entity{})
	if err != nil {
		return entityRepositoryInstance, fmt.Errorf("migration failed: %s", err)
	}

	//передаем функцию закрытия в клозер для graceful shutdown
	closer.CloseFunctions = append(closer.CloseFunctions, entityRepositoryInstance.CloseEntityRepository())
	return entityRepositoryInstance, nil
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

func (repository *entityRepository) SaveEntity(e model.Entity) error {

	//сохраняю Entity в базу
	result := repository.db.Create(&e)
	if result.Error != nil {
		return fmt.Errorf("failed to save entity into repository: %s", result.Error)
	}
	return nil
}

func (repository *entityRepository) GetEntities() ([]model.Entity, error) {

	//запрашиваем в базе все записи Entities
	var entities []model.Entity
	if result := repository.db.Find(&entities); result.Error != nil {
		log.Printf("failed to get all entities from repository: %s", result.Error)
		return entities, result.Error
	}
	return entities, nil
}

func (repository *entityRepository) UpdateEntity(entity model.Entity) error {
	//проверяем наличие записи в базе
	checkExistence := model.Entity{}
	result := repository.db.Where("id=?", entity.Id).First(&checkExistence)
	if result.Error != nil {
		log.Printf("failed to find row while updating entity: %s", result.Error)
		return result.Error
	}

	//если запись нашлась, то обновляю её
	if result = repository.db.Model(&entity).Where("id=?", entity.Id).Updates(map[string]interface{}{"id": entity.Id, "Test": entity.Test}); result.Error != nil {
		log.Printf("failed to update data in repository: %s", result.Error)
		return result.Error
	}
	log.Println("entity updated successfully")
	return nil
}

func (repository *entityRepository) DeleteEntity(id uuid.UUID) error {
	entity := model.Entity{}
	if result := repository.db.Where("Id = ?", id).Delete(&entity); result.Error != nil {
		err := fmt.Errorf("failed to delete data from repository: %s", result.Error)
		log.Println(err)
		return err
	}
	log.Println("deleting from database ended successfully")
	return nil
}

/*
func (repository *entityRepository) DeleteEntity(id string) model.Entity {
	entity := model.Entity{}
	if result := repository.db.Clauses(clause.Returning{}).Where("Id = ?", id).Delete(&entity); result.Error != nil {
		log.Println(result.Error)
	}
	log.Println("deleted from database ended: ", entity)
	return entity
}
*/
