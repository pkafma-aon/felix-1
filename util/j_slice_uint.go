package util

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

//JSONSliceUint gorm mysql >= 5.7 JSON字段映射的类型
type JSONSliceUint []uint

//Value .
func (o JSONSliceUint) Value() (driver.Value, error) {
	b, err := json.Marshal(o)
	return string(b), err
}

//Scan .
func (o *JSONSliceUint) Scan(input interface{}) error {
	switch v := input.(type) {
	case string:
		return json.Unmarshal([]byte(v), o)
	case []byte:
		return json.Unmarshal(v, o)
	default:
		return errors.New("gorm json 错误类型")
	}
}
