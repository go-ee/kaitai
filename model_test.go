package kaitai

import (
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
	data, err := it.ToJson()
	//data, err := it.ToJsonIndent("", " ")
	//data, err := json.MarshalIndent(it, "", " ")
	if err == nil {
		_ = ioutil.WriteFile(it.model.typeSpec.Id+".json", data, 0644)
	} else {
		panic(err)
	}
	logrus.Infof("end")
}

func item() (ret *TypeItem) {
	ksySpecPath := os.Getenv("KSY_SPEC")
	ksyDataPath := os.Getenv("KSY_DATA")

	logrus.Infof("start")
	var err error
	m, err := NewModel(ksySpecPath, &Options{LazyDecoding: false, PositionFill: false})
	if err != nil {
		panic(err)
	}
	logrus.Infof("model, created")
	ret, err = m.Read(ksyDataPath)

	if err != nil {
		panic(err)
	}
	logrus.Infof("version %v", ret.GetAttrValue(0))
	return
}
