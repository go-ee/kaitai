package kaitai

import (
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
)

type AttrReader interface {
	Attr() *Attr
	Accessor() interface{}
	ReadTo(fillItem *Item, reader Reader) (err error)
	NewItem(parent *Item, value interface{}) *Item
}

type ReadFix func(reader Reader) (ret interface{}, err error)
type ReadDynamic func(reader Reader, parent *Item) (ret interface{}, err error)
type Convert func(data []byte) (ret interface{}, err error)

type AttrReaderBase struct {
	attr     *Attr
	accessor interface{}
}

func (o *AttrReaderBase) Attr() *Attr {
	return o.attr
}

func (o *AttrReaderBase) Accessor() interface{} {
	return o.accessor
}

func (o *AttrReaderBase) NewItem(parent *Item, value interface{}) *Item {
	return &Item{Attr: o.attr, Accessor: o.accessor, Value: value, Parent: parent}
}

type Item struct {
	Attr     *Attr
	Accessor interface{}
	Value    interface{}
	Parent   *Item
}

func (o *Item) Expr(expr string) (ret *Item, err error) {
	ret = o
	for _, part := range strings.Split(expr, ".") {
		if value, ok := ret.Value.(map[string]*Item); ok {
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

type Reader struct {
	io.ReadSeeker
	buf [8]byte

	// Number of bits remaining in "bits" for sequential calls to ReadBitsInt
	bitsLeft uint8
	bits     uint64
}

func (o *Reader) ReadBytes(n uint8) (ret []byte, err error) {
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

func ReadFixAttr(attr *Attr, convert Convert) (ret ReadFix) {
	if attr.SizeEos == "true" {
		ret = ReadFixFull(convert)
	} else if attr.Size != "" {
		if length, err := strconv.Atoi(attr.Size); err == nil {
			ret = ReadFixLength(uint8(length), convert)
		}
	}
	return
}

func ReadDynamicAttr(attr *Attr, convert Convert) (ret ReadDynamic) {
	if attr.Size != "" {
		if _, err := strconv.Atoi(attr.Size); err != nil {
			ret = ReadDynamicLengthExpr(attr.Size, convert)
		}
	}
	return
}

func ReadFixFull(convert Convert) (ret ReadFix) {
	return func(reader Reader) (ret interface{}, err error) {
		return ReadFull(reader, convert)
	}
}

func ReadFull(reader Reader, convert Convert) (ret interface{}, err error) {
	var data []byte
	if data, err = reader.ReadBytesFull(); err == nil {
		ret, err = convert(data)
	}
	return
}

func ReadFixLength(length uint8, convert Convert) (ret ReadFix) {
	return func(reader Reader) (ret interface{}, err error) {
		return ReadLength(reader, length, convert)
	}
}

func ReadLength(reader Reader, length uint8, convert Convert) (ret interface{}, err error) {
	var data []byte
	if length > 0 {
		data, err = reader.ReadBytes(length)
	} else {
		data, err = reader.ReadBytesFull()
	}

	if err == nil {
		ret, err = convert(data)
	}
	return
}

func ReadDynamicLengthExpr(expr string, convert Convert) (ret ReadDynamic) {
	return func(reader Reader, parent *Item) (ret interface{}, err error) {
		var sizeItem *Item
		if sizeItem, err = parent.Expr(expr); err == nil {
			if length, ok := sizeItem.Value.(uint8); ok {
				ret, err = ReadLength(reader, length, convert)
			} else {
				err = fmt.Errorf("cant parse Size to uint8, expr=%v, valiue=%v", expr, sizeItem.Value)
			}
		}
		return
	}
}

func ToString(data []byte) (ret interface{}, _ error) {
	ret = string(data)
	return
}

func ToSame(data []byte) (ret interface{}, _ error) {
	ret = data
	return
}
