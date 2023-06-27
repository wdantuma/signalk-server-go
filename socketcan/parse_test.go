package socketcan

import (
	"fmt"
	"log"
	"testing"
)

func TestFrameParse(t *testing.T) {
	candumpString := "09FD0269#00A102B54DFAFFFF"
	frame, err := Parse(candumpString)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(frame.String())
	if frame.String() != candumpString {
		t.Error()
	}
}
