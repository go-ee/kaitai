package kaitai

import (
	"github.com/sirupsen/logrus"
)

var BigEndianBuildReadS = &EndianBuildReadS{BigEndianConverter}
var LittleEndianBuildReadS = &EndianBuildReadS{LittleEndianConverter}

type EndianBuildReadS struct {
	endianConverter EndianReader
}

func (o *EndianBuildReadS) BuildRead(length uint8) (ret Read) {
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

func (o *EndianBuildReadS) BuildRead8() Read {
	return func(reader *ReaderIO) (ret interface{}, err error) {
		var data []byte
		if data, err = reader.ReadBytes(8); err == nil {
			ret = int64(o.endianConverter.Uint64(data))
		}
		return
	}
}

func (o *EndianBuildReadS) BuildRead4() Read {
	return func(reader *ReaderIO) (ret interface{}, err error) {
		var data []byte
		if data, err = reader.ReadBytes(4); err == nil {
			ret = int64(o.endianConverter.Uint32(data))
		}
		return
	}
}

func (o *EndianBuildReadS) BuildRead2() Read {
	return func(reader *ReaderIO) (ret interface{}, err error) {
		var data []byte
		if data, err = reader.ReadBytes(2); err == nil {
			ret = int64(o.endianConverter.Uint16(data))
		}
		return
	}
}

func (o *EndianBuildReadS) BuildRead1() Read {
	return func(reader *ReaderIO) (ret interface{}, err error) {
		var data []byte
		if data, err = reader.ReadBytes(1); err == nil {
			ret = int64(o.endianConverter.Uint8(data))
		}
		return
	}
}
