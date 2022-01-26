package main

import (
	"flag"
	"github.com/go-ee/kaitai"
)

func main() {
	flag.Parse()
	ksyPath := flag.Args()[0]

	model, err := kaitai.NewModel(ksyPath)
	if err != nil {
		panic(err)
	}

	binaryFilePath := flag.Args()[1]
	item, err := model.Read(binaryFilePath)
	if err != nil {
		panic(err)
	}

	println(item)

}
