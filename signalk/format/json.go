package format

import (
	"encoding/json"

	"github.com/wdantuma/signalk-server-go/signalk"
)

func Json(input <-chan signalk.DeltaJson) <-chan []byte {
	output := make(chan []byte)
	go func() {
		for delta := range input {
			deltaBytes, _ := json.Marshal(delta)
			output <- deltaBytes
		}
		close(output)
	}()
	return output
}