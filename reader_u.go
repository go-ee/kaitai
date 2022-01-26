package kaitai

import (
	"fmt"
	"log"
)

func ReadU(endianConverter EndianReader, length uint8) (ret ReadFix) {
	switch length {
	case 1:
		ret = ReadU1()
	case 2:
		ret = ReadU2(endianConverter)
	case 4:
		ret = ReadU4(endianConverter)
	case 8:
		ret = ReadU8(endianConverter)
	default:
		log.Println(fmt.Sprintf("not supported Native(u,%v)", length))
	}
	return
}

func ReadU1() ReadFix {
	return func(reader ReaderIO) (ret interface{}, err error) {
		var data []byte
		if data, err = reader.ReadBytes(1); err == nil {
			ret = data[0]
		}
		return
	}
}

func ReadU2(endianConverter EndianReader) ReadFix {
	return func(reader ReaderIO) (ret interface{}, err error) {
		var data []byte
		if data, err = reader.ReadBytes(2); err == nil {
			ret = endianConverter.Uint16(data)
		}
		return
	}
}

func ReadU4(endianConverter EndianReader) ReadFix {
	return func(reader ReaderIO) (ret interface{}, err error) {
		var data []byte
		if data, err = reader.ReadBytes(4); err == nil {
			ret = endianConverter.Uint32(data)
		}
		return
	}
}

func ReadU8(endianConverter EndianReader) ReadFix {
	return func(reader ReaderIO) (ret interface{}, err error) {
		var data []byte
		if data, err = reader.ReadBytes(8); err == nil {
			ret = endianConverter.Uint64(data)
		}
		return
	}
}
