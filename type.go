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
			attrReaders:    readers,
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
	attrReaders []AttrReader
}

func (o *TypeReader) Read(reader Reader, parent *Item, root *Item) (ret *Item, err error) {
	data := map[string]*Item{}
	ret = o.newItem(data)

	parent = ret
	if root == nil {
		root = ret
	}

	for _, attrReader := range o.attrReaders {
		if item, currentErr := attrReader.Read(reader, parent, root); currentErr == nil {
			data[item.Attr.Id] = item
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
