package main

import (
	"fmt"
	"log"

	"github.com/wdantuma/signalk-server-go/s57"
	"github.com/wdantuma/signalk-server-go/s57/dataset"
)

func main() {
	datasets, err := dataset.GetS57Datasets("../../charts")
	if err != nil {
		log.Fatal(err)
	}

	tiler := s57.NewS57Tiler(datasets)
	tiles := tiler.GetTiles(datasets[0], 14)

	total := len(tiles)
	n := 0
	for k := range tiles {
		tiler.GenerateTile("../../static/charts", datasets[0], tiles[k])
		done := float64(n) / float64(total) * 100
		fmt.Printf("\rProcessed: %.0f %%    ", done)
		n++
	}
}
