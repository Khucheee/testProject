package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type Test struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}
type Entity struct {
	ID   string `gorm:"type:uuid;primaryKey"`
	Test Test   `gorm:"type:json"`
}

func (c *Controller) CreateCustomer(ctx *gin.Context) {
	//Парсинг полученной json
	t := &Test{}
	err := ctx.BindJSON(t)
	if err != nil {
		fmt.Println("Получена кривая jsonка, не получилось её распарсить в хендлере")
	}
	//собираем entity
	e := Entity{ID: uuid.NewString(), Test: *t}
	//сохраняем в кэш
	c.storage.SaveCustomer(e)
	//закидываем в топик
	c.producer.SendMessage(e)
	time.Sleep(time.Second * 1)
	c.consumer.Read()
	ctx.Status(http.StatusCreated)
}
func (c *Controller) GetCustomers(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, c.storage.GetCustomers())
}
