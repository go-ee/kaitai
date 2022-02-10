package kaitai

import (
	"github.com/sirupsen/logrus"
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
	//file, _ := json.Marshal(it)
	//_ = ioutil.WriteFile("file.json", file, 0644)
	logrus.Infof("end, %v", it["version"])
}

func item() (ret Item) {
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
