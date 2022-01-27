package kaitai

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
)

type AttrReader interface {
	NewItem(parent *Item) *Item
	ReadTo(fillItem *Item, reader *Reader) (err error)
}

type ReadTo func(fillItem *Item, reader *Reader) (err error)
type Parse func(data []byte) (interface{}, error)
type Decode func(fillItem *Item)

type ReadToReader struct {
	attr     *Attr
	accessor interface{}
	readTo   ReadTo
}

func (o *ReadToReader) ReadTo(fillItem *Item, reader *Reader) (err error) {
	return o.readTo(fillItem, reader)
}

func (o *ReadToReader) Attr() *Attr {
	return o.attr
}

func (o *ReadToReader) Accessor() interface{} {
	return o.accessor
}

func (o *ReadToReader) NewItem(parent *Item) *Item {
	return &Item{Attr: o.attr, Accessor: o.accessor, Parent: parent}
}

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

func (o *Item) SetStartPos(reader *Reader) {
	o.StartPos = reader.Position()
}

func (o *Item) SetEndPos(reader *Reader) {
	o.EndPos = reader.Position()
}

func (o *Item) SetValue(value interface{}) {
	o.value = value
}

type Reader struct {
	io.ReadSeeker

	offset int64
	buf    [8]byte

	bitsLeft uint8
	bits     uint64
}

func (o *Reader) ReadBytes(n uint16) (ret []byte, err error) {
	if n < 0 {
		err = fmt.Errorf("ReadBytes(%d): negative number of bytes to read", n)
		return
	}

	ret = make([]byte, n)
	_, err = io.ReadFull(o, ret)
	return
}

func (o *Reader) ReadBytesFull() ([]byte, error) {
	return ioutil.ReadAll(o)
}

func (o *Reader) ReadBytesAsReader(n uint16) (ret *Reader, raw []byte, err error) {
	currentPos := o.Position()
	if raw, err = o.ReadBytes(n); err == nil {
		ret = &Reader{ReadSeeker: bytes.NewReader(raw), offset: currentPos}
	}
	return
}

func (o *Reader) Position() (ret int64) {
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
			ret = ReadToLengthExpr(attr.Size, parse)
		}
	}
	return
}

type Options struct {
	LazyDecoding bool
	RawFill      bool
	PositionFill bool
}

func BuildReadToFull(parse Parse) (ret ReadTo) {
	return func(fillItem *Item, reader *Reader) (err error) {
		var data []byte
		if data, err = reader.ReadBytesFull(); err == nil {
			fillItem.value, err = parse(data)
		}
		return
	}
}

func BuildReadToLength(length uint16, parse Parse) (ret ReadTo) {
	return func(fillItem *Item, reader *Reader) (err error) {
		return ReadToLength(fillItem, reader, length, parse)
	}
}

func ReadToLength(fillItem *Item, reader *Reader, length uint16, parse Parse) (err error) {
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

func ReadToLengthExpr(expr string, parse Parse) (ret ReadTo) {
	return func(fillItem *Item, reader *Reader) (err error) {
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

func toUint16(value interface{}) (ret uint16, err error) {
	if number, ok := value.(int); ok {
		ret = uint16(number)
		return
	}

	if number, ok := value.(uint); ok {
		ret = uint16(number)
		return
	}

	if number, ok := value.(uint8); ok {
		ret = uint16(number)
		return
	}

	if number, ok := value.(uint16); ok {
		ret = number
		return
	}

	if number, ok := value.(uint32); ok {
		ret = uint16(number)
		return
	}

	if number, ok := value.(uint64); ok {
		ret = uint16(number)
		return
	}

	str := fmt.Sprintf("%v", value)
	var intValue int
	if intValue, err = strconv.Atoi(str); err == nil {
		ret = uint16(intValue)
	}
	return
}

func BuildReadToPositionWrapper(readTo ReadTo) ReadTo {
	return func(fillItem *Item, reader *Reader) (err error) {
		fillItem.SetStartPos(reader)
		err = readTo(fillItem, reader)
		fillItem.SetEndPos(reader)
		return
	}
}
