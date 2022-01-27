package kaitai

import (
	"github.com/sirupsen/logrus"
)

func ReadU(endianConverter EndianReader, length uint8) (ret ReadTo) {
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
		logrus.Infof("not supported Native(u,%v)", length)
	}
	return
}

func ReadU1() ReadTo {
	return func(fillItem *Item, reader *Reader) (err error) {
		fillItem.SetStartPos(reader)
		if fillItem.Raw, err = reader.ReadBytes(1); err == nil {
			fillItem.SetValue(fillItem.Raw[0])
		}
		fillItem.SetEndPos(reader)
		return
	}
}

func ReadU2(endianConverter EndianReader) ReadTo {
	return func(fillItem *Item, reader *Reader) (err error) {
		fillItem.SetStartPos(reader)
		if fillItem.Raw, err = reader.ReadBytes(2); err == nil {
			fillItem.SetValue(endianConverter.Uint16(fillItem.Raw))
		}
		fillItem.SetEndPos(reader)
		return
	}
}

func ReadU4(endianConverter EndianReader) ReadTo {
	return func(fillItem *Item, reader *Reader) (err error) {
		fillItem.SetStartPos(reader)
		if fillItem.Raw, err = reader.ReadBytes(4); err == nil {
			fillItem.SetValue(endianConverter.Uint32(fillItem.Raw))
		}
		fillItem.SetEndPos(reader)
		return
	}
}

func ReadU8(endianConverter EndianReader) ReadTo {
	return func(fillItem *Item, reader *Reader) (err error) {
		fillItem.SetStartPos(reader)
		if fillItem.Raw, err = reader.ReadBytes(8); err == nil {
			fillItem.SetValue(endianConverter.Uint64(fillItem.Raw))
		}
		fillItem.SetEndPos(reader)
		return
	}
}
