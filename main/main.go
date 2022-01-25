package main

import (
	"flag"
	"fmt"
	"github.com/go-ee/kaitai"
	"log"
)

func main() {
	flag.Parse()
	ksyPath := flag.Args()[0]
	model, err := kaitai.ParseToModelFromYamlFile(ksyPath)
	if err != nil {
		panic(err)
	}
	log.Println(fmt.Sprintf("resolved: %v, %v", model.Root.Id, model.Root.RefsResolved()))

	binaryFilePath := flag.Args()[1]
	item, err := model.Read(binaryFilePath)
	if err != nil {
		panic(err)
	}

	println(item)

}
