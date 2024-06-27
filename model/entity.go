package model

type Entity struct {
	Id   string `gorm:"type:uuid;primaryKey"`
	Test Test   `gorm:"type:json"`
}
