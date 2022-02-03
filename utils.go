package kaitai

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	day  = time.Minute * 60 * 24
	year = 365 * day
)

func ToUint16(value interface{}) (ret uint16, err error) {
	if number, ok := value.(int); ok {
		ret = uint16(number)
		return
	}

	if number, ok := value.(uint); ok {
		ret = uint16(number)
		return
	}

	if number, ok := value.(uint8); ok {
		ret = uint16(number)
		return
	}

	if number, ok := value.(uint16); ok {
		ret = number
		return
	}

	if number, ok := value.(uint32); ok {
		ret = uint16(number)
		return
	}

	if number, ok := value.(uint64); ok {
		ret = uint16(number)
		return
	}

	str := fmt.Sprintf("%v", value)
	var intValue int
	if intValue, err = strconv.Atoi(str); err == nil {
		ret = uint16(intValue)
	}
	return
}

func DurationToString(d time.Duration) string {
	if d < day {
		return d.String()
	}

	var b strings.Builder
	if d >= year {
		years := d / year
		fmt.Fprintf(&b, "%dy", years)
		d -= years * year
	}

	days := d / day
	d -= days * day
	fmt.Fprintf(&b, "%dd%s", days, d)

	return b.String()
}
