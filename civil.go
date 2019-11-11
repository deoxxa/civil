package civil

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type Date struct {
	Year  int
	Month time.Month
	Day   int
}

func DateOf(t time.Time) Date {
	var d Date
	d.Year, d.Month, d.Day = t.Date()
	return d
}

func DateOfNil(t *time.Time) *Date {
	if t == nil {
		return nil
	}

	v := DateOf(*t)

	return &v
}

func (d Date) Format(f string) string {
	return d.In(time.UTC).Format(f)
}

func ParseDate(s string) (Date, error) {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return Date{}, err
	}
	return DateOf(t), nil
}

func (d Date) String() string {
	return fmt.Sprintf("%04d-%02d-%02d", d.Year, d.Month, d.Day)
}

func (d Date) IsValid() bool {
	return DateOf(d.In(time.UTC)) == d
}

func (d Date) In(loc *time.Location) time.Time {
	return time.Date(d.Year, d.Month, d.Day, 0, 0, 0, 0, loc)
}

func (d Date) AddDays(n int) Date {
	return DateOf(d.In(time.UTC).AddDate(0, 0, n))
}

func (d Date) DaysSince(s Date) (days int) {
	deltaUnix := d.In(time.UTC).Unix() - s.In(time.UTC).Unix()
	return int(deltaUnix / 86400)
}

func (d Date) Before(other Date) bool {
	if d.Year != other.Year {
		return d.Year < other.Year
	}
	if d.Month != other.Month {
		return d.Month < other.Month
	}
	return d.Day < other.Day
}

func (d Date) After(other Date) bool {
	return other.Before(d)
}

func (d Date) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}

func (d *Date) UnmarshalText(text []byte) error {
	var err error
	*d, err = ParseDate(string(text))
	return err
}

func (d *Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *Date) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	v, err := ParseDate(s)
	if err != nil {
		return err
	}

	*d = v

	return nil
}

func (d *Date) Scan(src interface{}) error {
	switch v := src.(type) {
	case time.Time:
		*d = DateOf(v)
		return nil
	case string:
		t, err := ParseDate(v)
		if err != nil {
			return err
		}
		*d = t
		return nil
	default:
		return fmt.Errorf("civil.Date.Scan: can't scan into %T", src)
	}
}

func (d Date) Value() (driver.Value, error) {
	return d.String(), nil
}
