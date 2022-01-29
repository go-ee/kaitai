package kaitai

type Type struct {
	Id  string  `-`
	Seq []*Attr `yaml:"seq,omitempty"`
	Doc string  `yaml:"doc,omitempty"`
}

func (o *Type) BuildReader(attr *Attr, spec *Spec) (ret Reader, err error) {
	var seqReaders []Reader
	if seqReaders, err = o.buildSeqReaders(spec); err == nil {
		ret = WrapReader(&TypeReader{attr: attr, accessor: o, readers: seqReaders}, spec.Options)
	}
	return
}

func (o *Type) buildSeqReaders(spec *Spec) (ret []Reader, err error) {
	readers := make([]Reader, len(o.Seq))
	for i, attr := range o.Seq {
		if readers[i], err = attr.BuildReader(spec); err != nil {
			return
		}
	}

	if err == nil {
		ret = readers
	}
	return
}

type TypeReader struct {
	attr     *Attr
	accessor interface{}
	readers  []Reader
}

func (o *TypeReader) ReadTo(fillItem *Item, reader *ReaderIO) (err error) {
	data := map[string]*Item{}
	fillItem.SetValue(data)
	for _, attrReader := range o.readers {
		item := attrReader.NewItem(fillItem)
		data[item.Attr.Id] = item
		if err = attrReader.ReadTo(item, reader); err != nil {
			break
		}
	}
	return
}

func (o *TypeReader) Attr() *Attr {
	return o.attr
}

func (o *TypeReader) Accessor() interface{} {
	return o.accessor
}

func (o *TypeReader) NewItem(parent *Item) *Item {
	return &Item{Attr: o.attr, Accessor: o.accessor, Parent: parent}
}
