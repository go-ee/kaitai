package kaitai

import "fmt"

type TypeSwitch struct {
	SwitchOn string              `yaml:"switch-on,omitempty"`
	Cases    map[string]*TypeRef `yaml:"cases,omitempty"`
}

func (o *TypeSwitch) BuildReader(attr *Attr, spec *Spec) (ret Reader, err error) {
	typeSwitchReader := &TypeSwitchReader{
		ReaderBase:      &ReaderBase{attr: attr, accessor: o},
		findSwitchValue: o.buildSwitchValueFinder(),
	}
	if typeSwitchReader.cases, err = o.buildCaseReaders(attr, spec); err != nil {
		return
	}
	typeSwitchReader.defaultCase = typeSwitchReader.cases["_"]
	ret = typeSwitchReader
	return
}

func (o *TypeSwitch) buildSwitchValueFinder() func(attr *Attr, item *Item) (ret string, err error) {
	return func(attr *Attr, parent *Item) (ret string, err error) {
		var switchOnValue *Item
		if switchOnValue, err = parent.Expr(o.SwitchOn); err != nil {
			return
		}
		ret = fmt.Sprintf("%v", switchOnValue.Value())
		return
	}
}

func (o *TypeSwitch) buildCaseReaders(attr *Attr, spec *Spec) (ret map[string]Reader, err error) {
	ret = make(map[string]Reader, len(o.Cases))
	for name, caseItem := range o.Cases {
		var caseReader Reader
		if caseReader, err = caseItem.BuildReader(attr, spec); err != nil {
			return
		}
		ret[name] = caseReader
	}
	return
}

type TypeSwitchReader struct {
	*ReaderBase
	findSwitchValue func(attr *Attr, parent *Item) (string, error)
	cases           map[string]Reader
	defaultCase     Reader
}

func (o *TypeSwitchReader) ReadTo(fillItem *Item, reader *ReaderIO) (err error) {
	var switchValue string
	if switchValue, err = o.findSwitchValue(o.attr, fillItem.Parent); err != nil {
		return
	}

	itemReader := o.cases[switchValue]
	if itemReader == nil {
		itemReader = o.defaultCase
	}

	if itemReader != nil {
		err = itemReader.ReadTo(fillItem, reader)
	} else {
		err = fmt.Errorf("no case found for %v, %v", o.attr.Id, switchValue)
	}
	return
}
