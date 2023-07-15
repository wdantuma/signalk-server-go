package converter

import (
	"fmt"
	"log"
	"testing"

	"github.com/wdantuma/signalk-server-go/converter"
	"github.com/wdantuma/signalk-server-go/signalk/filter"
	"github.com/wdantuma/signalk-server-go/signalk/format"
	"github.com/wdantuma/signalk-server-go/signalkserver"
	"github.com/wdantuma/signalk-server-go/socketcan"
)

func TestParse(t *testing.T) {
	source, err := socketcan.NewCanDumpSource("../../data/n2kdump.txt")
	if err != nil {
		log.Fatal(err)
	}
	state := signalkserver.NewSignalkServer()
	converter, err := converter.NewCanToSignalk(state)
	if err != nil {
		log.Fatal(err)
	}

	sk := converter.Convert(state, source)
	filter := filter.NewFilter(state.GetSelf())
	f := filter.Filter(sk)
	json := format.Json(f, nil)

	for bytes := range json {
		fmt.Println(string(bytes))
	}
}
