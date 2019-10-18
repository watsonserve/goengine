package goengine

import (
    "time"
    "database/sql/driver"
)

type ActionFunc func(http.ResponseWriter, *Session, *http.Request)
type FilterFunc func(http.ResponseWriter, *Session, *http.Request) bool

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
