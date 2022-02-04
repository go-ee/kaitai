package kaitai

import (
	"os"
	"testing"
)

func TestParsing(t *testing.T) {
	toJson()
}

func BenchmarkParsing(t *testing.B) {
	toJson()
}

func toJson() {
	it := item()
	pintRecordsLen(it)
}

func pintRecordsLen(item *Item) {
	m := item.Value().(map[string]*Item)
	recordsItem := m["records"]
	records := recordsItem.Value().([]*Item)
	println(len(records))
}

func item() (ret *Item) {
	ksySpecPath := os.Getenv("KSY_SPEC")
	ksyDataPath := os.Getenv("KSY_DATA")

	var err error
	m, err := NewModel(ksySpecPath, &Options{LazyDecoding: false, PositionFill: false})
	if err != nil {
		panic(err)
	}

	ret, err = m.Read(ksyDataPath)

	if err != nil {
		panic(err)
	}
	return
}
