package app

import (
	"github.com/gin-gonic/gin"
)

func (c *Controller) Route() {
	router := gin.Default()
	router.GET("/getall", c.GetCustomers)
	router.POST("create", c.CreateCustomer)
	router.Run("localhost:8080")
}
