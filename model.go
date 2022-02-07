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
	itemReader *TypeReader
}

func NewModel(ksyPath string, options *Options) (ret *Model, err error) {
	var data []byte
	if data, err = ioutil.ReadFile(ksyPath); err != nil {
		return
	}

	if options == nil {
		options = &Options{}
	}

	ret = &Model{Spec: &Spec{Options: options}}
	if err = yaml.Unmarshal(data, ret.Spec); err != nil {
		return
	}

	err = ret.build()
	return
}

func (o *Model) Info() string {
	return fmt.Sprintf("%v", o.Root)
}

func (o *Model) build() (err error) {
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
	defer file.Close()

	var item interface{}
	if item, err = o.itemReader.Read(nil, &ReaderIO{ReadSeeker: file}); err == nil {
		ret = item.(*Item)
	}
	return
}

type Spec struct {
	Meta      *Meta                `yaml:"meta,omitempty"`
	Types     map[string]*Type     `yaml:"types,omitempty"`
	Seq       []*Attr              `yaml:"seq,omitempty"`
	Enums     map[string]*Enum     `yaml:"enums,omitempty"`
	Doc       string               `yaml:"doc,omitempty"`
	Instances map[string]*Instance `yaml:"instances,omitempty"`
	Options   *Options             `-`
}

func (o *Spec) crossInit() (err error) {
	o.Meta.crossInit()

	for id, item := range o.Types {
		item.Id = id
	}

	for id, item := range o.Enums {
		item.Id = id
	}

	for id, item := range o.Instances {
		item.Id = id
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

type Options struct {
	LazyDecoding bool
	PositionFill bool
}
