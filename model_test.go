package kaitai

import (
	"os"
	"testing"
)

func TestParsing(t *testing.T) {
	ksySpecPath := os.Getenv("KSY_SPEC")
	ksyDataPath := os.Getenv("KSY_DATA")

	model, err := NewModel(ksySpecPath, &Options{LazyDecoding: false, PositionFill: false})
	if err != nil {
		panic(err)
	}

	item, err := model.Read(ksyDataPath)
	if err != nil {
		panic(err)
	}

	m := item.Value().(map[string]*Item)
	recordsItem := m["records"]
	records := recordsItem.Value().([]*Item)

	println(len(records))
}

func BenchmarkParsing(t *testing.B) {
	ksySpecPath := os.Getenv("KSY_SPEC")
	ksyDataPath := os.Getenv("KSY_DATA")

	model, err := NewModel(ksySpecPath, &Options{LazyDecoding: false, PositionFill: false})
	if err != nil {
		panic(err)
	}

	item, err := model.Read(ksyDataPath)
	if err != nil {
		panic(err)
	}

	m := item.Value().(map[string]*Item)
	recordsItem := m["records"]
	records := recordsItem.Value().([]*Item)

	println(len(records))
}
