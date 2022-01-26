package kaitai

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type Model struct {
	Root       *Type
	Spec       *Spec
	itemReader ItemReader
}

func NewModel(ksyPath string) (ret *Model, err error) {
	var data []byte
	if data, err = ioutil.ReadFile(ksyPath); err != nil {
		return
	}
	ret = &Model{Spec: &Spec{}}
	if err = yaml.Unmarshal(data, ret.Spec); err != nil {
		return
	}
	err = ret.crossInit()
	return
}

func (o *Model) Info() string {
	return fmt.Sprintf("%v", o.Root)
}

func (o *Model) crossInit() (err error) {
	o.Root = &Type{Id: o.Spec.Meta.Id, Seq: o.Spec.Seq, Doc: o.Spec.Doc}
	if err = o.Spec.crossInit(); err != nil {
		return
	}
	o.itemReader, err = o.Root.BuildReader(nil, o.Spec)
	return
}

func (o *Model) Read(filePath string) (ret *Item, err error) {
	var file *os.File
	if file, err = os.Open(filePath); err != nil {
		return
	}
	ret, err = o.itemReader.Read(ReaderIO{ReadSeeker: file}, nil, nil)
	return
}

type Spec struct {
	Meta      *Meta                `yaml:"meta,omitempty"`
	Types     map[string]*Type     `yaml:"types,omitempty"`
	Seq       []*Attr              `yaml:"seq,omitempty"`
	Enums     map[string]*Enum     `yaml:"enums,omitempty"`
	Doc       string               `yaml:"doc,omitempty"`
	Instances map[string]*Instance `yaml:"instances,omitempty"`
}

func (o *Spec) crossInit() (err error) {
	o.Meta.crossInit()

	for _, item := range o.Seq {
		if err = item.crossInit(o); err != nil {
			return
		}
	}

	for id, item := range o.Types {
		item.Id = id
		if err = item.crossInit(o); err != nil {
			return
		}
	}

	for id, item := range o.Enums {
		item.Id = id
	}

	for id, item := range o.Instances {
		item.Id = id
		if err = item.crossInit(o); err != nil {
			return
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

func (o *Meta) crossInit() {
	o.EndianBe = parseEndian(o.Endian)
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
