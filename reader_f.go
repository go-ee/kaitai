package kaitai

import (
	"fmt"
	"log"
	"math"
)

func ReadF(endianConverter EndianReader, length uint8) (ret ReadFix) {
	switch length {
	case 4:
		ret = ReadF4(endianConverter)
	case 8:
		ret = ReadF8(endianConverter)
	default:
		log.Println(fmt.Sprintf("not supported Native(f,%v)", length))
	}
	return
}

func ReadF4(endianConverter EndianReader) ReadFix {
	return func(reader Reader) (ret interface{}, err error) {
		var data []byte
		if data, err = reader.ReadBytes(4); err == nil {
			ret = math.Float32frombits(endianConverter.Uint32(data))
		}
		return
	}
}

func ReadF8(endianConverter EndianReader) ReadFix {
	return func(reader Reader) (ret interface{}, err error) {
		var data []byte
		if data, err = reader.ReadBytes(8); err == nil {
			ret = math.Float64frombits(endianConverter.Uint64(data))
		}
		return
	}
}
