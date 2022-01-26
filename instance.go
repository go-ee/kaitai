package kaitai

import "fmt"

type Instance struct {
	Id    string `-`
	Attrs map[int]*Attr
}

func (o *Instance) BuildReader(attr *Attr, spec *Spec) (ret AttrReader, err error) {
	err = fmt.Errorf("read %v.Instance(%v) not implemented yet", attr.Id, o.Id)
	return
}

func (o *Instance) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	err = unmarshal(&o.Attrs)
	return
}
