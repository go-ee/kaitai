package kaitai

var BigEndianLazyBuildReadB = &EndianLazyBuildReadB{BigEndianConverter}
var LittleEndianLazyBuildReadB = &EndianLazyBuildReadB{LittleEndianConverter}

type EndianLazyBuildReadB struct {
	endianConverter EndianReader
}

func (o *EndianLazyBuildReadB) BuildRead(length uint8) (ret ReadTo) {
	switch length {
	case 1:
		ret = o.BuildRead1()
	case 2:
		ret = o.BuildRead2()
	default:
		ret = o.BuildReadUint64(length)
	}
	return
}

func (o *EndianLazyBuildReadB) BuildRead1() ReadTo {
	return func(fillItem *Item, reader *ReaderIO) (err error) {
		var value uint64
		if value, err = o.endianConverter.ReadBitsInt(reader, 1); err == nil {
			fillItem.SetValue(value != 0)
		}
		return
	}
}

func (o *EndianLazyBuildReadB) BuildRead2() ReadTo {
	return func(fillItem *Item, reader *ReaderIO) (err error) {
		var value uint64
		if value, err = o.endianConverter.ReadBitsInt(reader, 2); err == nil {
			fillItem.SetValue(uint(value))
		}
		return
	}
}

func (o *EndianLazyBuildReadB) BuildReadUint64(length uint8) ReadTo {
	return func(fillItem *Item, reader *ReaderIO) (err error) {
		var value uint64
		if value, err = o.endianConverter.ReadBitsInt(reader, length); err == nil {
			fillItem.SetValue(value)
		}
		return
	}
}
