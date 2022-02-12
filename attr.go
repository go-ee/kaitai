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
		itemReader = &AttrParentRead{attr: o, parentRead: BuildReadToFull(ToSame)}
	} else {
		err = fmt.Errorf("read attr: ELSE, not implemented yet")
	}

	if o.Repeat == "eos" {
		ret = &AttrCycleReader{itemReader: itemReader}
	} else if o.Size != "" {
		ret = &AttrSizeReader{attr: o, itemReader: itemReader}
	} else {
		ret = itemReader
	}
	return
}

type AttrCycleReader struct {
	itemReader Reader
}

func (o *AttrCycleReader) Read(parent *TypeItem, reader *ReaderIO) (ret interface{}, err error) {
	var items []interface{}
	for i := 0; err == nil; i++ {
		var item interface{}
		if item, err = o.itemReader.Read(parent, reader); err != nil {
			break
		}
		items = append(items, item)
	}

	if io.EOF == err {
		err = nil
	}

	if err == nil {
		ret = items
	}
	return
}

type AttrSizeReader struct {
	attr       *Attr
	itemReader Reader
}

func (o *AttrSizeReader) Read(parent *TypeItem, reader *ReaderIO) (ret interface{}, err error) {
	var size interface{}
	if size, err = parent.Expr(o.attr.Size); err != nil {
		return
	}

	var length uint16
	if length, err = ToUint16(size); err != nil {
		return
	}

	offset := reader.Position()

	var data []byte
	if data, err = reader.ReadBytes(length); err != nil {
		return
	}

	childReader := &ReaderIO{ReadSeeker: bytes.NewReader(data), offset: offset}
	ret, err = o.itemReader.Read(parent, childReader)

	if io.EOF == err {
		err = nil
	}
	return
}
