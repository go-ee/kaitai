package kaitai

import (
	"fmt"
	"github.com/kaitai-io/kaitai_struct_go_runtime/kaitai"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
)

type Model struct {
	Root *Type
	Spec *Spec
}

func ParseToModelFromYamlFile(ksyPath string) (ret *Model, err error) {
	var data []byte
	if data, err = ioutil.ReadFile(ksyPath); err != nil {
		return
	}
	ret = &Model{Spec: &Spec{}}
	if err = yaml.Unmarshal(data, ret.Spec); err != nil {
		return
	}
	ret.resolveRefs()
	return
}

func (o *Model) Read(filePath string) (ret *Item, err error) {
	var file *os.File
	if file, err = os.Open(filePath); err != nil {
		return
	}
	ret, err = o.Root.Read(kaitai.Stream{ReadSeeker: file}, nil, nil)
	return
}

func (o *Model) Info() string {
	return fmt.Sprintf("%v", o.Root)
}

func (o *Model) resolveRefs() {
	o.Root = &Type{Id: o.Spec.Meta.Id, Seq: o.Spec.Seq, Doc: o.Spec.Doc}
	o.Spec.resolveRefs()
}

type Spec struct {
	Meta      *Meta                `yaml:"meta,omitempty"`
	Types     map[string]*Type     `yaml:"types,omitempty"`
	Seq       []*Attr              `yaml:"seq,omitempty"`
	Enums     map[string]*Enum     `yaml:"enums,omitempty"`
	Doc       string               `yaml:"doc,omitempty"`
	Instances map[string]*Instance `yaml:"instances,omitempty"`
}

func (o *Spec) resolveRefs() {
	o.Meta.resolveRefs()

	for _, item := range o.Seq {
		item.resolveRefs(o)
	}

	for id, item := range o.Types {
		item.Id = id
		item.resolveRefs(o)
	}

	for id, item := range o.Enums {
		item.Id = id
	}

	for id, item := range o.Instances {
		item.Id = id
		item.resolveRefs(o)
	}
	return
}

type Type struct {
	Id  string  `-`
	Seq []*Attr `yaml:"seq,omitempty"`
	Doc string  `yaml:"doc,omitempty"`
}

func (o *Type) Read(stream kaitai.Stream, parent *Item, root *Item) (ret *Item, err error) {
	data := map[string]interface{}{}
	ret = &Item{Type: o, Value: data}

	parent = ret
	if root == nil {
		root = ret
	}

	for _, attr := range o.Seq {
		var item *Item
		if item, err = attr.Read(stream, parent, root); err != nil {
			break
		}
		data[attr.Id] = item
	}
	if err != nil {
		ret = nil
	}
	return
}

func (o *Type) resolveRefs(base *Spec) {
	for _, item := range o.Seq {
		item.resolveRefs(base)
	}
	return
}

func (o *Type) RefsResolved() (ret bool) {
	ret = true
	for _, c := range o.Seq {
		ret = c.RefsResolved()
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
	EndianBe      *bool  `-`
}

func (o *Meta) resolveRefs() {
	o.EndianBe = parseEndian(o.Endian)
}

type Enum struct {
	Id       string `-`
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
	Id    string `-`
	Attrs map[int]*Attr
}

func (o *Instance) resolveRefs(base *Spec) {
	for _, attr := range o.Attrs {
		attr.resolveRefs(base)
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

func (o *Attr) Read(stream kaitai.Stream, parent *Item, root *Item) (ret *Item, err error) {
	if o.Repeat == "eos" {
		var items []*Item
		for i := 0; err == nil; i++ {
			var item *Item
			if item, err = o.readSingle(stream, parent, root); err == nil {
				items = append(items, item)
			}
		}

		if io.EOF == err {
			err = nil
			ret = &Item{Type: o, Value: items}
		}
	} else {
		ret, err = o.readSingle(stream, parent, root)
	}
	return
}

func (o *Attr) readSingle(stream kaitai.Stream, parent *Item, root *Item) (ret *Item, err error) {
	if o.Type != nil {
		ret, err = o.Type.Read(stream, o, parent, root)
	} else if o.Contents != nil {
		ret, err = o.Contents.Read(stream, parent, root)
	} else {
		err = fmt.Errorf("attr(%v) ELSE, not implemented yet", o.Id)
	}
	return
}

func (o *Attr) resolveRefs(base *Spec) {
	if o.Type != nil {
		o.Type.resolveRefs(base)
	}
	return
}

func (o *Attr) RefsResolved() (ret bool) {
	ret = o.Type == nil || o.Type.RefsResolved()
	return
}

type Contents struct {
	Name   string `-`
	Values []interface{}
	Switch *TypeSwitch
}

func (o *Contents) Read(stream kaitai.Stream, parent *Item, root *Item) (ret *Item, err error) {
	if o.Values != nil {
		var data []byte
		if data, err = stream.ReadBytes(len(o.Values)); err == nil {
			ret = &Item{Type: o, Value: data}
		}
	} else if o.Switch != nil {
		err = fmt.Errorf("contents(%v) read Switch not implemented yet", o.Name)
	} else {
		err = fmt.Errorf("contents(%v) read ELSE not implemented yet", o.Name)
	}
	return
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

func (o *TypeSwitch) Read(stream kaitai.Stream, attr *Attr, parent *Item, root *Item) (ret *Item, err error) {
	var switchOnValue *Item
	if switchOnValue, err = parent.Expr(o.SwitchOn); err != nil {
		return
	}

	if value, ok := switchOnValue.Value.(string); ok {
		if typeRef := o.Cases[value]; typeRef != nil {
			ret, err = typeRef.Read(stream, attr, parent, root)
		} else {
			err = fmt.Errorf("can't find SwitchOn %v", value)
		}
	} else {
		err = fmt.Errorf("can't find SwitchOn %v", value)
	}
	return
}

func (o *TypeSwitch) resolveRefs(base *Spec) {
	for _, t := range o.Cases {
		t.resolveRefs(base)
	}
	return
}

func (o *TypeSwitch) RefsResolved() (ret bool) {
	ret = true
	for _, c := range o.Cases {
		ret = c.RefsResolved()
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

func (o *BuildIn) Read(stream kaitai.Stream, attr *Attr, parent *Item, root *Item) (ret *Item, err error) {
	if o.BytesLength > 0 {
		data := make([]byte, o.BytesLength)
		if _, err = stream.Read(data); err == nil {
			ret = &Item{Type: o, Value: data}
		}
	} else {
		var data []byte
		if data, err = ioutil.ReadAll(stream); err == nil {
			ret = &Item{Type: o, Value: data}
		}
	}
	return
}

var buildInRegExp *regexp.Regexp

func init() {
	buildInRegExp = regexp.MustCompile(`(b|f|s|u)([1-8])(be|le|)`)
}

type TypeRef struct {
	Name       string
	Type       *Type
	Enum       *Enum
	Instance   *Instance
	TypeSwitch *TypeSwitch
	BuildIn    *BuildIn
}

func (o *TypeRef) Read(stream kaitai.Stream, attr *Attr, parent *Item, root *Item) (ret *Item, err error) {
	if o.BuildIn != nil {
		ret, err = o.BuildIn.Read(stream, attr, parent, root)
	} else if o.Enum != nil {
		err = fmt.Errorf("TypeRef(%v) read Enum not implemented yet, %v", o.Name, o.Instance.Id)
	} else if o.Type != nil {
		ret, err = o.Type.Read(stream, parent, root)
	} else if o.TypeSwitch != nil {
		ret, err = o.TypeSwitch.Read(stream, attr, parent, root)
	} else if o.Instance != nil {
		err = fmt.Errorf("TypeRef(%v) read Instance not implemented yet, %v", o.Name, o.Instance.Id)
	} else {
		err = fmt.Errorf("TypeRef(%v) not resolved", o.Name)
	}
	return
}

func (o *TypeRef) RefsResolved() (ret bool) {
	ret = o.BuildIn != nil || o.Type != nil || o.Enum != nil || o.Instance != nil || o.TypeSwitch.RefsResolved()
	return
}

func (o *TypeRef) resolveRefs(base *Spec) {
	if o.BuildIn != nil {
		if o.BuildIn.EndianBe == nil {
			o.BuildIn.EndianBe = base.Meta.EndianBe
		}
	} else if o.TypeSwitch != nil {
		o.TypeSwitch.resolveRefs(base)
	} else {
		if o.Type = base.Types[o.Name]; o.Type == nil {
			if o.Enum = base.Enums[o.Name]; o.Enum == nil {
				o.Instance = base.Instances[o.Name]
			}
		}
	}
}

func (o *TypeRef) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	if err = unmarshal(&o.Name); err == nil {
		if o.Name == "str" || o.Name == "strz" {
			o.BuildIn = &BuildIn{Type: o.Name}
		} else {
			parts := buildInRegExp.FindStringSubmatch(o.Name)
			if parts != nil {
				length, _ := strconv.Atoi(parts[2])
				endian := parts[3]
				o.BuildIn = &BuildIn{Type: parts[1], BytesLength: length, EndianBe: parseEndian(endian)}
			}
		}
	} else {
		err = unmarshal(&o.TypeSwitch)
	}
	return
}

func parseEndian(endian string) (ret *bool) {
	if endian == "be" {
		endianBe := true
		ret = &endianBe
	} else if endian == "le" {
		endianBe := false
		ret = &endianBe
	}
	return
}
