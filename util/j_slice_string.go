package util

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

//JSONSliceString gorm mysql >= 5.7 JSON字段映射的类型
type JSONSliceString []string

func (m JSONSliceString) Value() (driver.Value, error) {
	b, err := json.Marshal(m)
	return string(b), err
}

func (m *JSONSliceString) Scan(input interface{}) error {
	switch v := input.(type) {
	case string:
		return json.Unmarshal([]byte(v), m)
	case []byte:
		return json.Unmarshal(v, m)
	default:
		return errors.New("gorm json 错误类型")
	}
}
