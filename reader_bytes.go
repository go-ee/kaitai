package kaitai

import (
	"fmt"
	"strconv"
)

func BuildReadAttr(attr *Attr, parse Parse) (ret Reader) {
	if attr.SizeEos == "true" {
		ret = &AttrParentRead{attr, BuildReadToFull(parse)}
	} else if attr.Size != "" {
		if length, err := strconv.Atoi(attr.Size); err == nil {
			ret = &AttrParentRead{attr, BuildReadToLength(uint16(length), parse)}
		} else {
			ret = &AttrParentRead{attr, BuildReadToLengthExpr(attr.Size, parse)}
		}
	}
	return
}

func BuildReadToFull(parse Parse) (ret ParentRead) {
	return func(parent Item, reader *ReaderIO) (ret interface{}, err error) {
		var data []byte
		if data, err = reader.ReadBytesFull(); err == nil {
			ret, err = parse(data)
		}
		return
	}
}

func BuildReadToLength(length uint16, parse Parse) (ret ParentRead) {
	return func(parent Item, reader *ReaderIO) (ret interface{}, err error) {
		return ReadToLength(reader, length, parse)
	}
}

func ReadToLength(reader *ReaderIO, length uint16, parse Parse) (ret interface{}, err error) {
	var data []byte
	if length > 0 {
		data, err = reader.ReadBytes(length)
	} else {
		data, err = reader.ReadBytesFull()
	}

	if err == nil {
		ret, err = parse(data)
	}
	return
}

func BuildReadToLengthExpr(expr string, parse Parse) (ret ParentRead) {
	return func(parent Item, reader *ReaderIO) (ret interface{}, err error) {
		var sizeItem interface{}
		if sizeItem, err = parent.ExprValue(expr); err == nil {
			var length uint16
			if length, err = ToUint16(sizeItem); err == nil {
				return ReadToLength(reader, length, parse)
			} else {
				err = fmt.Errorf("cant parse Size to uint16, expr=%v, valiue=%v, %v", expr, sizeItem, err)
			}
		}
		return
	}
}

func ToString(data []byte) (interface{}, error) {
	return string(data), nil
}

func ToSame(data []byte) (interface{}, error) {
	return data, nil
}
