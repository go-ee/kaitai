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

func (o *TypeReader) Read(_ Item, reader *ReaderIO) (ret interface{}, err error) {
	ret, err = o.readTo(Item{}, reader)
	return
}

func (o *TypeReader) readTo(item Item, reader *ReaderIO) (ret interface{}, err error) {
	ret = item
	for _, attrReader := range o.readers {
		attrName := attrReader.Attr().Id
		if attrValue, attrErr := attrReader.Read(item, reader); attrErr == nil {
			item[attrName] = attrValue
		} else {
			err = attrErr
			break
		}
	}
	return
}

func (o TypeReader) Attr() *Attr {
	return o.attr
}

type SetPositionTypeReader struct {
	*TypeReader
}

func (o *SetPositionTypeReader) Read(_ Item, reader *ReaderIO) (ret interface{}, err error) {
	item := Item{}
	item.SetStartPos(reader)
	ret, err = o.readTo(item, reader)
	item.SetEndPos(reader)
	return
}

func (o SetPositionTypeReader) Attr() *Attr {
	return o.attr
}
