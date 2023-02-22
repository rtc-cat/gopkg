package ctime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFixedZone(t *testing.T) {
	loc := time.FixedZone("+8", 8*3600)
	n := CTime{now().In(loc)}
	assert.Equal(t, "2022-11-08T22:52:33+08:00", n.String())
}

func TestTimezoneReg(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		want  string
	}{
		{"empty", "", "2022-11-08T14:52:33Z"},
		{"error", "zlkbdasklhnt", "2022-11-08T14:52:33Z"},
		{"utc", "UTC", "2022-11-08T14:52:33Z"},
		{"utc+8", "UTC+8", "2022-11-08T22:52:33+08:00"},
		{"utc-4", "UTC-4", "2022-11-08T10:52:33-04:00"},
		{"Shanghai", "Asia/Shanghai", "2022-11-08T22:52:33+08:00"},
		{"Los_Angeles", "America/Los_Angeles", "2022-11-08T06:52:33-08:00"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			loc := ParseTimezone(tc.input)
			ct := Now()
			ct = ct.In(loc)
			assert.Equal(t, tc.want, ct.String())
		})
	}

}
