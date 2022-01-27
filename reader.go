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
	Attr() *Attr
	Accessor() interface{}
	ReadTo(fillItem *Item, reader *Reader) (err error)
	NewItem(parent *Item, value interface{}) *Item
}

type ReadFix func(reader *Reader) (ret interface{}, err error)
type ReadDynamic func(reader *Reader, parent *Item) (ret interface{}, err error)
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
	Value    interface{}
	Parent   *Item
	StartPos int64
	EndPos   int64
	Attr     *Attr
	Accessor interface{}
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

func (o *Item) SetStartPos(reader *Reader) {
	o.StartPos = reader.Position()
}

func (o *Item) SetEndPos(reader *Reader) {
	o.EndPos = reader.Position()
}

type Reader struct {
	io.ReadSeeker
	offset int64
	buf    [8]byte

	// Number of bits remaining in "bits" for sequential calls to ReadBitsInt
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

func (o *Reader) ReadBytesAsReader(n uint16) (ret *Reader, err error) {
	currentPos := o.Position()
	var raw []byte
	if raw, err = o.ReadBytes(n); err == nil {
		ret = &Reader{ReadSeeker: bytes.NewReader(raw), offset: currentPos}
	}
	return
}

func (o *Reader) Position() (ret int64) {
	ret, _ = o.Seek(0, io.SeekCurrent)
	return o.offset + ret
}

func ReadFixAttr(attr *Attr, convert Convert) (ret ReadFix) {
	if attr.SizeEos == "true" {
		ret = ReadFixFull(convert)
	} else if attr.Size != "" {
		if length, err := strconv.Atoi(attr.Size); err == nil {
			ret = ReadFixLength(uint16(length), convert)
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
	return func(reader *Reader) (ret interface{}, err error) {
		return ReadFull(reader, convert)
	}
}

func ReadFull(reader *Reader, convert Convert) (ret interface{}, err error) {
	var data []byte
	if data, err = reader.ReadBytesFull(); err == nil {
		ret, err = convert(data)
	}
	return
}

func ReadFixLength(length uint16, convert Convert) (ret ReadFix) {
	return func(reader *Reader) (ret interface{}, err error) {
		return ReadLength(reader, length, convert)
	}
}

func ReadLength(reader *Reader, length uint16, convert Convert) (ret interface{}, err error) {
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
	return func(reader *Reader, parent *Item) (ret interface{}, err error) {
		var sizeItem *Item
		if sizeItem, err = parent.Expr(expr); err == nil {
			var length uint16
			if length, err = toUint16(sizeItem.Value); err == nil {
				ret, err = ReadLength(reader, length, convert)
			} else {
				err = fmt.Errorf("cant parse Size to uint64, expr=%v, valiue=%v, %v", expr, sizeItem.Value, err)
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

func toUint16(value interface{}) (ret uint16, err error) {
	var ok bool
	if ret, ok = value.(uint16); !ok {
		str := fmt.Sprintf("%v", value)
		var intValue int
		if intValue, err = strconv.Atoi(str); err != nil {
			ret = uint16(intValue)
		}
	}
	return
}
