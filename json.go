package kaitai

import (
	"bytes"
	"encoding/json"
)

func (o *TypeItem) ToJson() (ret []byte, err error) {
	buffer := bytes.NewBufferString("")
	err = o.FillJson(buffer)
	ret = buffer.Bytes()
	return
}

func (o *TypeItem) FillJson(buffer *bytes.Buffer) (err error) {
	buffer.WriteString("{")
	for i, value := range o.attrs {
		buffer.WriteString("\"")
		buffer.WriteString(o.model.IndexToAttrName(i))
		buffer.WriteString("\":")
		if childItem, ok := value.(*TypeItem); ok {
			if err = childItem.FillJson(buffer); err != nil {
				break
			}
		} else {
			switch t := value.(type) {
			case []interface{}:
				buffer.WriteString("[")
				for _, arrayItem := range t {
					if arrayItem, arrayOk := arrayItem.(*TypeItem); arrayOk {
						if err = arrayItem.FillJson(buffer); err != nil {
							break
						}
					} else {
						if jsonValue, attErr := json.Marshal(value); attErr == nil {
							buffer.Write(jsonValue)
						} else {
							err = attErr
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
				if jsonValue, attErr := json.Marshal(value); attErr == nil {
					buffer.Write(jsonValue)
				} else {
					err = attErr
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

func (o *TypeItem) MarshalJSON() (ret []byte, err error) {
	buffer := bytes.NewBufferString("")
	if err = o.FillJson(buffer); err == nil {
		ret = buffer.Bytes()
	}
	return
}
