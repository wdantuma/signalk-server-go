package format

import (
	"encoding/json"
	"log"

	"github.com/wdantuma/signalk-server-go/signalk"
)

func Json(input <-chan signalk.DeltaJson, output chan []byte) chan []byte {
	if output == nil {
		output = make(chan []byte)
	}

	go func() {
		for delta := range input {
			deltaBytes, err := json.Marshal(delta)
			if err == nil {
				output <- deltaBytes
			} else {
				log.Fatal(err)
			}
		}
		close(output)
	}()
	return output
}
