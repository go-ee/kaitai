package kaitai

type Type struct {
	Id  string  `-`
	Seq []*Attr `yaml:"seq,omitempty"`
	Doc string  `yaml:"doc,omitempty"`
}

func (o *Type) BuildReader(attr *Attr, spec *Spec) (ret ItemReader, err error) {
	if err = o.crossInit(spec); err != nil {
		return
	}

	var seqReaders []*AttrReader
	if seqReaders, err = o.buildSeqReaders(spec); err == nil {
		typeReader := &TypeReader{Attr: attr, Accessor: o, Seq: seqReaders}
		ret = typeReader
	}
	return
}

func (o *Type) buildSeqReaders(spec *Spec) (ret []*AttrReader, err error) {
	seqReaders := make([]*AttrReader, len(o.Seq))
	for i, seqAttr := range o.Seq {
		if itemReader, currentErr := seqAttr.BuildReader(spec); currentErr == nil {
			seqReaders[i] = &AttrReader{Attr: seqAttr, ItemReader: itemReader}
		} else {
			err = currentErr
			return
		}
	}

	if err == nil {
		ret = seqReaders
	}
	return
}

func (o *Type) crossInit(base *Spec) (err error) {
	for _, item := range o.Seq {
		if err = item.crossInit(base); err != nil {
			break
		}
	}
	return
}

type TypeReader struct {
	Attr     *Attr
	Accessor *Type
	Seq      []*AttrReader
}

func (o *TypeReader) Read(reader ReaderIO, parent *Item, root *Item) (ret *Item, err error) {
	data := map[string]*Item{}
	ret = &Item{Attr: o.Attr, Accessor: o.Accessor, Value: data}

	parent = ret
	if root == nil {
		root = ret
	}

	for _, attrReader := range o.Seq {
		if item, currentErr := attrReader.ItemReader.Read(reader, parent, root); currentErr == nil {
			data[attrReader.Attr.Id] = item
		} else {
			err = currentErr
			break
		}
	}

	if err != nil {
		ret = nil
	}
	return
}

type AttrReader struct {
	Attr       *Attr
	ItemReader ItemReader
}
