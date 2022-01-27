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

func (o *Contents) BuildReader(attr *Attr, spec *Spec) (ret AttrReader, err error) {
	if o.ContentString != "" {
		ret = &ContentStringReader{
			AttrReaderBase: &AttrReaderBase{attr, o},
			value:          o.ContentString,
			validate:       true,
		}
	} else if o.ContentArray != nil {
		ret = &ContentArrayReader{
			AttrReaderBase: &AttrReaderBase{attr, o},
			array:          o.ContentArray,
			validate:       true,
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
	*AttrReaderBase
	validate bool
	value    string
}

func (o *ContentStringReader) ReadTo(fillItem *Item, reader *Reader) (err error) {

	fillItem.SetStartPos(reader)

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

	fillItem.SetEndPos(reader)

	return
}

type ContentArrayReader struct {
	*AttrReaderBase
	validate bool
	array    []byte
}

func (o *ContentArrayReader) ReadTo(fillItem *Item, reader *Reader) (err error) {

	fillItem.SetStartPos(reader)

	//each character as a byte
	if fillItem.Raw, err = reader.ReadBytes(uint16(len(o.array))); err == nil {
		if o.validate && bytes.Compare(fillItem.Raw, o.array) != 0 {
			err = fmt.Errorf("content is different, '%v' != '%v'", fillItem.Raw, o.array)
		}
		if err == nil {
			fillItem.SetValue(fillItem.Raw)
		}
	}

	fillItem.SetEndPos(reader)

	return
}
