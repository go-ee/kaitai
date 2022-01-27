package kaitai

import (
	"fmt"
	"log"
)

func ReadS(endianConverter EndianReader, length uint8) (ret ReadFix) {
	switch length {
	case 1:
		ret = ReadS1()
	case 2:
		ret = ReadS2(endianConverter)
	case 4:
		ret = ReadS4(endianConverter)
	case 8:
		ret = ReadS8(endianConverter)
	default:
		log.Println(fmt.Sprintf("not supported Native(s,%v)", length))
	}
	return
}

func ReadS8(endianConverter EndianReader) ReadFix {
	return func(reader *Reader) (ret interface{}, err error) {
		var data []byte
		if data, err = reader.ReadBytes(8); err == nil {
			ret = int64(endianConverter.Uint64(data))
		}
		return
	}
}

func ReadS4(endianConverter EndianReader) ReadFix {
	return func(reader *Reader) (ret interface{}, err error) {
		var data []byte
		if data, err = reader.ReadBytes(4); err == nil {
			ret = int32(endianConverter.Uint32(data))
		}
		return
	}
}

func ReadS2(endianConverter EndianReader) ReadFix {
	return func(reader *Reader) (ret interface{}, err error) {
		var data []byte
		if data, err = reader.ReadBytes(2); err == nil {
			ret = int16(endianConverter.Uint16(data))
		}
		return
	}
}

func ReadS1() ReadFix {
	return func(reader *Reader) (ret interface{}, err error) {
		var data []byte
		if data, err = reader.ReadBytes(1); err == nil {
			ret = int8(data[0])
		}
		return
	}
}
