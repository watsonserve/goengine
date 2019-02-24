package goengine

import (
    "time"
    "database/sql/driver"
)

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
