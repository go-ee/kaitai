package kaitai

import "fmt"

type TypeSwitch struct {
	SwitchOn string              `yaml:"switch-on,omitempty"`
	Cases    map[string]*TypeRef `yaml:"cases,omitempty"`
}

func (o *TypeSwitch) BuildReader(attr *Attr, spec *Spec) (ret AttrReader, err error) {
	typeSwitchReader := &TypeSwitchReader{
		AttrReaderBase:  &AttrReaderBase{attr, o},
		findSwitchValue: o.buildSwitchValueFinder(),
	}
	if typeSwitchReader.cases, err = o.buildCaseReaders(attr, spec); err != nil {
		return
	}
	typeSwitchReader.defaultCase = typeSwitchReader.cases["_"]
	ret = typeSwitchReader
	return
}

func (o *TypeSwitch) buildSwitchValueFinder() func(attr *Attr, parent *Item, root *Item) (ret string, err error) {
	return func(attr *Attr, parent *Item, root *Item) (ret string, err error) {
		var switchOnValue *Item
		if switchOnValue, err = parent.Expr(o.SwitchOn); err != nil {
			return
		}
		ret = fmt.Sprintf("%v", switchOnValue.Value)
		return
	}
}

func (o *TypeSwitch) buildCaseReaders(attr *Attr, spec *Spec) (ret map[string]AttrReader, err error) {
	ret = make(map[string]AttrReader, len(o.Cases))
	for _, caseItem := range o.Cases {
		var caseReader AttrReader
		if caseReader, err = caseItem.BuildReader(attr, spec); err != nil {
			return
		}
		ret[caseItem.Name] = caseReader
	}
	return
}

type TypeSwitchReader struct {
	*AttrReaderBase

	findSwitchValue func(attr *Attr, parent *Item, root *Item) (string, error)
	cases           map[string]AttrReader
	defaultCase     AttrReader
}

func (o *TypeSwitchReader) Read(reader Reader, parent *Item, root *Item) (ret *Item, err error) {
	var switchValue string
	if switchValue, err = o.findSwitchValue(o.attr, parent, root); err != nil {
		return
	}

	itemReader := o.cases[switchValue]
	if itemReader != nil {
		itemReader = o.defaultCase
	}

	if itemReader != nil {
		ret, err = itemReader.Read(reader, parent, root)
	} else {
		err = fmt.Errorf("no case found for %v, %v", o.attr.Id, switchValue)
	}
	return
}
