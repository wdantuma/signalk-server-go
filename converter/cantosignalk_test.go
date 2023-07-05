package converter

import (
	"fmt"
	"log"
	"testing"

	"github.com/wdantuma/signalk-server-go/signalk/filter"
	"github.com/wdantuma/signalk-server-go/signalk/format"
	"github.com/wdantuma/signalk-server-go/socketcan"
)

func TestParse(t *testing.T) {
	source, err := socketcan.NewCanDumpSource("../data/n2kdump.txt")
	if err != nil {
		log.Fatal(err)
	}
	converter, err := NewCanToSignalk()
	if err != nil {
		log.Fatal(err)
	}

	sk := converter.Convert(source)
	filter := filter.NewFilter()
	f := filter.Filter(sk)
	json := format.Json(f)

	for bytes := range json {
		fmt.Println(string(bytes))
	}
}
