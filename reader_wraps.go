package kaitai

type AttrParentRead struct {
	attr       *Attr
	parentRead ParentRead
}

func (o *AttrParentRead) Read(parent Item, reader *ReaderIO) (ret interface{}, err error) {
	return o.parentRead(parent, reader)
}

func (o *AttrParentRead) Attr() *Attr {
	return o.attr
}
