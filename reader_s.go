package kaitai

import (
	"github.com/sirupsen/logrus"
)

func BuildReadS(endianConverter EndianReader, length uint8) (ret ReadTo) {
	switch length {
	case 1:
		ret = BuildReadS1()
	case 2:
		ret = BuildReadS2(endianConverter)
	case 4:
		ret = BuildReadS4(endianConverter)
	case 8:
		ret = BuildReadS8(endianConverter)
	default:
		logrus.Infof("not supported Native(s,%v)", length)
	}
	return
}

func BuildReadS8(endianConverter EndianReader) ReadTo {
	return func(fillItem *Item, reader *Reader) (err error) {
		fillItem.SetStartPos(reader)
		if fillItem.Raw, err = reader.ReadBytes(8); err == nil {
			fillItem.SetValue(int64(endianConverter.Uint64(fillItem.Raw)))
		}
		fillItem.SetEndPos(reader)
		return
	}
}

func BuildReadS4(endianConverter EndianReader) ReadTo {
	return func(fillItem *Item, reader *Reader) (err error) {
		fillItem.SetStartPos(reader)
		if fillItem.Raw, err = reader.ReadBytes(4); err == nil {
			fillItem.SetValue(int32(endianConverter.Uint32(fillItem.Raw)))
		}
		fillItem.SetEndPos(reader)
		return
	}
}

func BuildReadS2(endianConverter EndianReader) ReadTo {
	return func(fillItem *Item, reader *Reader) (err error) {
		fillItem.SetStartPos(reader)
		if fillItem.Raw, err = reader.ReadBytes(2); err == nil {
			fillItem.SetValue(int16(endianConverter.Uint16(fillItem.Raw)))
		}
		fillItem.SetEndPos(reader)
		return
	}
}

func BuildReadS1() ReadTo {
	return func(fillItem *Item, reader *Reader) (err error) {
		fillItem.SetStartPos(reader)
		if fillItem.Raw, err = reader.ReadBytes(1); err == nil {
			fillItem.SetValue(int8(fillItem.Raw[0]))
		}
		fillItem.SetEndPos(reader)
		return
	}
}
