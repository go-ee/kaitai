package kaitai

import (
	"fmt"
	"log"
)

var BigEndianBuildReadF = &EndianBuildReadF{BigEndianConverter}
var LittleEndianBuildReadF = &EndianBuildReadF{LittleEndianConverter}
var BigEndianBuildLazyReadF = &EndianBuildLazyReadF{BigEndianConverter}
var LittleEndianBuildLazyReadF = &EndianBuildLazyReadF{LittleEndianConverter}

type EndianBuildReadF struct {
	endianConverter EndianReader
}

func (o *EndianBuildReadF) BuildRead(length uint8) (ret ReadTo) {
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

func (o *EndianBuildReadF) BuildRead4() ReadTo {
	return func(fillItem *Item, reader *ReaderIO) (err error) {
		if fillItem.Raw, err = reader.ReadBytes(4); err == nil {
			fillItem.SetValue(o.endianConverter.Float32fromBits(fillItem.Raw))
		}
		return
	}
}

func (o *EndianBuildReadF) BuildRead8() ReadTo {
	return func(fillItem *Item, reader *ReaderIO) (err error) {
		if fillItem.Raw, err = reader.ReadBytes(8); err == nil {
			fillItem.SetValue(o.endianConverter.Float64fromBits(fillItem.Raw))
		}
		return
	}
}

type EndianBuildLazyReadF struct {
	endianConverter EndianReader
}

func (o *EndianBuildLazyReadF) BuildRead(length uint8) (ret ReadTo) {
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

func (o *EndianBuildLazyReadF) BuildRead4() ReadTo {
	return func(fillItem *Item, reader *ReaderIO) (err error) {
		if fillItem.Raw, err = reader.ReadBytes(4); err == nil {
			fillItem.Decode = o.DecodeFloat32fromBits
		}
		return
	}
}

func (o *EndianBuildLazyReadF) BuildRead8() ReadTo {
	return func(fillItem *Item, reader *ReaderIO) (err error) {
		if fillItem.Raw, err = reader.ReadBytes(8); err == nil {
			fillItem.Decode = o.DecodeFloat64fromBits
		}
		return
	}
}

func (o *EndianBuildLazyReadF) DecodeFloat32fromBits(fillItem *Item) {
	fillItem.SetValue(o.endianConverter.Float32fromBits(fillItem.Raw))
}

func (o *EndianBuildLazyReadF) DecodeFloat64fromBits(fillItem *Item) {
	fillItem.SetValue(o.endianConverter.Float64fromBits(fillItem.Raw))
}
