package kaitai

import "fmt"

type Instance struct {
	Id    string `-`
	Attrs map[int]*Attr
}

func (o *Instance) BuildReader(attr *Attr, spec *Spec) (ret ItemReader, err error) {
	err = fmt.Errorf("read %v.Instance(%v) not implemented yet", attr.Id, o.Id)
	return
}

func (o *Instance) crossInit(base *Spec) (err error) {
	for _, attr := range o.Attrs {
		if err = attr.crossInit(base); err != nil {
			return
		}
	}
	return
}

func (o *Instance) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	err = unmarshal(&o.Attrs)
	return
}
