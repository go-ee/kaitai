package kaitai

import (
	"bytes"
	"fmt"
)

type Contents struct {
	Name          string `-`
	ContentString string
	ContentArray  []byte
	TypeSwitch    *TypeSwitch
}

func (o *Contents) BuildReader(attr *Attr, spec *Spec) (ret Reader, err error) {
	if o.ContentString != "" {
		ret = &ContentStringReader{
			ReaderBase: &ReaderBase{attr: attr, accessor: o}, value: o.ContentString, validate: true,
		}
	} else if o.ContentArray != nil {
		ret = &ContentArrayReader{
			ReaderBase: &ReaderBase{attr: attr, accessor: o}, array: o.ContentArray, validate: true,
		}
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
	*ReaderBase
	validate bool
	value    string
}

func (o *ContentStringReader) ReadTo(fillItem *Item, reader *ReaderIO) (err error) {
	//each character as a byte
	if fillItem.Raw, err = reader.ReadBytes(uint16(len(o.value))); err == nil {
		currentValue := string(fillItem.Raw)
		if o.validate && currentValue != o.value {
			err = fmt.Errorf("content is different, '%v' != '%v'", currentValue, o.value)
		}
		if err == nil {
			fillItem.SetValue(currentValue)
		}
	}
	return
}

type ContentArrayReader struct {
	*ReaderBase
	validate bool
	array    []byte
}

func (o *ContentArrayReader) ReadTo(fillItem *Item, reader *ReaderIO) (err error) {
	//each character as a byte
	if fillItem.Raw, err = reader.ReadBytes(uint16(len(o.array))); err == nil {
		if o.validate && bytes.Compare(fillItem.Raw, o.array) != 0 {
			err = fmt.Errorf("content is different, '%v' != '%v'", fillItem.Raw, o.array)
		}
		if err == nil {
			fillItem.SetValue(fillItem.Raw)
		}
	}
	return
}
