package kaitai

import (
	"github.com/sirupsen/logrus"
)

var BigEndianLazyBuildReadU = &EndianLazyBuildReadU{BigEndianConverter}
var LittleEndianBuildLazyReadU = &EndianLazyBuildReadU{LittleEndianConverter}

type EndianLazyBuildReadU struct {
	endianConverter EndianReader
}

func (o *EndianLazyBuildReadU) BuildRead(length uint8) (ret ReadTo) {
	switch length {
	case 1:
		ret = o.BuildRead1()
	case 2:
		ret = o.BuildRead2()
	case 4:
		ret = o.BuildRead4()
	case 8:
		ret = o.BuildRead8()
	default:
		logrus.Infof("not supported Native(u,%v)", length)
	}
	return
}

func (o *EndianLazyBuildReadU) BuildRead1() ReadTo {
	return func(fillItem *Item, reader *ReaderIO) (err error) {
		if fillItem.Raw, err = reader.ReadBytes(1); err == nil {
			fillItem.SetValue(fillItem.Raw[0])
		}
		return
	}
}

func (o *EndianLazyBuildReadU) BuildRead2() ReadTo {
	return func(fillItem *Item, reader *ReaderIO) (err error) {
		if fillItem.Raw, err = reader.ReadBytes(2); err == nil {
			fillItem.SetValue(o.endianConverter.Uint16(fillItem.Raw))
		}
		return
	}
}

func (o *EndianLazyBuildReadU) BuildRead4() ReadTo {
	return func(fillItem *Item, reader *ReaderIO) (err error) {
		if fillItem.Raw, err = reader.ReadBytes(4); err == nil {
			fillItem.SetValue(o.endianConverter.Uint32(fillItem.Raw))
		}
		return
	}
}

func (o *EndianLazyBuildReadU) BuildRead8() ReadTo {
	return func(fillItem *Item, reader *ReaderIO) (err error) {
		if fillItem.Raw, err = reader.ReadBytes(8); err == nil {
			fillItem.SetValue(o.endianConverter.Uint64(fillItem.Raw))
		}
		return
	}
}
