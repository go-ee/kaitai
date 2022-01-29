package kaitai

import (
	"fmt"
	"strconv"
)

func toUint16(value interface{}) (ret uint16, err error) {
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
