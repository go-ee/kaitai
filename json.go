package kaitai

import (
	"bytes"
	"encoding/json"
)

func (o *TypeItem) ToJsonIndent(prefix, indent string) (ret []byte, err error) {
	if ret, err = o.ToJson(); err == nil {
		var buf bytes.Buffer
		if err = json.Indent(&buf, ret, prefix, indent); err == nil {
			ret = buf.Bytes()
		}
	}
	return
}

func (o *TypeItem) ToJson() (ret []byte, err error) {
	buffer := bytes.NewBufferString("")
	err = o.FillJson(buffer)
	ret = buffer.Bytes()
	return
}

func (o *TypeItem) FillJson(buffer *bytes.Buffer) (err error) {
	buffer.WriteString("{")
	for i, value := range o.attrs {
		o.toJsonAttrName(buffer, i)
		if childItem, ok := value.(*TypeItem); ok {
			if err = childItem.FillJson(buffer); err != nil {
				break
			}
		} else {
			switch t := value.(type) {
			case []interface{}:
				buffer.WriteString("[")
				for _, arrayItem := range t {
					if arrayTypeItem, arrayOk := arrayItem.(*TypeItem); arrayOk {
						if err = arrayTypeItem.FillJson(buffer); err != nil {
							break
						}
					} else {
						if err = o.toJsonNative(buffer, value); err != nil {
							break
						}
					}
					buffer.WriteString(",")
				}
				buffer.Truncate(buffer.Len() - 1)
				buffer.WriteString("]")
			case []*TypeItem:
				buffer.WriteString("[")
				for _, arrayItem := range t {
					if err = arrayItem.FillJson(buffer); err != nil {
						break
					}
					buffer.WriteString(",")
				}
				buffer.Truncate(buffer.Len() - 1)
				buffer.WriteString("]")
			default:
				if err = o.toJsonNative(buffer, value); err != nil {
					break
				}
			}
		}
		buffer.WriteString(",")
	}

	if err == nil {
		buffer.Truncate(buffer.Len() - 1)
		buffer.WriteString("}")
	}
	return
}

func (o *TypeItem) toJsonAttrName(buffer *bytes.Buffer, i int) {
	buffer.WriteString("\"")
	buffer.WriteString(o.model.IndexToAttrName(i))
	buffer.WriteString("\":")
}

func (o *TypeItem) toJsonNative(buffer *bytes.Buffer, value interface{}) (err error) {
	var jsonValue []byte
	if jsonValue, err = json.Marshal(value); err == nil {
		buffer.Write(jsonValue)
	}
	return
}

func (o *TypeItem) MarshalJSON() (ret []byte, err error) {
	buffer := bytes.NewBufferString("")
	if err = o.FillJson(buffer); err == nil {
		ret = buffer.Bytes()
	}
	return
}
