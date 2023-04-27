package tools

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
	"time"
)

type CustomTime struct {
	time.Time
}

const ReceivedDTFormat = "2006-01-02T15:04:05"
const CustomDTLayout = "2006-01-02 15:04:05"

// UnmarshalJSON Parses the json string in the custom format
func (ct *CustomTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`)
	d, err := time.Parse(ReceivedDTFormat, s)
	var r error
	ct.Time, r = time.Parse(CustomDTLayout, d.Format(CustomDTLayout))

	if r != nil {
		panic(r)
	}
	return
}

// MarshalJSON writes a quoted string in the custom format
func (ct CustomTime) MarshalJSON() ([]byte, error) {
	return []byte(ct.String()), nil
}

// String returns the time in the custom format
func (ct *CustomTime) String() string {
	t := time.Time(ct.Time)
	return fmt.Sprintf("%q", t.Format(CustomDTLayout))
}

func (ct *CustomTime) Scan(value interface{}) error {
	if t, ok := value.(time.Time); ok {

		ct.Time = t
		return nil
	}
	return errors.New("Cant convert data")
}

func (ct CustomTime) Value() (driver.Value, error) {

	if ct.Time.IsZero() {
		return nil, nil
	}
	return ct.Time.Format(CustomDTLayout), nil
}
