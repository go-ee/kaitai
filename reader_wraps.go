package kaitai

type ReadToWrapper func(readTo ReadTo) ReadTo

type AttrAccessorReadToReader struct {
	attr     *Attr
	accessor interface{}
	readTo   ReadTo
}

func (o *AttrAccessorReadToReader) ReadTo(fillItem *Item, reader *ReaderIO) (err error) {
	return o.readTo(fillItem, reader)
}

func (o *AttrAccessorReadToReader) NewItem(parent *Item) *Item {
	return &Item{Attr: o.attr, Accessor: o.accessor, Parent: parent}
}

func (o *AttrAccessorReadToReader) Attr() *Attr {
	return o.attr
}

func (o *AttrAccessorReadToReader) Accessor() interface{} {
	return o.accessor
}

type ReadToWrapperReader struct {
	reader Reader
	readTo ReadTo
}

func NewReadToWrapperReader(reader Reader, readToWrapper ReadToWrapper) *ReadToWrapperReader {
	return &ReadToWrapperReader{reader: reader, readTo: readToWrapper(reader.ReadTo)}
}

func (o *ReadToWrapperReader) ReadTo(fillItem *Item, reader *ReaderIO) (err error) {
	return o.readTo(fillItem, reader)
}

func (o *ReadToWrapperReader) NewItem(parent *Item) *Item {
	return o.reader.NewItem(parent)
}

func ReadToPositionWrapper(readTo ReadTo) ReadTo {
	return func(fillItem *Item, reader *ReaderIO) (err error) {
		fillItem.SetStartPos(reader)
		err = readTo(fillItem, reader)
		fillItem.SetEndPos(reader)
		return
	}
}

func WrapReader(reader Reader, options *Options) (ret Reader) {
	if options.PositionFill {
		ret = NewReadToWrapperReader(reader, ReadToPositionWrapper)
	} else {
		ret = reader
	}
	return ret
}
