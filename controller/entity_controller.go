package controller

import (
	"context"
	"customers_kuber/closer"
	"customers_kuber/model"
	"customers_kuber/service"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
)

var entityControllerInstance *entityController

type EntityController interface {
	Route()
	CloseController() func()
}

type entityController struct {
	service service.EntityService
	server  *http.Server //храню сервер для graceful shutdown
}

func GetEntityController() EntityController {

	//ессли сущность уже есть, возвращаем её
	if entityControllerInstance != nil {
		return entityControllerInstance
	}

	//получаем сервис
	entityService, err := service.GetEntityService()
	entityControllerInstance = &entityController{service: entityService}
	if err != nil {
		log.Fatalf("Failed to create controller: %s", err)
	}

	//передаю функцию остановки для graceful shutdown
	closer.CloseFunctions = append(closer.CloseFunctions, entityControllerInstance.CloseController())
	return entityControllerInstance
}

func (controller *entityController) CloseController() func() {
	return func() {
		ctx := context.Background()
		go func() {
			if err := controller.server.Shutdown(ctx); err != nil {
				log.Println("failed to shut server down")
			}
		}()
		log.Println("entityController closed successfully")
	}
}

func (controller *entityController) Route() {

	//роутинг
	router := gin.Default()
	router.GET("/getAll", controller.GetAllEntities)
	router.POST("/create", controller.SaveEntity)
	router.PUT("/:id", controller.UpdateEntity)
	router.DELETE("/:id", controller.DeleteEntity)

	//Запускаем сервер
	server := &http.Server{Addr: ":8080", Handler: router}
	controller.server = server
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Printf("server down: %s", err)
	}
}

func (controller *entityController) SaveEntity(ctx *gin.Context) {

	//Парсинг полученной json
	test := &model.Test{}

	//валидация полей в теле запроса
	if err := ctx.BindJSON(test); err != nil {
		log.Printf("wrong JSON received in saveEntity: %s", err)
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	//сохранение данных, если сервис вернет ошибку, то 500
	if err := controller.service.SaveEntity(*test); err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	//если все ок, 200
	ctx.Status(http.StatusAccepted)
}

func (controller *entityController) GetAllEntities(ctx *gin.Context) {

	//получение всех существующих entity, если вернется ошибка, то 500
	entities, err := controller.service.GetAllEntities(ctx.Request.URL.Path)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, entities)
}

func (controller *entityController) UpdateEntity(ctx *gin.Context) {

	//парсю json и передаю сервису, если не прошли валидации, 400
	entity := &model.EntityForUpdate{}
	if err := ctx.ShouldBindJSON(entity); err != nil {
		log.Println("updateEntity:wrong JSON received in controller:", err)
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	//проверка на uuid, если в URL передан не uuid, то 400
	uuidFromURL, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, "param must be uuid")
		return
	}

	//сравниваем id в теле и в урле, если не совпадают, 400
	if entity.Id != uuidFromURL {
		ctx.JSON(http.StatusTeapot, "wrong id")
		return
	}

	//обновляем данные
	//если вернется not found, то 404
	//при любой другой ошибке 500
	err = controller.service.UpdateEntity(model.Entity{Id: entity.Id, Test: model.Test{Name: entity.Test.Name, Age: entity.Test.Age}})
	if err != nil {
		if err.Error() == "record not found" {
			ctx.Status(http.StatusNotFound)
			return
		}
		log.Println("fai")
		ctx.Status(http.StatusInternalServerError)
		return
	}

	//если дошли до сюда, значит все обновилось, отправляем 200
	ctx.Status(http.StatusOK)

}

func (controller *entityController) DeleteEntity(ctx *gin.Context) {
	uuidFromURL, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "param must be uuid")
	}

	if err := controller.service.DeleteEntity(uuidFromURL); err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
	ctx.Status(http.StatusOK)
}
