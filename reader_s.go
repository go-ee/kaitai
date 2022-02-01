package kaitai

import (
	"github.com/sirupsen/logrus"
)

var BigEndianLazyBuildReadS = &EndianLazyBuildReadS{BigEndianConverter}
var LittleEndianBuildLazyReadS = &EndianLazyBuildReadS{LittleEndianConverter}

type EndianLazyBuildReadS struct {
	endianConverter EndianReader
}

func (o *EndianLazyBuildReadS) BuildRead(length uint8) (ret ReadTo) {
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
		logrus.Infof("not supported Native(s,%v)", length)
	}
	return
}

func (o *EndianLazyBuildReadS) BuildRead8() ReadTo {
	return func(fillItem *Item, reader *ReaderIO) (err error) {
		if fillItem.Raw, err = reader.ReadBytes(8); err == nil {
			fillItem.SetValue(int64(o.endianConverter.Uint64(fillItem.Raw)))
		}
		return
	}
}

func (o *EndianLazyBuildReadS) BuildRead4() ReadTo {
	return func(fillItem *Item, reader *ReaderIO) (err error) {
		if fillItem.Raw, err = reader.ReadBytes(4); err == nil {
			fillItem.SetValue(int32(o.endianConverter.Uint32(fillItem.Raw)))
		}
		return
	}
}

func (o *EndianLazyBuildReadS) BuildRead2() ReadTo {
	return func(fillItem *Item, reader *ReaderIO) (err error) {
		if fillItem.Raw, err = reader.ReadBytes(2); err == nil {
			fillItem.SetValue(int16(o.endianConverter.Uint16(fillItem.Raw)))
		}
		return
	}
}

func (o *EndianLazyBuildReadS) BuildRead1() ReadTo {
	return func(fillItem *Item, reader *ReaderIO) (err error) {
		if fillItem.Raw, err = reader.ReadBytes(1); err == nil {
			fillItem.SetValue(int8(fillItem.Raw[0]))
		}
		return
	}
}
