package kaitai

type Type struct {
	Id  string  `-`
	Seq []*Attr `yaml:"seq,omitempty"`
	Doc string  `yaml:"doc,omitempty"`
}

func (o *Type) BuildReader(attr *Attr, spec *Spec) (ret *TypeReader, err error) {
	var seqReaders []AttrReader
	if seqReaders, err = o.buildSeqReaders(spec); err == nil {
		ret = &TypeReader{o, attr, seqReaders}
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

func (o *TypeReader) Read(_ *Item, reader *ReaderIO) (ret interface{}, err error) {
	data := map[string]interface{}{}
	item := &Item{value: data, Type: o.Type}
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
