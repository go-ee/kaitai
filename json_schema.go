package kaitai

import (
	"bytes"
)

func (o *Type) ToJsonSchema() (ret string, err error) {
	buffer := bytes.NewBufferString("")
	buffer.WriteString("\"$schema\": \"http://json-schema.org/draft-04/schema#\",")
	buffer.WriteString("\"definitions\": {")
	buffer.WriteString("},")
	err = o.toJsonSchemaDef(buffer)
	ret = buffer.String()
	return
}

func (o *Type) toJsonSchemaDef(buffer *bytes.Buffer) (err error) {
	buffer.WriteString("{")
	buffer.WriteString("\"type\": \"object\",")
	buffer.WriteString("\"properties\":{")
	for _, attr := range o.Seq {
		buffer.WriteString("\"")
		buffer.WriteString(attr.Id)
		buffer.WriteString("\":")
		buffer.WriteString("{")
		if err = attr.toJsonSchemaType(buffer); err != nil {
			return
		}
		buffer.WriteString("}")
	}

	buffer.WriteString("}")
	buffer.WriteString("}")
	return
}

func (o *Type) toJsonSchemaType() (ret string, err error) {
	ret = "\"$ref\": \"#/definitions/" + o.Id
	return
}

func (o *TypeRef) toJsonSchemaType() (ret string, err error) {
	if o.Type != nil {
		ret, err = o.Type.toJsonSchemaType()
	} else if o.Enum != nil {
		ret, err = o.Enum.toJsonSchemaType()
	} else if o.Instance != nil {
		ret, err = o.Instance.toJsonSchemaType()
	} else if o.TypeSwitch != nil {
		ret, err = o.TypeSwitch.toJsonSchemaType()
	} else if o.Native != nil {
		ret, err = o.Native.toJsonSchemaType()
	}
	return
}

func (o *Enum) toJsonSchemaType() (ret string, err error) {
	ret = "\"$ref\": \"#/definitions/" + o.Id
	return
}

func (o *Instance) toJsonSchemaType() (ret string, err error) {
	ret = "\"$ref\": \"#/definitions/" + o.Id
	return
}

func (o *TypeSwitch) toJsonSchemaType() (ret string, err error) {
	ret = "\"type\": \"object\"
	return
}

func (o *Native) toJsonSchemaType() (ret string, err error) {
	if o.Type == "str" || o.Type == "strz" {
		ret = "\"type\": \"string\"
	} else {
		switch o.Type {
		case "b":
			ret, err = o.toJsonSchemaTypeB()
		case "u":
			ret, err = o.toJsonSchemaTypeU()
		case "s":
			ret, err = o.toJsonSchemaTypeS()
		case "f":
			ret, err = o.toJsonSchemaTypeF()
		default:
			ret = "\"type\": \"string\"
		}
	}
	return
}

func (o *Native) toJsonSchemaTypeB() (ret string, err error) {
	switch o.Length {
	case 1:
		ret = "\"type\": \"bool\"
	case 2:
		ret = "\"type\": \"number\"
	default:
		ret = "\"type\": \"number\"
	}
	return
}

func (o *Native) toJsonSchemaTypeU() (ret string, err error) {
	switch o.Length {
	case 1:
		ret = "\"type\": \"bool\"
	case 2:
		ret = "\"type\": \"number\"
	default:
		ret = "\"type\": \"number\"
	}
	return
}

func (o *Native) toJsonSchemaTypeS() (ret string, err error) {
	switch o.Length {
	case 1:
		ret = "\"type\": \"bool\"
	case 2:
		ret = "\"type\": \"number\"
	default:
		ret = "\"type\": \"number\"
	}
	return
}

func (o *Native) toJsonSchemaTypeF() (ret string, err error) {
	switch o.Length {
	case 1:
		ret = "\"type\": \"bool\"
	case 2:
		ret = "\"type\": \"number\"
	default:
		ret = "\"type\": \"number\"
	}
	return
}

func (o *Attr) toJsonSchemaType(buffer *bytes.Buffer) (err error) {
	var itemType string
	if o.Type != nil {
		itemType, err = o.Type.toJsonSchemaType()
	} else if o.Contents != nil {
		itemType, err = o.Contents.toJsonSchemaType()
	} else if o.SizeEos == "true" {
		itemType = "\"type\": \"string\"
	} else {
		itemType = "\"type\": \"string\"
	}

	if o.Repeat == "eos" {
		itemType = "\"type\": \"array, \"items\": {" + itemType + "}"
	}
	return
}

func (o *Contents) toJsonSchemaType() (ret string, err error) {
	return
}

/*{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "definitions": {
    "pet": {
      "type": "object",
      "properties": {
        "name":  { "type": "string" },
        "breed": { "type": "string" },
        "age":  { "type": "string" }
      },
      "required": ["name", "breed", "age"]
    }
  },
  "type": "object",
  "properties": {
    "cat": { "$ref": "#/definitions/pet" },
    "dog": { "$ref": "#/definitions/pet" }
  }
}*/
