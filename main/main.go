package main

import (
	"flag"
	"github.com/go-ee/kaitai"
	"log"
)

func main() {
	flag.Parse()
	path := flag.Args()[0]
	model, err := kaitai.NewModelFromYamlFile(path)
	if err != nil {
		panic(err)
	}
	log.Println(model.Root.ReferencesResolved())
	log.Println(model.Info())
}
