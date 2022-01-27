package kaitai

func BuildReadB(endianConverter EndianReader, length uint8) (ret ReadTo) {
	switch length {
	case 1:
		ret = BuildReadB1(endianConverter)
	case 2:
		ret = BuildReadB2(endianConverter)
	default:
		ret = BuildReadBUint64(endianConverter, length)
	}
	return
}

func BuildReadB1(endianConverter EndianReader) ReadTo {
	return func(fillItem *Item, reader *Reader) (err error) {
		fillItem.SetStartPos(reader)
		var value uint64
		if value, err = endianConverter.ReadBitsInt(reader, 1); err == nil {
			fillItem.SetValue(value != 0)
		}
		fillItem.SetEndPos(reader)
		return
	}
}

func BuildReadB2(endianConverter EndianReader) ReadTo {
	return func(fillItem *Item, reader *Reader) (err error) {
		fillItem.SetStartPos(reader)
		var value uint64
		if value, err = endianConverter.ReadBitsInt(reader, 2); err == nil {
			fillItem.SetValue(uint(value))
		}
		fillItem.SetEndPos(reader)
		return
	}
}

func BuildReadBUint64(endianConverter EndianReader, length uint8) ReadTo {
	return func(fillItem *Item, reader *Reader) (err error) {
		fillItem.SetStartPos(reader)
		var value uint64
		if value, err = endianConverter.ReadBitsInt(reader, length); err == nil {
			fillItem.SetValue(value)
		}
		fillItem.SetEndPos(reader)
		return
	}
}
