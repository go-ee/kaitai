package kaitai

import (
	"fmt"
	"log"
)

var BigEndianBuildReadF = &EndianBuildReadF{BigEndianConverter}
var LittleEndianBuildReadF = &EndianBuildReadF{LittleEndianConverter}

type EndianBuildReadF struct {
	endianConverter EndianReader
}

func (o *EndianBuildReadF) BuildRead(length uint8) (ret ParentRead) {
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

func (o *EndianBuildReadF) BuildRead4() ParentRead {
	return func(parent Item, reader *ReaderIO) (ret interface{}, err error) {
		var data []byte
		if data, err = reader.ReadBytes(4); err == nil {
			ret = o.endianConverter.Float32fromBits(data)
		}
		return
	}
}

func (o *EndianBuildReadF) BuildRead8() ParentRead {
	return func(parent Item, reader *ReaderIO) (ret interface{}, err error) {
		var data []byte
		if data, err = reader.ReadBytes(8); err == nil {
			ret = o.endianConverter.Float64fromBits(data)
		}
		return
	}
}
