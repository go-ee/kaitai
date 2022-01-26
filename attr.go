package kaitai

import "io"

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

	itemReader ItemReader `-`
}

func (o *Attr) BuildReader(attr *Attr, spec *Spec) (ret ItemReader, err error) {
	return
}

func (o *Attr) Read(reader ReaderIO, parent *Item, root *Item) (ret *Item, err error) {
	if o.Repeat == "eos" {
		var items []*Item
		for i := 0; err == nil; i++ {
			var item *Item
			if item, err = o.itemReader.Read(reader, parent, root); err == nil {
				items = append(items, item)
			}
		}

		if io.EOF == err {
			err = nil
			ret = &Item{Accessor: o, Value: items}
		}
	} else {
		ret, err = o.itemReader.Read(reader, parent, root)
	}
	return
}

func (o *Attr) crossInit(spec *Spec) (err error) {
	if o.Type != nil {
		o.itemReader, err = o.Type.BuildReader(o, spec)
	} else if o.Contents != nil {
		o.itemReader, err = o.Contents.BuildReader(o, spec)
	} else if o.SizeEos == "true" {
		//o.itemReader = ItemReader
	} else {
		//o.itemReader = &AttrReader{o.Id}
	}
	return
}
