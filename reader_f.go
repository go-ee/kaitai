package kaitai

import (
	"fmt"
	"log"
	"math"
)

func ReadF(endianConverter EndianReader, length uint8) (ret ReadTo) {
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

func ReadF4(endianConverter EndianReader) ReadTo {
	return func(fillItem *Item, reader *Reader) (err error) {
		fillItem.SetStartPos(reader)
		if fillItem.Raw, err = reader.ReadBytes(4); err == nil {
			fillItem.Value = math.Float32frombits(endianConverter.Uint32(fillItem.Raw))
		}
		fillItem.SetEndPos(reader)
		return
	}
}

func ReadF8(endianConverter EndianReader) ReadTo {
	return func(fillItem *Item, reader *Reader) (err error) {
		fillItem.SetStartPos(reader)
		if fillItem.Raw, err = reader.ReadBytes(8); err == nil {
			fillItem.Value = math.Float64frombits(endianConverter.Uint64(fillItem.Raw))
		}
		fillItem.SetEndPos(reader)
		return
	}
}
