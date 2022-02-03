package kaitai

import (
	"fmt"
	"strconv"
)

func BuildReadAttr(attr *Attr, parse Parse) (ret ReadTo) {
	if attr.SizeEos == "true" {
		ret = BuildReadToFull(parse)
	} else if attr.Size != "" {
		if length, err := strconv.Atoi(attr.Size); err == nil {
			ret = BuildReadToLength(uint16(length), parse)
		} else {
			ret = BuildReadToLengthExpr(attr.Size, parse)
		}
	}
	return
}

func BuildReadToFull(parse Parse) (ret ReadTo) {
	return func(fillItem *Item, reader *ReaderIO) (err error) {
		if fillItem.Raw, err = reader.ReadBytesFull(); err == nil {
			fillItem.value, err = parse(fillItem.Raw)
		}
		return
	}
}

func BuildReadToLength(length uint16, parse Parse) (ret ReadTo) {
	return func(fillItem *Item, reader *ReaderIO) (err error) {
		return ReadToLength(fillItem, reader, length, parse)
	}
}

func ReadToLength(fillItem *Item, reader *ReaderIO, length uint16, parse Parse) (err error) {
	if length > 0 {
		fillItem.Raw, err = reader.ReadBytes(length)
	} else {
		fillItem.Raw, err = reader.ReadBytesFull()
	}

	if err == nil {
		fillItem.value, err = parse(fillItem.Raw)
	}
	return
}

func BuildReadToLengthExpr(expr string, parse Parse) (ret ReadTo) {
	return func(fillItem *Item, reader *ReaderIO) (err error) {
		var sizeItem *Item
		if sizeItem, err = fillItem.Parent.Expr(expr); err == nil {
			var length uint16
			if length, err = toUint16(sizeItem.Value()); err == nil {
				return ReadToLength(fillItem, reader, length, parse)
			} else {
				err = fmt.Errorf("cant parse Size to uint16, expr=%v, valiue=%v, %v", expr, sizeItem.Value(), err)
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

func BuildLazyReadAttr(attr *Attr, decode Decode) (ret ReadTo) {
	if attr.SizeEos == "true" {
		ret = BuildLazyReadToFull(decode)
	} else if attr.Size != "" {
		if length, err := strconv.Atoi(attr.Size); err == nil {
			ret = BuildLazyReadToLength(uint16(length), decode)
		} else {
			ret = BuildLazyReadToLengthExpr(attr.Size, decode)
		}
	}
	return
}

func BuildLazyReadToFull(decode Decode) (ret ReadTo) {
	return func(fillItem *Item, reader *ReaderIO) (err error) {
		if fillItem.Raw, err = reader.ReadBytesFull(); err == nil {
			fillItem.Decode = decode
		}
		return
	}
}

func BuildLazyReadToLength(length uint16, decode Decode) (ret ReadTo) {
	return func(fillItem *Item, reader *ReaderIO) (err error) {
		return LazyReadToLength(fillItem, reader, length, decode)
	}
}

func LazyReadToLength(fillItem *Item, reader *ReaderIO, length uint16, decode Decode) (err error) {
	if length > 0 {
		fillItem.Raw, err = reader.ReadBytes(length)
	} else {
		fillItem.Raw, err = reader.ReadBytesFull()
	}

	if err == nil {
		fillItem.Decode = decode
	}
	return
}

func BuildLazyReadToLengthExpr(expr string, decode Decode) (ret ReadTo) {
	return func(fillItem *Item, reader *ReaderIO) (err error) {
		var sizeItem *Item
		if sizeItem, err = fillItem.Parent.Expr(expr); err == nil {
			var length uint16
			if length, err = toUint16(sizeItem.Value()); err == nil {
				return LazyReadToLength(fillItem, reader, length, decode)
			} else {
				err = fmt.Errorf("cant decode Size to uint16, expr=%v, valiue=%v, %v", expr, sizeItem.Value(), err)
			}
		}
		return
	}
}

func DecodeToString(fillItem *Item) {
	fillItem.value = string(fillItem.Raw)
}

func DecodeToSame(fillItem *Item) {
	fillItem.value = fillItem.Raw
}
