package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/wdantuma/signalk-server-go/s57"
	"github.com/wdantuma/signalk-server-go/s57/dataset"
)

func main() {

	outputPath := flag.String("out", "./static/charts", "Output directory for vector tiles")
	inputPath := flag.String("in", "./charts", "Input path S-57 maps")
	minzoom := flag.Int("minzoom", 14, "Min zoom")
	maxzoom := flag.Int("maxzoom", 14, "Max zoom")
	flag.Parse()

	datasets, err := dataset.GetS57Datasets(*inputPath)
	if err != nil {
		log.Fatal(err)
	}

	tiler := s57.NewS57Tiler(datasets)

	for z := *minzoom; z <= *maxzoom; z++ {
		tiles := tiler.GetTiles(datasets[0], z)

		total := len(tiles)
		n := 0
		for k := range tiles {
			tiler.GenerateTile(*outputPath, datasets[0], tiles[k])
			done := float64(n) / float64(total) * 100
			fmt.Printf("\rZoom: %d, Processed: %.0f %%    ", z, done)
			n++
		}
	}

}
