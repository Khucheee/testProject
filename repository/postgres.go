package repository

import (
	"Customers/model"
	"fmt"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *Postgres

type Postgres struct {
	connect *gorm.DB
}

func NewPostgres() *Postgres {
	if Db != nil {
		return Db
	} else {
		//открываем соединение
		dsn := "host=localhost user=admin password=root dbname=mydatabase port=5432 sslmode=disable"
		dbCon, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			fmt.Println("Не открылось соединение с базой", err)
		}
		//делаем миграцию
		err = dbCon.AutoMigrate(&model.Entity{})
		if err != nil {
			fmt.Println("Упала миграция postgres")
		}
		//закрываем коннект,если что-то идет не так
		fmt.Println("Инициализация базы прошла успешно", Db)
		Db = &Postgres{dbCon}
		return Db
	}
}

func (p *Postgres) SaveCustomer(e model.Entity) {
	fmt.Println("Запущен метод сохранения в базу")
	fmt.Println(e)
	result := p.connect.Create(&e)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
}

func (p *Postgres) GetCustomers() []model.Entity {
	var entities []model.Entity

	return entities
}
