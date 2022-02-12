package kaitai

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
)

type Reader interface {
	Read(parent *TypeItem, reader *ReaderIO) (ret interface{}, err error)
}

type ReadItem func(parent *TypeItem, reader *ReaderIO) (ret *TypeItem, err error)
type ParentRead func(parent *TypeItem, reader *ReaderIO) (ret interface{}, err error)
type Parse func(data []byte) (interface{}, error)
type Decode func() (interface{}, error)

type TypeModel struct {
	attrs       []*Attr
	names       []string
	indexByName map[string]int
}

func NewTypeModel(attrCount int) (ret *TypeModel) {
	ret = &TypeModel{
		attrs:       make([]*Attr, attrCount),
		names:       make([]string, attrCount),
		indexByName: make(map[string]int, attrCount),
	}
	return
}

func (o *TypeModel) AddAttr(attrIndex int, attr *Attr) {
	o.attrs[attrIndex] = attr
	o.names[attrIndex] = attr.Id
	o.indexByName[attr.Id] = attrIndex
}

func (o *TypeModel) IndexToAttrName(index int) string {
	return o.names[index]
}

func (o *TypeModel) IndexToAttr(index int) *Attr {
	return o.attrs[index]
}

func (o *TypeModel) AttrToIndex(attr string) int {
	return o.indexByName[attr]
}

type TypeItem struct {
	model       *TypeModel
	attrsDecode []Decode
	attrs       []interface{}
	startPos    int64
	endPos      int64
}

func NewTypeItem(model *TypeModel) *TypeItem {
	return &TypeItem{model: model, attrs: make([]interface{}, len(model.attrs))}
}

func (o *TypeItem) IndexToAttr(index int) string {
	return o.model.IndexToAttrName(index)
}

func (o *TypeItem) AttrToIndex(attr string) int {
	return o.model.AttrToIndex(attr)
}

func (o *TypeItem) SetAttrValue(attrIndex int, attrValue interface{}) {
	o.attrs[attrIndex] = attrValue
}

func (o *TypeItem) GetAttrValue(attrIndex int) (ret interface{}) {
	ret = o.attrs[attrIndex]
	if ret == nil && o.attrsDecode != nil {
		if decode := o.attrsDecode[attrIndex]; decode != nil {
			if value, err := decode(); err != nil {
				o.attrs[attrIndex] = err
			} else {
				o.attrs[attrIndex] = value
			}
		}
	}
	return
}

func (o *TypeItem) ExprValue(expr string) (ret interface{}, err error) {
	ret, err = o.Expr(expr)
	return
}

func (o *TypeItem) Expr(expr string) (ret interface{}, err error) {
	ret = o
	for _, part := range strings.Split(expr, ".") {
		if value, ok := ret.(*TypeItem); ok {
			index := value.AttrToIndex(part)
			if ret = value.attrs[index]; ret == nil {
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

func (o *TypeItem) SetStartPos(reader *ReaderIO) {
	o.startPos = reader.Position()
}

func (o *TypeItem) SetEndPos(reader *ReaderIO) {
	o.endPos = reader.Position()
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
	ret = o.offset + ret
	return
}
