package gorm

import (
	"database/sql/driver"
	"encoding/json"
)

type AInt32 []int32

func (c AInt32) Value() (driver.Value, error) {
	b, err := json.Marshal(c)
	return string(b), err
}
func (c *AInt32) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), c)
}

type AInt64 []int64

func (c AInt64) Value() (driver.Value, error) {
	b, err := json.Marshal(c)
	return string(b), err
}
func (c *AInt64) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), c)
}

type AString []string

func (c AString) Value() (driver.Value, error) {
	b, err := json.Marshal(c)
	return string(b), err
}

func (c *AString) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), c)
}
