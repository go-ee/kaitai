package kaitai

import (
	"github.com/sirupsen/logrus"
)

func BuildReadU(endianConverter EndianReader, length uint8) (ret ReadTo) {
	switch length {
	case 1:
		ret = BuildReadU1()
	case 2:
		ret = BuildReadU2(endianConverter)
	case 4:
		ret = BuildReadU4(endianConverter)
	case 8:
		ret = BuildReadU8(endianConverter)
	default:
		logrus.Infof("not supported Native(u,%v)", length)
	}
	return
}

func BuildReadU1() ReadTo {
	return func(fillItem *Item, reader *Reader) (err error) {
		if fillItem.Raw, err = reader.ReadBytes(1); err == nil {
			fillItem.SetValue(fillItem.Raw[0])
		}
		return
	}
}

func BuildReadU2(endianConverter EndianReader) ReadTo {
	return func(fillItem *Item, reader *Reader) (err error) {
		if fillItem.Raw, err = reader.ReadBytes(2); err == nil {
			fillItem.SetValue(endianConverter.Uint16(fillItem.Raw))
		}
		return
	}
}

func BuildReadU4(endianConverter EndianReader) ReadTo {
	return func(fillItem *Item, reader *Reader) (err error) {
		if fillItem.Raw, err = reader.ReadBytes(4); err == nil {
			fillItem.SetValue(endianConverter.Uint32(fillItem.Raw))
		}
		return
	}
}

func BuildReadU8(endianConverter EndianReader) ReadTo {
	return func(fillItem *Item, reader *Reader) (err error) {
		if fillItem.Raw, err = reader.ReadBytes(8); err == nil {
			fillItem.SetValue(endianConverter.Uint64(fillItem.Raw))
		}
		return
	}
}
