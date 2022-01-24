package kaitai

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"regexp"
	"strconv"
)

type Model struct {
	Root *Type
	Spec *Spec
}

func NewModelFromYamlFile(ksyPath string) (ret *Model, err error) {
	var data []byte
	if data, err = ioutil.ReadFile(ksyPath); err != nil {
		return
	}
	ret = &Model{Spec: &Spec{}}
	if err = yaml.Unmarshal(data, ret.Spec); err != nil {
		return
	}
	ret.ResolveReferences()
	return
}

func (o *Model) Info() string {
	return fmt.Sprintf("%v", o.Root)
}

func (o *Model) ResolveReferences() {
	o.Root = &Type{Seq: o.Spec.Seq, Doc: o.Spec.Doc}
	o.Spec.resolveReferences()
}

type Spec struct {
	Meta      *Meta                `yaml:"meta,omitempty"`
	Types     map[string]*Type     `yaml:"types,omitempty"`
	Seq       []*Attr              `yaml:"seq,omitempty"`
	Enums     map[string]*Enum     `yaml:"enums,omitempty"`
	Doc       string               `yaml:"doc,omitempty"`
	Instances map[string]*Instance `yaml:"instances,omitempty"`
}

func (o *Spec) resolveReferences() {
	for _, attr := range o.Seq {
		attr.resolveReferences(o)
	}

	for _, t := range o.Types {
		t.resolveReferences(o)
	}

	for _, instance := range o.Instances {
		instance.resolveReferences(o)
	}
	return
}

type Type struct {
	Seq []*Attr `yaml:"seq,omitempty"`
	Doc string  `yaml:"doc,omitempty"`
}

func (o *Type) resolveReferences(base *Spec) {
	for _, attr := range o.Seq {
		attr.resolveReferences(base)
	}
	return
}

func (o *Type) ReferencesResolved() (ret bool) {
	ret = true
	for _, c := range o.Seq {
		ret = c.ReferencesResolved()
		if !ret {
			break
		}
	}
	return
}

type Meta struct {
	Id            string `yaml:"id,omitempty"`
	Title         string `yaml:"title,omitempty"`
	Application   string `yaml:"application,omitempty"`
	Imports       string `yaml:"imports,omitempty"`
	Encoding      string `yaml:"encoding,omitempty"`
	Endian        string `yaml:"endian,omitempty"`
	KsVersion     string `yaml:"ks-version,omitempty"`
	KsDebug       string `yaml:"ks-debug,omitempty"`
	KsOpaqueTypes string `yaml:"ksopaquetypes,omitempty"`
	Licence       string `yaml:"licence,omitempty"`
	FileExtension string `yaml:"fileextension,omitempty"`
}

type Enum struct {
	Literals map[int]*Literal
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

type Instance struct {
	Attrs map[int]*Attr
}

func (o *Instance) resolveReferences(base *Spec) {
	for _, attr := range o.Attrs {
		attr.resolveReferences(base)
	}
	return
}

func (o *Instance) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	err = unmarshal(&o.Attrs)
	return
}

type Attr struct {
	Id          string    `yaml:"id,omitempty"`
	Type        *TypeRef  `yaml:"type"`
	Size        string    `yaml:"size,omitempty"`
	SizeEos     string    `yaml:"size-eos,omitempty"`
	Doc         string    `yaml:"doc,omitempty"`
	Repeat      string    `yaml:"repeat,omitempty"`
	RepeatExpr  string    `yaml:"repeat-expr,omitempty"`
	RepeatUntil string    `yaml:"repeat-until,omitempty"`
	Contents    *Contents `yaml:"contents,omitempty"`
	Value       string    `yaml:"value,omitempty"`
	Pos         string    `yaml:"pos,omitempty"`
	Whence      string    `yaml:"whence,omitempty"`
	Enum        string    `yaml:"enum,omitempty"`
	If          string    `yaml:"if,omitempty"`
	Process     string    `yaml:"process,omitempty"`
	Terminator  string    `yaml:"terminator,omitempty"`
	Consume     string    `yaml:"consume,omitempty"`
	Include     string    `yaml:"include,omitempty"`
	EosError    string    `yaml:"eos-error,omitempty"`
	PadRight    string    `yaml:"pad-right,omitempty"`
	Encoding    string    `yaml:"encoding,omitempty"`
}

func (o *Attr) resolveReferences(base *Spec) {
	if o.Type != nil {
		o.Type.resolveReferences(base)
	}
	return
}

func (o *Attr) ReferencesResolved() (ret bool) {
	ret = o.Type == nil || o.Type.ReferencesResolved()
	return
}

type Contents struct {
	Values []interface{}
	Switch *TypeSwitch
}

func (o *Contents) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	var value string
	if err = unmarshal(&value); err != nil {
		if err = unmarshal(&o.Values); err != nil {
			err = unmarshal(&o.Switch)
		}
	} else {
		o.Values = append(o.Values, value)
	}
	return
}

func (o *Contents) Len() (ret int) {
	ret = len(o.Values)
	if ret > 0 {
		switch v := o.Values[0].(type) {
		case string:
			return len(v)
		}
	}
	return
}

type TypeSwitch struct {
	SwitchOn string              `yaml:"switch-on,omitempty"`
	Cases    map[string]*TypeRef `yaml:"cases,omitempty"`
}

func (o *TypeSwitch) resolveReferences(base *Spec) {
	for _, t := range o.Cases {
		t.resolveReferences(base)
	}
	return
}

func (o *TypeSwitch) ReferencesResolved() (ret bool) {
	ret = true
	for _, c := range o.Cases {
		ret = c.ReferencesResolved()
		if !ret {
			break
		}
	}
	return
}

type BuildIn struct {
	Type        string
	BytesLength int
	EndianBe    *bool
}

var buildInRegExp *regexp.Regexp

func init() {
	buildInRegExp = regexp.MustCompile(`(strz|str|b|f|s|u)([1-8])(be|le|)`)
}

type TypeRef struct {
	Name       string
	Type       *Type
	Enum       *Enum
	Instance   *Instance
	TypeSwitch *TypeSwitch
	BuildIn    *BuildIn
}

func (o *TypeRef) ReferencesResolved() (ret bool) {
	ret = o.BuildIn != nil || o.Type != nil || o.Enum != nil || o.Instance != nil || o.TypeSwitch.ReferencesResolved()
	return
}

func (o *TypeRef) resolveReferences(base *Spec) {
	if o.BuildIn == nil {
		if o.Type = base.Types[o.Name]; o.Type == nil {
			if o.Enum = base.Enums[o.Name]; o.Enum == nil {
				o.Instance = base.Instances[o.Name]
			}
		}
	} else if o.TypeSwitch != nil {
		o.TypeSwitch.resolveReferences(base)
	}
}

func (o *TypeRef) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	if err = unmarshal(&o.Name); err == nil {
		parts := buildInRegExp.FindStringSubmatch(o.Name)

		if parts != nil {
			length, _ := strconv.Atoi(parts[2])
			var endianBe bool
			endian := parts[3]
			if endian == "be" {
				endianBe = true
			} else if endian == "be" {
				endianBe = false
			}
			o.BuildIn = &BuildIn{Type: parts[1], BytesLength: length, EndianBe: &endianBe}
		}
	} else {
		err = unmarshal(&o.TypeSwitch)
	}
	return
}
