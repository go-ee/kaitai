package kaitai

func ReadB(endianConverter EndianReader, length uint8) (ret ReadTo) {
	switch length {
	case 1:
		ret = ReadB1(endianConverter)
	case 2:
		ret = ReadB2(endianConverter)
	default:
		ret = ReadBUint64(endianConverter, length)
	}
	return
}

func ReadB1(endianConverter EndianReader) ReadTo {
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

func ReadB2(endianConverter EndianReader) ReadTo {
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

func ReadBUint64(endianConverter EndianReader, length uint8) ReadTo {
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
