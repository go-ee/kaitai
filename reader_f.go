package kaitai

import (
	"fmt"
	"log"
)

var BigEndianLazyBuildReadF = &EndianLazyBuildReadF{BigEndianConverter}
var LittleEndianBuildLazyReadF = &EndianLazyBuildReadF{LittleEndianConverter}

type EndianLazyBuildReadF struct {
	endianConverter EndianReader
}

func (o *EndianLazyBuildReadF) BuildRead(length uint8) (ret ReadTo) {
	switch length {
	case 4:
		ret = o.BuildRead4()
	case 8:
		ret = o.BuildRead8()
	default:
		log.Println(fmt.Sprintf("not supported Native(f,%v)", length))
	}
	return
}

func (o *EndianLazyBuildReadF) BuildRead4() ReadTo {
	return func(fillItem *Item, reader *ReaderIO) (err error) {
		if fillItem.Raw, err = reader.ReadBytes(4); err == nil {
			fillItem.SetValue(o.endianConverter.Float32fromBits(fillItem.Raw))
		}
		return
	}
}

func (o *EndianLazyBuildReadF) BuildRead8() ReadTo {
	return func(fillItem *Item, reader *ReaderIO) (err error) {
		if fillItem.Raw, err = reader.ReadBytes(8); err == nil {
			fillItem.SetValue(o.endianConverter.Float64fromBits(fillItem.Raw))
		}
		return
	}
}
