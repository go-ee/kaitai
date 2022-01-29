package kaitai

import "fmt"

type Enum struct {
	Id       string `-`
	Literals map[int]*Literal
}

func (o *Enum) BuildReader(attr *Attr, spec *Spec) (ret Reader, err error) {
	err = fmt.Errorf("read %v.Enum(%v) not implemented yet", attr.Id, o.Id)
	return
}

func (o *Enum) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	err = unmarshal(&o.Literals)
	return
}

type Literal struct {
	Id  string `yaml:"id,omitempty"`
	Doc string `yaml:"doc,omitempty"`
}

func (o *Literal) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	var lit literal
	if err = unmarshal(&lit); err != nil {
		err = unmarshal(&o.Id)
	} else {
		o.Id = lit.Id
		o.Doc = lit.Doc
	}
	return
}

type literal struct {
	Id  string `yaml:"id,omitempty"`
	Doc string `yaml:"doc,omitempty"`
}
