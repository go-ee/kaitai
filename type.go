package kaitai

type Type struct {
	Id  string  `-`
	Seq []*Attr `yaml:"seq,omitempty"`
	Doc string  `yaml:"doc,omitempty"`
}

func (o *Type) BuildReader(attr *Attr, spec *Spec) (ret AttrReader, err error) {
	var seqReaders []*AttrReaderHolder
	if seqReaders, err = o.buildSeqReaders(spec); err == nil {
		typeReader := &TypeReader{
			AttrReaderBase: &AttrReaderBase{attr, o},
			seqReaders:     seqReaders,
		}
		ret = typeReader
	}
	return
}

func (o *Type) buildSeqReaders(spec *Spec) (ret []*AttrReaderHolder, err error) {
	seqReaders := make([]*AttrReaderHolder, len(o.Seq))
	for i, seqAttr := range o.Seq {
		if itemReader, currentErr := seqAttr.BuildReader(spec); currentErr == nil {
			seqReaders[i] = &AttrReaderHolder{Attr: seqAttr, AttrReader: itemReader}
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

type TypeReader struct {
	*AttrReaderBase
	seqReaders []*AttrReaderHolder
}

func (o *TypeReader) Read(reader Reader, parent *Item, root *Item) (ret *Item, err error) {
	data := map[string]*Item{}
	ret = o.newItem(data)

	parent = ret
	if root == nil {
		root = ret
	}

	for _, attrReader := range o.seqReaders {
		if item, currentErr := attrReader.AttrReader.Read(reader, parent, root); currentErr == nil {
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

type AttrReaderHolder struct {
	Attr       *Attr
	AttrReader AttrReader
}
