package kaitai

type Type struct {
	Id  string  `-`
	Seq []*Attr `yaml:"seq,omitempty"`
	Doc string  `yaml:"doc,omitempty"`
}

func (o *Type) BuildReader(attr *Attr, spec *Spec) (ret AttrReader, err error) {
	var seqReaders []AttrReader
	if seqReaders, err = o.buildSeqReaders(spec); err == nil {
		typeReader := &TypeReader{o, attr, seqReaders}
		if spec.Options.PositionFill {
			ret = &SetPositionTypeReader{typeReader}
		} else {
			ret = typeReader
		}
	}
	return
}

func (o *Type) buildSeqReaders(spec *Spec) (ret []AttrReader, err error) {
	readers := make([]AttrReader, len(o.Seq))
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
	*Type
	attr    *Attr
	readers []AttrReader
}

func (o *TypeReader) Read(parent *Item, reader *ReaderIO) (ret interface{}, err error) {
	item := o.buildItem(parent)
	ret, err = o.readTo(item, reader)
	return
}

func (o *TypeReader) buildItem(parent *Item) *Item {
	return &Item{Attr: o.attr, Type: o.Type, Parent: parent, value: map[string]interface{}{}}
}

func (o *TypeReader) readTo(item *Item, reader *ReaderIO) (ret interface{}, err error) {
	data := item.value.(map[string]interface{})
	for _, attrReader := range o.readers {
		if attrValue, attrErr := attrReader.Read(item, reader); attrErr == nil {
			data[attrReader.Attr().Id] = attrValue
		} else {
			err = attrErr
			break
		}
	}

	if err == nil {
		ret = item
	}
	return
}

func (o TypeReader) Attr() *Attr {
	return o.attr
}

type SetPositionTypeReader struct {
	*TypeReader
}

func (o *SetPositionTypeReader) Read(parent *Item, reader *ReaderIO) (ret interface{}, err error) {
	item := o.buildItem(parent)
	item.SetStartPos(reader)
	ret, err = o.readTo(item, reader)
	item.SetEndPos(reader)
	return
}

func (o SetPositionTypeReader) Attr() *Attr {
	return o.attr
}
