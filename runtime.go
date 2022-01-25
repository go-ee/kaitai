package kaitai

import (
	"fmt"
	"strings"
)

type Item struct {
	Type  interface{}
	Value interface{}
}

func (o *Item) Expr(expr string) (ret *Item, err error) {
	ret = o
	for _, part := range strings.Split(expr, ".") {
		if value, ok := o.Value.(map[string]*Item); ok {
			if ret = value[part]; ret == nil {
				err = fmt.Errorf("can't resolve '%v' of expression '%v'", part, expr)
				break
			}
		} else {
			err = fmt.Errorf("can't resolve '%v' of expression '%v'", part, expr)
			break
		}
	}
	return
}
