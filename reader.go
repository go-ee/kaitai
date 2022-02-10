package kaitai

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
)

type Reader interface {
	Read(parent *Item, reader *ReaderIO) (ret interface{}, err error)
}

type ReadItem func(parent *Item, reader *ReaderIO) (ret *Item, err error)
type ParentRead func(parent *Item, reader *ReaderIO) (ret interface{}, err error)
type Parse func(data []byte) (interface{}, error)
type Decode func(fillItem *Item)

type TypeModel struct {
	indexToAttr []*Attr
	indexToName []string
	nameToIndex map[string]int
}

func NewTypeModel(attrCount int) (ret *TypeModel) {
	ret = &TypeModel{
		indexToAttr: make([]*Attr, attrCount),
		indexToName: make([]string, attrCount),
		nameToIndex: make(map[string]int, attrCount),
	}
	return
}

func (o *TypeModel) SetAttr(attrIndex int, attr *Attr) {
	o.indexToAttr[attrIndex] = attr
	o.indexToName[attrIndex] = attr.Id
	o.nameToIndex[attr.Id] = attrIndex
}

func (o *TypeModel) IndexToAttrName(index int) string {
	return o.indexToName[index]
}

func (o *TypeModel) IndexToAttr(index int) *Attr {
	return o.indexToAttr[index]
}

func (o *TypeModel) AttrToIndex(attr string) int {
	return o.nameToIndex[attr]
}

type Item struct {
	Model    *TypeModel
	Attrs    []interface{}
	StartPos int64
	EndPos   int64
}

func (o *Item) IndexToAttr(index int) string {
	return o.Model.IndexToAttrName(index)
}

func (o *Item) AttrToIndex(attr string) int {
	return o.Model.AttrToIndex(attr)
}

func (o *Item) ExprValue(expr string) (ret interface{}, err error) {
	ret, err = o.Expr(expr)
	return
}

func (o *Item) Expr(expr string) (ret interface{}, err error) {
	ret = o
	for _, part := range strings.Split(expr, ".") {
		if value, ok := ret.(*Item); ok {
			index := value.AttrToIndex(part)
			if ret = value.Attrs[index]; ret == nil {
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
