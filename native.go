package kaitai

import (
	"fmt"
)

type Native struct {
	Type     string
	Length   uint8
	EndianBe *bool
}

func (o *Native) BuildReader(attr *Attr, spec *Spec) (ret ItemReader, err error) {
	if o.EndianBe == nil {
		o.EndianBe = spec.Meta.EndianBe
	}
	var endianConverter EndianReader
	if *o.EndianBe {
		endianConverter = BigEndianConverter
	} else {
		endianConverter = LittleEndianConverter
	}

	var fix ReadFix
	var dynamic ReadDynamic

	switch o.Type {
	case "str":
		fix = ReadFixAttr(attr, ToString)
		dynamic = ReadDynamicAttr(attr, ToString)
	case "strz":
		fix = ReadFixAttr(attr, ToString)
		dynamic = ReadDynamicAttr(attr, ToString)
	case "b":
		fix = ReadB(endianConverter, o.Length)
	case "u":
		fix = ReadU(endianConverter, o.Length)
	case "s":
		fix = ReadS(endianConverter, o.Length)
	case "f":
		fix = ReadF(endianConverter, o.Length)
	default:
		err = fmt.Errorf("not supported Native(%v,%v)", o.Type, o.Length)
	}

	if fix != nil {
		ret = &NativeReaderFix{Attr: attr, Native: o, fix: fix}
	} else if dynamic != nil {
		ret = &NativeReaderDynamic{Attr: attr, Native: o, dynamic: dynamic}
	}
	return
}

type NativeReaderFix struct {
	Attr   *Attr
	Native *Native

	fix ReadFix
}

func (o *NativeReaderFix) Read(reader ReaderIO, _ *Item, _ *Item) (ret *Item, err error) {
	if value, currentErr := o.fix(reader); currentErr == nil {
		ret = &Item{Attr: o.Attr, Accessor: o.Native, Value: value}
	} else {
		err = currentErr
	}
	return
}

type NativeReaderDynamic struct {
	Attr   *Attr
	Native *Native

	dynamic ReadDynamic
}

func (o *NativeReaderDynamic) Read(reader ReaderIO, parent *Item, root *Item) (ret *Item, err error) {
	if value, currentErr := o.dynamic(reader, parent, root); currentErr == nil {
		ret = &Item{Attr: o.Attr, Accessor: o.Native, Value: value}
	} else {
		err = currentErr
	}
	return
}
