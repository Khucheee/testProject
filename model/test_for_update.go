package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

//в отличие от test, поля необязательные для парсинга, можно было отлавливать ошибки типа required field
//но так показалось красивее

type TestForUpdate struct {
	Name *string `json:"name"`
	Age  *int    `json:"age" binding:"min=0,max=100"`
}

func (test *TestForUpdate) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to parse test from base into struct")
	}
	return json.Unmarshal(bytes, &test)
}

func (test *TestForUpdate) Value() (driver.Value, error) {
	return json.Marshal(test)
}
