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

func pintRecordsLen(root *Item) {
	recordsValue, _ := root.ExprValue("records")

	if records, ok := recordsValue.([]interface{}); ok {
		println(len(records))
		//for _, record := range records {
		//item := record.(*Item)
		//value := item.value
		//println(value)
		//}
	} else {
		println(recordsValue)
	}
}

func item() (ret *Item) {
	ksySpecPath := os.Getenv("KSY_SPEC")
	ksyDataPath := os.Getenv("KSY_DATA")

	var err error
	m, err := NewModel(ksySpecPath, &Options{LazyDecoding: true, PositionFill: false})
	if err != nil {
		panic(err)
	}

	ret, err = m.Read(ksyDataPath)

	if err != nil {
		panic(err)
	}
	return
}
