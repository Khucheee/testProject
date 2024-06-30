package controller

import (
	"Customers/closer"
	"Customers/model"
	"Customers/service"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
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
	if entityControllerInstance != nil {
		return entityControllerInstance
	}
	entityService := service.GetEntityService()
	entityControllerInstance = &entityController{service: entityService}

	//передаю функцию остановки для graceful shutdown
	closer.CloseFunctions = append(closer.CloseFunctions, entityControllerInstance.CloseController())
	return entityControllerInstance
}

func (controller *entityController) CloseController() func() {
	return func() {
		ctx := context.Background()
		if err := controller.server.Shutdown(ctx); err != nil {
			log.Println("failed to shut server down")
		}
		log.Println("entityController closed successfully")
	}
}

func (controller *entityController) Route() {
	router := gin.Default()
	router.GET("/getAll", controller.GetAllEntities)
	router.GET("/getLastOne", controller.GetOneEntity)
	router.POST("/create", controller.SaveEntity)

	//Запускаем сервер
	server := &http.Server{Addr: ":8080", Handler: router}
	controller.server = server
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Println("Server down:", err)
	}
}

func (controller *entityController) SaveEntity(ctx *gin.Context) {
	//Парсинг полученной json
	test := &model.Test{}
	if err := ctx.BindJSON(test); err != nil {
		log.Println("Wrong JSON received in controller:", err)
	}

	//сохранение данных
	controller.service.SaveEntity(*test)

	//ответ от сервераа
	ctx.Status(http.StatusCreated)
}

func (controller *entityController) GetAllEntities(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, controller.service.GetAllEntities())
}

func (controller *entityController) GetOneEntity(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, controller.service.GetOneEntity())
}
