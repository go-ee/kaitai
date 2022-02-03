package kaitai

import (
	"bytes"
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

func (o *Attr) BuildReader(spec *Spec) (ret Reader, err error) {
	var itemReader Reader

	if o.Type != nil {
		itemReader, err = o.Type.BuildReader(o, spec)
	} else if o.Contents != nil {
		itemReader, err = o.Contents.BuildReader(o, spec)
	} else if o.SizeEos == "true" {
		itemReader = &AttrAccessorReadToReader{ReaderBase: &ReaderBase{attr: o}, readTo: BuildReadToFull(ToSame)}
	} else {
		err = fmt.Errorf("read attr: ELSE, not implemented yet")
	}

	if o.Repeat == "eos" {
		ret = &AttrCycleReader{ReaderBase: &ReaderBase{attr: o, accessor: o}, itemReader: itemReader}
	} else if o.Size != "" {
		if spec.Options.LazyDecoding {
			ret = &AttrSizeLazyReader{ReaderBase: &ReaderBase{attr: o}, itemReader: itemReader}
		} else {
			ret = &AttrSizeReader{ReaderBase: &ReaderBase{attr: o}, itemReader: itemReader}
		}
	} else {
		ret = itemReader
	}
	return
}

type AttrCycleReader struct {
	*ReaderBase
	itemReader Reader
}

func (o *AttrCycleReader) ReadTo(fillItem *Item, reader *ReaderIO) (err error) {
	var items []*Item
	for i := 0; err == nil; i++ {
		item := o.itemReader.NewItem(fillItem)
		items = append(items, item)
		fillItem.SetValue(items)
		err = o.itemReader.ReadTo(item, reader)
	}

	if io.EOF == err {
		err = nil
	}
	return
}

type AttrSizeReader struct {
	*ReaderBase
	itemReader Reader
}

func (o *AttrSizeReader) ReadTo(fillItem *Item, reader *ReaderIO) (err error) {
	var sizeItem *Item
	if sizeItem, err = fillItem.Parent.Expr(o.attr.Size); err != nil {
		return
	}

	var length uint16
	if length, err = ToUint16(sizeItem.Value()); err != nil {
		return
	}

	var data []byte
	if data, err = reader.ReadBytes(length); err != nil {
		return
	}

	childReader := &ReaderIO{ReadSeeker: bytes.NewReader(data), offset: reader.Position()}
	err = o.itemReader.ReadTo(fillItem, childReader)

	if io.EOF == err {
		err = nil
	}
	return
}

type AttrSizeLazyReader struct {
	*ReaderBase
	itemReader Reader
}

func (o *AttrSizeLazyReader) ReadTo(fillItem *Item, reader *ReaderIO) (err error) {
	var sizeItem *Item
	if sizeItem, err = fillItem.Parent.Expr(o.attr.Size); err != nil {
		return
	}

	var length uint16
	if length, err = ToUint16(sizeItem.Value()); err != nil {
		return
	}

	parser := RawReaderParser{offset: reader.Position(), itemReader: o.itemReader}
	fillItem.Decode = parser.Decode
	fillItem.Raw, err = reader.ReadBytes(length)
	return
}

type RawReaderParser struct {
	offset     int64
	itemReader Reader
}

func (o *RawReaderParser) Decode(fillItem *Item) {
	reader := &ReaderIO{ReadSeeker: bytes.NewReader(fillItem.Raw), offset: o.offset}
	err := o.itemReader.ReadTo(fillItem, reader)

	if io.EOF == err {
		err = nil
	}

	if err != nil {
		fillItem.Err = err
	}
	return
}
