package kaitai

type Type struct {
	Id  string  `-`
	Seq []*Attr `yaml:"seq,omitempty"`
	Doc string  `yaml:"doc,omitempty"`
}

func (o *Type) BuildReader(attr *Attr, spec *Spec) (ret Reader, err error) {
	var seqReaders []Reader
	var typeModel *TypeModel
	if seqReaders, typeModel, err = o.buildSeqReaders(spec); err == nil {
		typeReader := &TypeReader{o, typeModel, attr, seqReaders}
		if spec.Options.PositionFill {
			ret = &SetPositionTypeReader{typeReader}
		} else {
			ret = typeReader
		}
	}
	return
}

func (o *Type) buildSeqReaders(spec *Spec) (ret []Reader, model *TypeModel, err error) {
	model = NewTypeModel()
	readers := make([]Reader, len(o.Seq))
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
	attr    *Attr
	readers []Reader
}

func (o *TypeReader) Read(_ Item, reader *ReaderIO) (ret interface{}, err error) {
	ret, err = o.readTo(o.NewItem(), reader)
	return
}

func (o *TypeReader) readTo(item Item, reader *ReaderIO) (ret interface{}, err error) {
	ret = item
	for i, attrReader := range o.readers {
		if attrValue, attrErr := attrReader.Read(item, reader); attrErr == nil {
			item[i] = attrValue
		} else {
			err = attrErr
			break
		}
	}
	return
}

func (o TypeReader) NewItem() (ret Item) {
	ret = Item{}
	ret.SetModel(o.model)
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
