package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

//Отдельная структура нужна для валидаций

type EntityForUpdate struct {
	Id      uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey" binding:"required"`
	Test    TestForUpdate  `json:"test" gorm:"type:json" binding:"required"`
	Deleted gorm.DeletedAt `json:"-"`
}
