package kaitai

type ReadToWrapper func(readTo ReadTo) ReadTo

type AttrAccessorReadToReader struct {
	*ReaderBase
	readTo ReadTo
}

func (o *AttrAccessorReadToReader) ReadTo(fillItem *Item, reader *ReaderIO) (err error) {
	return o.readTo(fillItem, reader)
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
