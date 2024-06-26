package controller

import (
	"Customers/model"
	"Customers/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

var controller *Controller

type Controller struct{}

func NewController() *Controller {
	if controller == nil {
		return &Controller{}
	}
	return controller
}

func (c *Controller) SaveCustomer(ctx *gin.Context) {
	//Парсинг полученной json
	test := &model.Test{}
	if err := ctx.BindJSON(test); err != nil {
		fmt.Println("Получена кривая jsonка, не получилось её распарсить в хендлере")
	}
	//сохранение данных
	service.SaveCustomer(*test)
	//ответ от сервераа
	ctx.Status(http.StatusCreated)
}

func (c *Controller) GetCustomers(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, service.GetAllCustomers()) //здесь надо ходить в service
}
