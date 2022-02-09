package kaitai

import (
	"encoding/json"
)

func (o *Item) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.Value())
}
