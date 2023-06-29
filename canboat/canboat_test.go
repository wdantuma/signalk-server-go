package canboat

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	xmlFile, err := os.Open("canboat.xml")
	if err != nil {
		log.Fatal(err)
	}
	var pgnDefinitions PGNDefinitions
	byteValue, _ := ioutil.ReadAll(xmlFile)

	// we initialize our Users array

	err = xml.Unmarshal(byteValue, &pgnDefinitions)
	if err != nil {
		t.Error()
	}
	defer xmlFile.Close()
}
