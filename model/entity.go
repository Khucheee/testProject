package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Entity struct {
	Id      uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey"`
	Test    Test           `json:"test" gorm:"type:json"`
	Deleted gorm.DeletedAt `json:"-"`
}
