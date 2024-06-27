package controller

import (
	"Customers/model"
	"Customers/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

var entityControllerInstance *entityController

type EntityController interface {
	Route()
}

type entityController struct {
	service service.EntityService
}

func GetEntityController() EntityController {
	if entityControllerInstance != nil {
		return entityControllerInstance
	}
	entityService := service.GetEntityService()
	entityControllerInstance = &entityController{entityService}
	return entityControllerInstance
}

func (controller *entityController) Route() {
	router := gin.Default()
	router.GET("/getAll", controller.GetAllEntities)
	router.GET("/getLastOne", controller.GetOneEntity)
	router.POST("/create", controller.SaveEntity)
	if err := router.Run("localhost:8080"); err != nil {
		log.Println("Failed to raise the server:", err)
	}
}

func (controller *entityController) SaveEntity(ctx *gin.Context) {

	//Парсинг полученной json
	test := &model.Test{}
	if err := ctx.BindJSON(test); err != nil {
		log.Println("Wrong JSON recieved in controller:", err)
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
