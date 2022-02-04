package kaitai

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
)

type Reader interface {
	NewItem(parent *Item) *Item
	ReadTo(fillItem *Item, reader *ReaderIO) (err error)
}

type NativeReader interface {
	Read(parent *Item, reader *ReaderIO) (ret interface{}, err error)
}

type Read func(reader *ReaderIO) (ret interface{}, err error)
type ParentRead func(parent *Item, reader *ReaderIO) (ret interface{}, err error)
type ReadTo func(fillItem *Item, reader *ReaderIO) (err error)
type Parse func(data []byte) (interface{}, error)
type Decode func(fillItem *Item)

type Item struct {
	value    interface{}
	Err      error
	Parent   *Item
	StartPos *int64
	EndPos   *int64
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

func (o *Item) ExprValue(expr string) (ret interface{}, err error) {
	if ret, err = o.Expr(expr); err == nil {
		if item, ok := ret.(*Item); ok {
			ret = item.Value()
		}
	}
	return
}

func (o *Item) Expr(expr string) (ret interface{}, err error) {
	ret = o
	for _, part := range strings.Split(expr, ".") {
		if item, ok := ret.(*Item); ok {
			ret = item.Value()
		}

		if value, ok := ret.(map[string]interface{}); ok {
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
	pos := reader.Position()
	o.StartPos = &pos
}

func (o *Item) SetEndPos(reader *ReaderIO) {
	pos := reader.Position()
	o.EndPos = &pos
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

type ReaderBase struct {
	attr     *Attr
	accessor interface{}
}

func (o *ReaderBase) Attr() *Attr {
	return o.attr
}

func (o *ReaderBase) Accessor() interface{} {
	return o.accessor
}

func (o *ReaderBase) NewItem(parent *Item) *Item {
	return &Item{Attr: o.attr, Accessor: o.accessor, Parent: parent}
}
