package goengine

import (
	"database/sql/driver"
	"net/http"
	"time"
)

type ActionFunc func(http.ResponseWriter, *http.Request)

// @return go on
type FilterFunc func(http.ResponseWriter, *http.Request) bool

type NullTime struct {
	Time  time.Time
	Valid bool
}

func (this *NullTime) Scan(value interface{}) error {
	this.Time, this.Valid = value.(time.Time)
	return nil
}

func (this *NullTime) Value() (driver.Value, error) {
	if !this.Valid {
		return nil, nil
	}
	return this.Time, nil
}
