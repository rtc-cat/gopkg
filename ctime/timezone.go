package ctime

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

func TimezoneHeader() string {
	return "X-Timezone"
}

func TimezoneKey() string {
	return "timezone"
}

var timezoneReg *regexp.Regexp

func init() {
	timezoneReg = regexp.MustCompile(`^UTC$|^UTC[\+-][0-9]{1,2}$`)
}

// ParseTimezone parse timezone with string. Support location name and utc offset
// e.g.
// ParseTimezone("Asia/Shanghai")
// ParseTimezone("America/Los_Angeles")
// ParseTimezone("UTC")
// ParseTimezone("UTC+8")
// ParseTimezone("UTC-6")
func ParseTimezone(timezone string) *time.Location {
	timezone = strings.TrimSpace(timezone)
	loc, err := time.LoadLocation(timezone)
	if err == nil {
		return loc
	}

	if !timezoneReg.MatchString(timezone) {
		return time.UTC
	}

	if timezone == "UTC" {
		return time.UTC
	}

	var offset int
	offsetStr := strings.TrimPrefix(timezone, "UTC")
	sign := offsetStr[0]
	digitStr := offsetStr[1:]
	digit, err := strconv.Atoi(digitStr)
	if err != nil {
		return time.UTC
	}
	switch sign {
	case '+':
		offset = digit * 3600
	case '-':
		offset = -digit * 3600
	default:
		return time.UTC
	}
	return time.FixedZone(timezone, offset)
}
