package kaitai

type AttrParentRead struct {
	attr       *Attr
	parentRead ParentRead
}

func (o *AttrParentRead) Read(parent *Item, reader *ReaderIO) (ret interface{}, err error) {
	return o.parentRead(parent, reader)
}

func (o *AttrParentRead) Attr() *Attr {
	return o.attr
}

func ItemReadToPositionWrapper(itemRead ItemRead) ItemRead {
	return func(parent *Item, reader *ReaderIO) (ret *Item, err error) {
		startPos := reader.Position()
		if ret, err = itemRead(parent, reader); err == nil {
			ret.StartPos = &startPos
			ret.SetEndPos(reader)
		}
		return
	}
}

func WrapReader(reader ItemRead, options *Options) (ret ItemRead) {
	if options.PositionFill {
		//ret = NewReadToWrapperReader(reader, ItemReadToPositionWrapper)
	} else {
		ret = reader
	}
	return ret
}

func ReadToParentRead(read Read) ParentRead {
	return func(parent *Item, reader *ReaderIO) (ret interface{}, err error) {
		return read(reader)
	}
}
