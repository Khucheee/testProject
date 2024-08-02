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

func GetEntityController() (EntityController, error) {

	//ессли сущность уже есть, возвращаем её
	if entityControllerInstance != nil {
		return entityControllerInstance, nil
	}

	//получаем сервис
	entityService, err := service.GetEntityService()
	entityControllerInstance = &entityController{service: entityService}
	if err != nil {
		return entityControllerInstance, err
	}

	//передаю функцию остановки для graceful shutdown
	closer.CloseFunctions = append(closer.CloseFunctions, entityControllerInstance.CloseController())
	return entityControllerInstance, nil
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
	if err := ctx.BindJSON(test); err != nil {
		log.Printf("wrong JSON received in saveEntity: %s", err)
		ctx.Status(http.StatusBadRequest)
	}

	//сохранение данных
	if err := controller.service.SaveEntity(*test); err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	//если все ок
	ctx.Status(http.StatusAccepted)
}

func (controller *entityController) GetAllEntities(ctx *gin.Context) {
	entities, err := controller.service.GetAllEntities(ctx.Request.URL.Path)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, entities)
}

func (controller *entityController) UpdateEntity(ctx *gin.Context) {

	//парсю json и передаю сервису, если не получилось распарсить, ошибка
	entity := &model.Entity{}
	if err := ctx.BindJSON(entity); err != nil {
		log.Println("updateEntity:wrong JSON received in controller:", err)
	}

	//проверка на uuid
	if err := uuid.Validate(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusBadRequest, "param is not uuid")
	}

	//сравниваем id в теле и в урле, если не совпадают, ошибка
	if entity.Id != ctx.Param("id") {
		ctx.JSON(http.StatusTeapot, "wrong id")
		return
	}

	//обновляем данные
	test := entity.Test
	err := controller.service.UpdateEntity(model.Entity{Id: ctx.Param("id"), Test: test})
	if err != nil {
		if err.Error() == "record not found" {
			ctx.Status(http.StatusNotFound)
			return
		}
		ctx.Status(http.StatusInternalServerError)
		return
	}

	//если дошли до сюда, значит все обновилось, отправляем 200
	ctx.Status(http.StatusOK)

}

func (controller *entityController) DeleteEntity(ctx *gin.Context) {
	if err := controller.service.DeleteEntity(ctx.Param("id")); err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
	ctx.Status(http.StatusOK)
}
