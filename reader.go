package kaitai

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
)

type Reader interface {
	NewItem(parent *Item) *Item
	ReadTo(fillItem *Item, reader *ReaderIO) (err error)
}

type ReadTo func(fillItem *Item, reader *ReaderIO) (err error)
type Parse func(data []byte) (interface{}, error)
type Decode func(fillItem *Item)

type Item struct {
	value    interface{}
	Err      error
	Parent   *Item
	StartPos int64
	EndPos   int64
	Raw      []byte
	Attr     *Attr
	Accessor interface{}
	Decode   Decode
}

func (o *Item) Value() interface{} {
	if o.value == nil && o.Err == nil && o.Raw != nil {
		o.Decode(o)
	}
	return o.value
}

func (o *Item) Expr(expr string) (ret *Item, err error) {
	ret = o
	for _, part := range strings.Split(expr, ".") {
		if value, ok := ret.Value().(map[string]*Item); ok {
			if ret = value[part]; ret == nil {
				err = fmt.Errorf("can't resolve '%v' of expression '%v'", part, expr)
				break
			}
		} else {
			err = fmt.Errorf("can't resolve '%v' of expression '%v'", part, expr)
			break
		}
	}
	return
}

func (o *Item) SetStartPos(reader *ReaderIO) {
	o.StartPos = reader.Position()
}

func (o *Item) SetEndPos(reader *ReaderIO) {
	o.EndPos = reader.Position()
}

func (o *Item) SetValue(value interface{}) {
	o.value = value
}

type ReaderIO struct {
	io.ReadSeeker

	offset int64
	buf    [8]byte

	bitsLeft uint8
	bits     uint64
}

func (o *ReaderIO) ReadBytes(n uint16) (ret []byte, err error) {
	if n < 0 {
		err = fmt.Errorf("ReadBytes(%d): negative number of bytes to read", n)
		return
	}

	ret = make([]byte, n)
	_, err = io.ReadFull(o, ret)
	return
}

func (o *ReaderIO) ReadBytesFull() ([]byte, error) {
	return ioutil.ReadAll(o)
}

func (o *ReaderIO) ReadBytesAsReader(n uint16) (ret *ReaderIO, raw []byte, err error) {
	currentPos := o.Position()
	if raw, err = o.ReadBytes(n); err == nil {
		ret = &ReaderIO{ReadSeeker: bytes.NewReader(raw), offset: currentPos}
	}
	return
}

func (o *ReaderIO) Position() (ret int64) {
	ret, _ = o.Seek(0, io.SeekCurrent)
	return o.offset + ret
}

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
		var data []byte
		if data, err = reader.ReadBytesFull(); err == nil {
			fillItem.value, err = parse(data)
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
	var data []byte
	if length > 0 {
		data, err = reader.ReadBytes(length)
	} else {
		data, err = reader.ReadBytesFull()
	}

	if err == nil {
		fillItem.value, err = parse(data)
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
