package model

type Test struct {
	Name string `json:"Name"`
	Age  int    `json:"Age"`
}
type Entity struct {
	Id   string `gorm:"type:uuid;primaryKey"`
	Test Test   `gorm:"type:json"`
}
