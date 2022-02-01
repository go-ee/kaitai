package kaitai

type Type struct {
	Id  string  `-`
	Seq []*Attr `yaml:"seq,omitempty"`
	Doc string  `yaml:"doc,omitempty"`
}

func (o *Type) BuildReader(attr *Attr, spec *Spec) (ret Reader, err error) {
	var seqReaders []Reader
	if seqReaders, err = o.buildSeqReaders(spec); err == nil {
		ret = WrapReader(&TypeReader{
			ReaderBase: &ReaderBase{attr: attr, accessor: o}, readers: seqReaders,
		}, spec.Options)
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
	*ReaderBase
	readers []Reader
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
