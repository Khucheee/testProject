package PostgreDB

import (
	"Customers/internal/app"
	"fmt"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *Postgres

type Postgres struct {
	db *gorm.DB
}

func NewPostgres() *Postgres {
	if db != nil {
		fmt.Println("Коннект уже создан")
		return db
	} else {
		//открываем соединение
		dsn := "host=localhost user=admin password=root dbname=mydatabase port=5432 sslmode=disable"
		dbCon, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			fmt.Println("Не открылось соединение с базой", err)
		}
		//делаем миграцию
		err = dbCon.AutoMigrate(&app.Entity{})
		if err != nil {
			fmt.Println("Упала миграция postgres")
		}
		//закрываем коннект,если что-то идет не так
		fmt.Println("Инициализация базы прошла успешно", db)
		db = &Postgres{dbCon}
		return db
	}
}

func (p *Postgres) SaveCustomer(e app.Entity) {
	fmt.Println("Запущен метод сохранения в базу")
	fmt.Println(e)
	result := p.db.Create(&e)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
}

/*func (p *Postgres) Restore() []Entity {
	var entity Entity
	result := p.db.First(&enti)
	fmt.Println(entity)
	if result.Error != nil {
		fmt.Println("Не удалось восстановить кэш из базы", result.Error)
	}
	var entities []Entity
	return entities
}
*/
