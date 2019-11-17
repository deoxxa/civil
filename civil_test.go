package civil

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDates(t *testing.T) {
	for _, test := range []struct {
		date     Date
		loc      *time.Location
		wantStr  string
		wantTime time.Time
	}{
		{
			date:     Date{2014, 7, 29},
			loc:      time.Local,
			wantStr:  "2014-07-29",
			wantTime: time.Date(2014, time.July, 29, 0, 0, 0, 0, time.Local),
		},
		{
			date:     DateOf(time.Date(2014, 8, 20, 15, 8, 43, 1, time.Local)),
			loc:      time.UTC,
			wantStr:  "2014-08-20",
			wantTime: time.Date(2014, 8, 20, 0, 0, 0, 0, time.UTC),
		},
		{
			date:     DateOf(time.Date(999, time.January, 26, 0, 0, 0, 0, time.Local)),
			loc:      time.UTC,
			wantStr:  "0999-01-26",
			wantTime: time.Date(999, 1, 26, 0, 0, 0, 0, time.UTC),
		},
	} {
		if got := test.date.String(); got != test.wantStr {
			t.Errorf("%#v.String() = %q, want %q", test.date, got, test.wantStr)
		}
		if got := test.date.In(test.loc); !got.Equal(test.wantTime) {
			t.Errorf("%#v.In(%v) = %v, want %v", test.date, test.loc, got, test.wantTime)
		}
	}
}

func TestDateIsValid(t *testing.T) {
	for _, test := range []struct {
		date Date
		want bool
	}{
		{Date{2014, 7, 29}, true},
		{Date{2000, 2, 29}, true},
		{Date{10000, 12, 31}, true},
		{Date{1, 1, 1}, true},
		{Date{0, 1, 1}, true},  // year zero is OK
		{Date{-1, 1, 1}, true}, // negative year is OK
		{Date{1, 0, 1}, false},
		{Date{1, 1, 0}, false},
		{Date{2016, 1, 32}, false},
		{Date{2016, 13, 1}, false},
		{Date{1, -1, 1}, false},
		{Date{1, 1, -1}, false},
	} {
		got := test.date.IsValid()
		if got != test.want {
			t.Errorf("%#v: got %t, want %t", test.date, got, test.want)
		}
	}
}

func TestParseDate(t *testing.T) {
	for _, test := range []struct {
		str  string
		want Date // if empty, expect an error
	}{
		{"2016-01-02", Date{2016, 1, 2}},
		{"2016-12-31", Date{2016, 12, 31}},
		{"0003-02-04", Date{3, 2, 4}},
		{"999-01-26", Date{}},
		{"", Date{}},
		{"2016-01-02x", Date{}},
		{"2016-01-02T00:00:00.000Z", Date{2016, 1, 2}},
		{"2016-01-02T23:59:59.999Z", Date{2016, 1, 2}},
	} {
		got, err := ParseDate(test.str)
		if got != test.want {
			t.Errorf("ParseDate(%q) = %+v, want %+v", test.str, got, test.want)
		}
		if err != nil && test.want != (Date{}) {
			t.Errorf("Unexpected error %v from ParseDate(%q)", err, test.str)
		}
	}
}

func TestDateArithmetic(t *testing.T) {
	for _, test := range []struct {
		desc  string
		start Date
		end   Date
		days  int
	}{
		{
			desc:  "zero days noop",
			start: Date{2014, 5, 9},
			end:   Date{2014, 5, 9},
			days:  0,
		},
		{
			desc:  "crossing a year boundary",
			start: Date{2014, 12, 31},
			end:   Date{2015, 1, 1},
			days:  1,
		},
		{
			desc:  "negative number of days",
			start: Date{2015, 1, 1},
			end:   Date{2014, 12, 31},
			days:  -1,
		},
		{
			desc:  "full leap year",
			start: Date{2004, 1, 1},
			end:   Date{2005, 1, 1},
			days:  366,
		},
		{
			desc:  "full non-leap year",
			start: Date{2001, 1, 1},
			end:   Date{2002, 1, 1},
			days:  365,
		},
		{
			desc:  "crossing a leap second",
			start: Date{1972, 6, 30},
			end:   Date{1972, 7, 1},
			days:  1,
		},
		{
			desc:  "dates before the unix epoch",
			start: Date{101, 1, 1},
			end:   Date{102, 1, 1},
			days:  365,
		},
	} {
		if got := test.start.AddDays(test.days); got != test.end {
			t.Errorf("[%s] %#v.AddDays(%v) = %#v, want %#v", test.desc, test.start, test.days, got, test.end)
		}
		if got := test.end.DaysSince(test.start); got != test.days {
			t.Errorf("[%s] %#v.Sub(%#v) = %v, want %v", test.desc, test.end, test.start, got, test.days)
		}
	}
}

func TestDateBefore(t *testing.T) {
	for _, test := range []struct {
		d1, d2 Date
		want   bool
	}{
		{Date{2016, 12, 31}, Date{2017, 1, 1}, true},
		{Date{2016, 1, 1}, Date{2016, 1, 1}, false},
		{Date{2016, 12, 30}, Date{2016, 12, 31}, true},
	} {
		if got := test.d1.Before(test.d2); got != test.want {
			t.Errorf("%v.Before(%v): got %t, want %t", test.d1, test.d2, got, test.want)
		}
	}
}

func TestDateAfter(t *testing.T) {
	for _, test := range []struct {
		d1, d2 Date
		want   bool
	}{
		{Date{2016, 12, 31}, Date{2017, 1, 1}, false},
		{Date{2016, 1, 1}, Date{2016, 1, 1}, false},
		{Date{2016, 12, 30}, Date{2016, 12, 31}, false},
	} {
		if got := test.d1.After(test.d2); got != test.want {
			t.Errorf("%v.After(%v): got %t, want %t", test.d1, test.d2, got, test.want)
		}
	}
}

func TestMarshalJSON(t *testing.T) {
	for _, test := range []struct {
		value interface{}
		want  string
	}{
		{Date{1987, 4, 15}, `"1987-04-15"`},
	} {
		bgot, err := json.Marshal(test.value)
		if err != nil {
			t.Fatal(err)
		}
		if got := string(bgot); got != test.want {
			t.Errorf("%#v: got %s, want %s", test.value, got, test.want)
		}
	}
}

func TestUnmarshalJSON(t *testing.T) {
	var d Date
	for _, test := range []struct {
		data string
		ptr  interface{}
		want interface{}
	}{
		{`"1987-04-15"`, &d, &Date{1987, 4, 15}},
		{`"1987-04-\u0031\u0035"`, &d, &Date{1987, 4, 15}},
	} {
		if err := json.Unmarshal([]byte(test.data), test.ptr); err != nil {
			t.Fatalf("%s: %v", test.data, err)
		}
		assert.Equal(t, test.want, test.ptr)
	}

	for _, bad := range []string{"", `""`, `"bad"`, `"1987-04-15x"`,
		`19870415`,     // a JSON number
		`11987-04-15x`, // not a JSON string

	} {
		if json.Unmarshal([]byte(bad), &d) == nil {
			t.Errorf("%q, Date: got nil, want error", bad)
		}
	}
}
