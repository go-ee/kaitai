package kaitai

import (
	"fmt"
	"io"
)

type Attr struct {
	Id          string    `yaml:"id,omitempty"`
	Type        *TypeRef  `yaml:"type"`
	Size        string    `yaml:"size,omitempty"`
	SizeEos     string    `yaml:"size-eos,omitempty"`
	Doc         string    `yaml:"doc,omitempty"`
	Repeat      string    `yaml:"repeat,omitempty"`
	RepeatExpr  string    `yaml:"repeat-expr,omitempty"`
	RepeatUntil string    `yaml:"repeat-until,omitempty"`
	Contents    *Contents `yaml:"contents,omitempty"`
	Value       string    `yaml:"value,omitempty"`
	Pos         string    `yaml:"pos,omitempty"`
	Whence      string    `yaml:"whence,omitempty"`
	Enum        string    `yaml:"enum,omitempty"`
	If          string    `yaml:"if,omitempty"`
	Process     string    `yaml:"process,omitempty"`
	Terminator  string    `yaml:"terminator,omitempty"`
	Consume     string    `yaml:"consume,omitempty"`
	Include     string    `yaml:"include,omitempty"`
	EosError    string    `yaml:"eos-error,omitempty"`
	PadRight    string    `yaml:"pad-right,omitempty"`
	Encoding    string    `yaml:"encoding,omitempty"`
}

func (o *Attr) BuildReader(spec *Spec) (ret AttrReader, err error) {
	var itemReader AttrReader

	if o.Type != nil {
		itemReader, err = o.Type.BuildReader(o, spec)
	} else if o.Contents != nil {
		itemReader, err = o.Contents.BuildReader(o, spec)
	} else if o.SizeEos == "true" {
		itemReader = &NativeReaderFix{&AttrReaderBase{attr: o}, ReadFixFull(ToSame)}
	} else {
		err = fmt.Errorf("read attr: ELSE, not implemented yet")
	}

	if o.Repeat == "eos" {
		ret = &AttrCycleReader{&AttrReaderBase{attr: o}, itemReader}
	} else if o.Size != "" {
		ret = &AttrSizeReader{&AttrReaderBase{attr: o}, itemReader}
	} else {
		ret = itemReader
	}
	return
}

type AttrCycleReader struct {
	*AttrReaderBase
	itemReader AttrReader
}

func (o *AttrCycleReader) ReadTo(fillItem *Item, reader *Reader) (err error) {

	fillItem.SetStartPos(reader)

	var items []*Item
	for i := 0; err == nil; i++ {
		item := o.itemReader.NewItem(fillItem, nil)
		items = append(items, item)
		fillItem.Value = items
		err = o.itemReader.ReadTo(item, reader)
	}

	if io.EOF == err {
		err = nil
	}

	fillItem.SetEndPos(reader)

	return
}

type AttrSizeReader struct {
	*AttrReaderBase
	itemReader AttrReader
}

func (o *AttrSizeReader) ReadTo(fillItem *Item, reader *Reader) (err error) {

	fillItem.SetStartPos(reader)

	var sizeItem *Item
	if sizeItem, err = fillItem.Parent.Expr(o.attr.Size); err != nil {
		return
	}

	var length uint16
	if length, err = toUint16(sizeItem.Value); err != nil {
		return
	}

	var sizeReader *Reader
	if sizeReader, err = reader.ReadBytesAsReader(length); err != nil {
		return
	}

	err = o.itemReader.ReadTo(fillItem, sizeReader)

	if io.EOF == err {
		err = nil
	}

	fillItem.SetEndPos(reader)

	return
}
