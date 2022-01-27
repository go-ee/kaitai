package kaitai

func ReadB(endianConverter EndianReader, length uint8) (ret ReadFix) {
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

func ReadB1(endianConverter EndianReader) ReadFix {
	return func(reader *Reader) (ret interface{}, err error) {
		var value uint64
		if value, err = endianConverter.ReadBitsInt(reader, 1); err == nil {
			ret = value != 0
		}
		return
	}
}

func ReadB2(endianConverter EndianReader) ReadFix {
	return func(reader *Reader) (ret interface{}, err error) {
		var value uint64
		if value, err = endianConverter.ReadBitsInt(reader, 2); err == nil {
			ret = uint(value)
		}
		return
	}
}

func ReadBUint64(endianConverter EndianReader, length uint8) ReadFix {
	return func(reader *Reader) (ret interface{}, err error) {
		var value uint64
		if value, err = endianConverter.ReadBitsInt(reader, length); err == nil {
			ret = value
		}
		return
	}
}
