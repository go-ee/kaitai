package kaitai

type Type struct {
	Id  string  `-`
	Seq []*Attr `yaml:"seq,omitempty"`
	Doc string  `yaml:"doc,omitempty"`
}

func (o *Type) BuildReader(_ *Attr, spec *Spec) (ret Reader, err error) {
	var seqReaders []Reader
	var typeModel *TypeModel
	if seqReaders, typeModel, err = o.buildSeqReaders(spec); err == nil {
		typeReader := &TypeReader{o, typeModel, seqReaders}
		if spec.Options.PositionFill {
			ret = &SetPositionTypeReader{typeReader}
		} else {
			ret = typeReader
		}
	}
	return
}

func (o *Type) buildSeqReaders(spec *Spec) (ret []Reader, model *TypeModel, err error) {
	attrCount := len(o.Seq)
	model = NewTypeModel(attrCount)
	readers := make([]Reader, attrCount)
	for i, attr := range o.Seq {
		model.AddAttr(i, attr)
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
	model   *TypeModel
	readers []Reader
}

func (o *TypeReader) Read(_ *TypeItem, reader *ReaderIO) (ret interface{}, err error) {
	ret, err = o.readTo(NewTypeItem(o.model), reader)
	return
}

func (o *TypeReader) readTo(item *TypeItem, reader *ReaderIO) (ret interface{}, err error) {
	ret = item
	for i, attrReader := range o.readers {
		if attrValue, attrErr := attrReader.Read(item, reader); attrErr == nil {
			item.SetAttrValue(i, attrValue)
		} else {
			err = attrErr
			break
		}
	}
	return
}

type SetPositionTypeReader struct {
	*TypeReader
}

func (o *SetPositionTypeReader) Read(_ *TypeItem, reader *ReaderIO) (ret interface{}, err error) {
	item := NewTypeItem(o.model)
	item.SetStartPos(reader)
	ret, err = o.readTo(item, reader)
	item.SetEndPos(reader)
	return
}
