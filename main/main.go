package main

import (
	"flag"
	"fmt"
	"github.com/go-ee/kaitai"
	"strings"
	"time"
)

const (
	day  = time.Minute * 60 * 24
	year = 365 * day
)

func main() {
	flag.Parse()
	ksyPath := flag.Args()[0]

	model, err := kaitai.NewModel(ksyPath)
	if err != nil {
		panic(err)
	}

	start := time.Now()
	binaryFilePath := flag.Args()[1]
	item, err := model.Read(binaryFilePath)
	if err != nil {
		panic(err)
	}

	m := item.Value().(map[string]*kaitai.Item)
	recordsItem := m["records"]
	records := recordsItem.Value().([]*kaitai.Item)

	println(len(records), duration(time.Now().Sub(start)))

}

func duration(d time.Duration) string {
	if d < day {
		return d.String()
	}

	var b strings.Builder
	if d >= year {
		years := d / year
		fmt.Fprintf(&b, "%dy", years)
		d -= years * year
	}

	days := d / day
	d -= days * day
	fmt.Fprintf(&b, "%dd%s", days, d)

	return b.String()
}
