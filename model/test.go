package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Test struct {
	Name *string `json:"name" binding:"required"`
	Age  *int    `json:"age" binding:"required,min=0,max=100"`
}

func (test *Test) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to parse test from base into struct")
	}
	return json.Unmarshal(bytes, &test)
}

func (test *Test) Value() (driver.Value, error) {
	return json.Marshal(test)
}
