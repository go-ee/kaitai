package kaitai

import (
	"fmt"
	"log"
)

func ReadS(endianConverter EndianReader, length uint8) (ret ReadTo) {
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

func ReadS8(endianConverter EndianReader) ReadTo {
	return func(fillItem *Item, reader *Reader) (err error) {
		fillItem.SetStartPos(reader)
		if fillItem.Raw, err = reader.ReadBytes(8); err == nil {
			fillItem.Value = int64(endianConverter.Uint64(fillItem.Raw))
		}
		fillItem.SetEndPos(reader)
		return
	}
}

func ReadS4(endianConverter EndianReader) ReadTo {
	return func(fillItem *Item, reader *Reader) (err error) {
		fillItem.SetStartPos(reader)
		if fillItem.Raw, err = reader.ReadBytes(4); err == nil {
			fillItem.Value = int32(endianConverter.Uint32(fillItem.Raw))
		}
		fillItem.SetEndPos(reader)
		return
	}
}

func ReadS2(endianConverter EndianReader) ReadTo {
	return func(fillItem *Item, reader *Reader) (err error) {
		fillItem.SetStartPos(reader)
		if fillItem.Raw, err = reader.ReadBytes(2); err == nil {
			fillItem.Value = int16(endianConverter.Uint16(fillItem.Raw))
		}
		fillItem.SetEndPos(reader)
		return
	}
}

func ReadS1() ReadTo {
	return func(fillItem *Item, reader *Reader) (err error) {
		fillItem.SetStartPos(reader)
		if fillItem.Raw, err = reader.ReadBytes(1); err == nil {
			fillItem.Value = int8(fillItem.Raw[0])
		}
		fillItem.SetEndPos(reader)
		return
	}
}
