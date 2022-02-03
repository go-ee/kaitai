package kaitai

import (
	"github.com/sirupsen/logrus"
)

var BigEndianBuildReadU = &EndianBuildReadU{BigEndianConverter}
var LittleEndianBuildReadU = &EndianBuildReadU{LittleEndianConverter}
var BigEndianBuildLazyReadU = &EndianBuildLazyReadU{BigEndianConverter}
var LittleEndianBuildLazyReadU = &EndianBuildLazyReadU{LittleEndianConverter}

type EndianBuildReadU struct {
	endianConverter EndianReader
}

func (o *EndianBuildReadU) BuildRead(length uint8) (ret ReadTo) {
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

func (o *EndianBuildReadU) BuildRead1() ReadTo {
	return func(fillItem *Item, reader *ReaderIO) (err error) {
		if fillItem.Raw, err = reader.ReadBytes(1); err == nil {
			fillItem.SetValue(o.endianConverter.Uint8(fillItem.Raw))
		}
		return
	}
}

func (o *EndianBuildReadU) BuildRead2() ReadTo {
	return func(fillItem *Item, reader *ReaderIO) (err error) {
		if fillItem.Raw, err = reader.ReadBytes(2); err == nil {
			fillItem.SetValue(o.endianConverter.Uint16(fillItem.Raw))
		}
		return
	}
}

func (o *EndianBuildReadU) BuildRead4() ReadTo {
	return func(fillItem *Item, reader *ReaderIO) (err error) {
		if fillItem.Raw, err = reader.ReadBytes(4); err == nil {
			fillItem.SetValue(o.endianConverter.Uint32(fillItem.Raw))
		}
		return
	}
}

func (o *EndianBuildReadU) BuildRead8() ReadTo {
	return func(fillItem *Item, reader *ReaderIO) (err error) {
		if fillItem.Raw, err = reader.ReadBytes(8); err == nil {
			fillItem.SetValue(o.endianConverter.Uint64(fillItem.Raw))
		}
		return
	}
}

type EndianBuildLazyReadU struct {
	endianConverter EndianReader
}

func (o *EndianBuildLazyReadU) BuildRead(length uint8) (ret ReadTo) {
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

func (o *EndianBuildLazyReadU) BuildRead1() ReadTo {
	return func(fillItem *Item, reader *ReaderIO) (err error) {
		if fillItem.Raw, err = reader.ReadBytes(1); err == nil {
			fillItem.Decode = o.DecodeUint8
		}
		return
	}
}

func (o *EndianBuildLazyReadU) BuildRead2() ReadTo {
	return func(fillItem *Item, reader *ReaderIO) (err error) {
		if fillItem.Raw, err = reader.ReadBytes(2); err == nil {
			fillItem.Decode = o.DecodeUint16
		}
		return
	}
}

func (o *EndianBuildLazyReadU) BuildRead4() ReadTo {
	return func(fillItem *Item, reader *ReaderIO) (err error) {
		if fillItem.Raw, err = reader.ReadBytes(4); err == nil {
			fillItem.Decode = o.DecodeUint32
		}
		return
	}
}

func (o *EndianBuildLazyReadU) BuildRead8() ReadTo {
	return func(fillItem *Item, reader *ReaderIO) (err error) {
		if fillItem.Raw, err = reader.ReadBytes(8); err == nil {
			fillItem.Decode = o.DecodeUint64
		}
		return
	}
}

func (o *EndianBuildLazyReadU) DecodeUint8(fillItem *Item) {
	fillItem.SetValue(fillItem.Raw[0])
}

func (o *EndianBuildLazyReadU) DecodeUint16(fillItem *Item) {
	fillItem.SetValue(o.endianConverter.Uint16(fillItem.Raw))
}

func (o *EndianBuildLazyReadU) DecodeUint32(fillItem *Item) {
	fillItem.SetValue(o.endianConverter.Uint32(fillItem.Raw))
}

func (o *EndianBuildLazyReadU) DecodeUint64(fillItem *Item) {
	fillItem.SetValue(o.endianConverter.Uint64(fillItem.Raw))
}
