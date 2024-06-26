package storage

import (
	"Customers/internal/app"
	"fmt"
)

var storage *Storage

type Storage struct {
	cash []app.Entity
}

func NewStorage() *Storage {
	if storage != nil {
		fmt.Println("Хранилище уже создано")
		return storage
	}
	cash := []app.Entity{}
	storage = &Storage{cash}
	return storage
}

func (s *Storage) SaveCustomer(e app.Entity) {
	s.cash = append(s.cash, e)
}

func (s *Storage) GetCustomers() []app.Entity {
	fmt.Println("Получены значения:", s.cash)
	return s.cash
}
