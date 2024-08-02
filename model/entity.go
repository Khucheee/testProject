package model

import "gorm.io/gorm"

type Entity struct {
	Id      string         `gorm:"type:uuid;primaryKey"`
	Test    Test           `gorm:"type:json"`
	Deleted gorm.DeletedAt `json:"-"`
}
