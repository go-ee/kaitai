package kaitai

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io/ioutil"
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
	logrus.Infof("start")
	it := item()
	file, _ := json.Marshal(it)
	_ = ioutil.WriteFile(it.Type.Id+".json", file, 0644)
	logrus.Infof("end")
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
