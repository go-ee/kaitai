package kaitai

import (
	"fmt"
	"strconv"
	"strings"
)

type TypeRef struct {
	Name       string      `-`
	Type       *Type       `-`
	Enum       *Enum       `-`
	Instance   *Instance   `-`
	TypeSwitch *TypeSwitch `-`
	Native     *Native     `-`
}

func (o *TypeRef) BuildReader(attr *Attr, spec *Spec) (ret AttrReader, err error) {
	if o.Native != nil {
		ret, err = o.Native.BuildReader(attr, spec)
	} else if o.TypeSwitch != nil {
		ret, err = o.TypeSwitch.BuildReader(attr, spec)
	} else {
		if o.Type = spec.Types[o.Name]; o.Type == nil {
			if o.Enum = spec.Enums[o.Name]; o.Enum == nil {
				if o.Instance = spec.Instances[o.Name]; o.Instance != nil {
					ret, err = o.Instance.BuildReader(attr, spec)
				} else {
					err = fmt.Errorf("no accessor for attr(%v), typeRef(%v)", attr.Id, o.Name)
				}
			} else {
				ret, err = o.Enum.BuildReader(attr, spec)
			}
		} else {
			ret, err = o.Type.BuildReader(attr, spec)
		}
	}
	return
}

func (o *TypeRef) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	if err = unmarshal(&o.Name); err == nil {
		if strings.HasSuffix(o.Name, "str") {
			o.Native = &Native{Type: o.Name}
		} else {
			parts := buildInRegExp.FindStringSubmatch(o.Name)
			if parts != nil {
				length, _ := strconv.Atoi(parts[2])
				endian := parts[3]

				o.Native = &Native{Type: parts[1], Length: uint8(length), EndianBe: parseEndian(endian)}
			}
		}
	} else {
		err = unmarshal(&o.TypeSwitch)
	}
	return
}
