package kaitai

type Type struct {
	Id  string  `-`
	Seq []*Attr `yaml:"seq,omitempty"`
	Doc string  `yaml:"doc,omitempty"`
}

func (o *Type) BuildReader(attr *Attr, spec *Spec) (ret AttrReader, err error) {
	var readers []AttrReader
	if readers, err = o.buildSeqReaders(spec); err == nil {
		typeReader := &TypeReader{
			AttrReaderBase: &AttrReaderBase{attr, o},
			readers:        readers,
		}
		ret = typeReader
	}
	return
}

func (o *Type) buildSeqReaders(spec *Spec) (ret []AttrReader, err error) {
	readers := make([]AttrReader, len(o.Seq))
	for i, seqAttr := range o.Seq {
		if readers[i], err = seqAttr.BuildReader(spec); err != nil {
			return
		}
	}

	if err == nil {
		ret = readers
	}
	return
}

type TypeReader struct {
	*AttrReaderBase
	readers []AttrReader
}

func (o *TypeReader) ReadTo(fillItem *Item, reader Reader) (err error) {
	data := map[string]*Item{}
	fillItem.Value = data

	for _, attrReader := range o.readers {
		item := attrReader.NewItem(fillItem, nil)
		data[item.Attr.Id] = item
		if err = attrReader.ReadTo(item, reader); err != nil {
			break
		}
	}
	return
}
