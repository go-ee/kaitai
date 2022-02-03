package kaitai

import (
	"github.com/sirupsen/logrus"
)

var BigEndianBuildReadS = &EndianBuildReadS{BigEndianConverter}
var LittleEndianBuildReadS = &EndianBuildReadS{LittleEndianConverter}
var BigEndianBuildLazyReadS = &EndianBuildLazyReadS{BigEndianConverter}
var LittleEndianBuildLazyReadS = &EndianBuildLazyReadS{LittleEndianConverter}

type EndianBuildReadS struct {
	endianConverter EndianReader
}

func (o *EndianBuildReadS) BuildRead(length uint8) (ret ReadTo) {
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

func (o *EndianBuildReadS) BuildRead8() ReadTo {
	return func(fillItem *Item, reader *ReaderIO) (err error) {
		if fillItem.Raw, err = reader.ReadBytes(8); err == nil {
			fillItem.SetValue(int64(o.endianConverter.Uint64(fillItem.Raw)))
		}
		return
	}
}

func (o *EndianBuildReadS) BuildRead4() ReadTo {
	return func(fillItem *Item, reader *ReaderIO) (err error) {
		if fillItem.Raw, err = reader.ReadBytes(4); err == nil {
			fillItem.SetValue(int64(o.endianConverter.Uint32(fillItem.Raw)))
		}
		return
	}
}

func (o *EndianBuildReadS) BuildRead2() ReadTo {
	return func(fillItem *Item, reader *ReaderIO) (err error) {
		if fillItem.Raw, err = reader.ReadBytes(2); err == nil {
			fillItem.SetValue(int64(o.endianConverter.Uint16(fillItem.Raw)))
		}
		return
	}
}

func (o *EndianBuildReadS) BuildRead1() ReadTo {
	return func(fillItem *Item, reader *ReaderIO) (err error) {
		if fillItem.Raw, err = reader.ReadBytes(1); err == nil {
			fillItem.SetValue(int64(o.endianConverter.Uint8(fillItem.Raw)))
		}
		return
	}
}

type EndianBuildLazyReadS struct {
	endianConverter EndianReader
}

func (o *EndianBuildLazyReadS) BuildRead(length uint8) (ret ReadTo) {
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

func (o *EndianBuildLazyReadS) BuildRead8() ReadTo {
	return func(fillItem *Item, reader *ReaderIO) (err error) {
		if fillItem.Raw, err = reader.ReadBytes(8); err == nil {
			fillItem.Decode = o.DecodeInt64
		}
		return
	}
}

func (o *EndianBuildLazyReadS) BuildRead4() ReadTo {
	return func(fillItem *Item, reader *ReaderIO) (err error) {
		if fillItem.Raw, err = reader.ReadBytes(4); err == nil {
			fillItem.Decode = o.DecodeInt32
		}
		return
	}
}

func (o *EndianBuildLazyReadS) BuildRead2() ReadTo {
	return func(fillItem *Item, reader *ReaderIO) (err error) {
		if fillItem.Raw, err = reader.ReadBytes(2); err == nil {
			fillItem.Decode = o.DecodeInt16
		}
		return
	}
}

func (o *EndianBuildLazyReadS) BuildRead1() ReadTo {
	return func(fillItem *Item, reader *ReaderIO) (err error) {
		if fillItem.Raw, err = reader.ReadBytes(1); err == nil {
			fillItem.Decode = o.DecodeInt8
		}
		return
	}
}

func (o *EndianBuildLazyReadS) DecodeInt64(fillItem *Item) {
	fillItem.SetValue(int64(o.endianConverter.Uint64(fillItem.Raw)))
}

func (o *EndianBuildLazyReadS) DecodeInt32(fillItem *Item) {
	fillItem.SetValue(int32(o.endianConverter.Uint32(fillItem.Raw)))
}

func (o *EndianBuildLazyReadS) DecodeInt16(fillItem *Item) {
	fillItem.SetValue(int16(o.endianConverter.Uint16(fillItem.Raw)))
}

func (o *EndianBuildLazyReadS) DecodeInt8(fillItem *Item) {
	fillItem.SetValue(int8(fillItem.Raw[0]))
}
