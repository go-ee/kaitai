package kaitai

import (
	"fmt"
	"log"
)

func BuildReadF(endianConverter EndianReader, length uint8) (ret ReadTo) {
	switch length {
	case 4:
		ret = BuildReadF4(endianConverter)
	case 8:
		ret = BuildReadF8(endianConverter)
	default:
		log.Println(fmt.Sprintf("not supported Native(f,%v)", length))
	}
	return
}

func BuildReadF4(endianConverter EndianReader) ReadTo {
	return func(fillItem *Item, reader *ReaderIO) (err error) {
		if fillItem.Raw, err = reader.ReadBytes(4); err == nil {
			fillItem.SetValue(endianConverter.Float32fromBits(fillItem.Raw))
		}
		return
	}
}

func BuildReadF8(endianConverter EndianReader) ReadTo {
	return func(fillItem *Item, reader *ReaderIO) (err error) {
		if fillItem.Raw, err = reader.ReadBytes(8); err == nil {
			fillItem.SetValue(endianConverter.Float64fromBits(fillItem.Raw))
		}
		return
	}
}
