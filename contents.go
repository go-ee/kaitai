package kaitai

import (
	"bytes"
	"fmt"
)

type Contents struct {
	Name          string `-`
	ContentString *string
	ContentArray  []byte
	TypeSwitch    *TypeSwitch
}

func (o *Contents) BuildReader(attr *Attr, spec *Spec) (ret ItemReader, err error) {
	if o.ContentString != nil {
		ret = &ContentStringReader{Attr: attr, Accessor: o, value: *o.ContentString}
	} else if o.ContentArray != nil {
		ret = &ContentArrayReader{Attr: attr, Accessor: o, array: o.ContentArray}
	} else if o.TypeSwitch != nil {
		ret, err = o.TypeSwitch.BuildReader(attr, spec)
	} else {
		err = fmt.Errorf("contents(%v) read ELSE not implemented yet", o.Name)
	}

	return
}

func (o *Contents) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	if err = unmarshal(&o.ContentString); err != nil {
		if err = unmarshal(&o.ContentArray); err != nil {
			err = unmarshal(&o.TypeSwitch)
		}
	}
	return
}

type ContentStringReader struct {
	Attr     *Attr
	Accessor *Contents

	validate bool
	value    string
}

func (o *ContentStringReader) Read(reader ReaderIO, _ *Item, _ *Item) (ret *Item, err error) {
	var data []byte
	//each character as a byte
	if data, err = reader.ReadBytes(uint8(len(o.value))); err == nil {
		currentValue := string(data)
		if o.validate && currentValue != o.value {
			err = fmt.Errorf("content is different, '%v' != '%v'", currentValue, o.value)
		}
		if err == nil {
			ret = &Item{Attr: o.Attr, Accessor: o.Accessor, Value: currentValue}
		}
	}
	return
}

type ContentArrayReader struct {
	Attr     *Attr
	Accessor *Contents

	validate bool
	array    []byte
}

func (o *ContentArrayReader) Read(reader ReaderIO, _ *Item, _ *Item) (ret *Item, err error) {
	var data []byte
	//each character as a byte
	if data, err = reader.ReadBytes(uint8(len(o.array))); err == nil {
		if o.validate && bytes.Compare(data, o.array) != 0 {
			err = fmt.Errorf("content is different, '%v' != '%v'", data, o.array)
		}
		if err == nil {
			ret = &Item{Attr: o.Attr, Accessor: o.Accessor, Value: data}
		}
	}
	return
}
