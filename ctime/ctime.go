package ctime

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

var (
	_ json.Marshaler   = (*CTime)(nil)
	_ json.Unmarshaler = (*CTime)(nil)
	_ sql.Scanner      = (*CTime)(nil)
	_ driver.Valuer    = (*CTime)(nil)
)

// return current time, instead of `time.Now()`
// help us to test this package
var now func() time.Time

func init() {
	now = time.Now
}

// CTime is Custom Time, declare all time objects as CTime type in this project
// * use CTime struct directly, DO NOT use this pointer
//
// e.g.
//
//	type User struct{
//		Name string
//		CreatedAt CTime
//		UpdatedAt CTime
//	}
type CTime struct {
	time.Time
}

// New returns the instance with the UTC timezone
func New(t time.Time) CTime {
	return CTime{t.UTC()}
}

// Now returns the UTC timezone by default, replace the Golang std `time.Now()`.
func Now() CTime {
	return CTime{now().UTC()}
}

// String formats with `time.RFC3339`
func (t CTime) String() string {
	return t.Format(time.RFC3339)
}

// In returns a copy of t, but change the timezone with location,
//
// If locName is invalid, returns the origin value
func (t CTime) In(loc *time.Location) CTime {
	if loc == nil {
		return t
	}
	return CTime{t.Time.In(loc)}
}

// MarshalJSON implements the `json.Marshaler`
// returns null if t is zero value
func (t CTime) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte("null"), nil
	}
	return []byte("\"" + t.String() + "\""), nil
}

// UnmarshalJSON implements the `json.Unmarshaler`
// Backward compatible with previous layout.
func (t *CTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		t.Time = time.Time{}
		return nil
	}
	result, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return err
	}
	t.Time = result.UTC()
	return nil
}

// Scan implements the `sql.Scanner`
func (t *CTime) Scan(src interface{}) error {
	if src == nil {
		t.Time = time.Time{}
		return nil
	}
	switch v := src.(type) {
	case time.Time:
		t.Time = v.UTC()
		return nil
	default:
		return fmt.Errorf("invalid time value: %v", src)
	}
}

// Value implements the `driver.Valuer`
func (t CTime) Value() (driver.Value, error) {
	if t.IsZero() {
		return nil, nil
	}
	return t.Time.UTC(), nil
}

// GormDataType specify data type for GORM
// * please read this doc https://gorm.io/docs/data_types.html for more information
func (t CTime) GormDataType() string {
	return "datetime"
}
