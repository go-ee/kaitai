package kaitai

import (
	"github.com/sirupsen/logrus"
)

var BigEndianBuildReadU = &EndianBuildReadU{BigEndianConverter}
var LittleEndianBuildReadU = &EndianBuildReadU{LittleEndianConverter}

type EndianBuildReadU struct {
	endianConverter EndianReader
}

func (o *EndianBuildReadU) BuildRead(length uint8) (ret ParentRead) {
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

func (o *EndianBuildReadU) BuildRead1() ParentRead {
	return func(parent *Item, reader *ReaderIO) (ret interface{}, err error) {
		var data []byte
		if data, err = reader.ReadBytes(1); err == nil {
			ret = o.endianConverter.Uint8(data)
		}
		return
	}
}

func (o *EndianBuildReadU) BuildRead2() ParentRead {
	return func(parent *Item, reader *ReaderIO) (ret interface{}, err error) {
		var data []byte
		if data, err = reader.ReadBytes(2); err == nil {
			ret = o.endianConverter.Uint16(data)
		}
		return
	}
}

func (o *EndianBuildReadU) BuildRead4() ParentRead {
	return func(parent *Item, reader *ReaderIO) (ret interface{}, err error) {
		var data []byte
		if data, err = reader.ReadBytes(4); err == nil {
			ret = o.endianConverter.Uint32(data)
		}
		return
	}
}

func (o *EndianBuildReadU) BuildRead8() ParentRead {
	return func(parent *Item, reader *ReaderIO) (ret interface{}, err error) {
		var data []byte
		if data, err = reader.ReadBytes(8); err == nil {
			ret = o.endianConverter.Uint64(data)
		}
		return
	}
}
