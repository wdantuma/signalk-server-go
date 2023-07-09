package canboat

import (
	"embed"
	"encoding/xml"
	"io/ioutil"
	"log"
)

type Canboat struct {
	pgnDefinitions  *PGNDefinitions
	pgnIndex        map[int]int
	lookupEnumIndex map[string]int
}

//go:embed canboat.xml
var canboatxml embed.FS

func NewCanboat() (*Canboat, error) {
	c := Canboat{pgnDefinitions: &PGNDefinitions{}, pgnIndex: make(map[int]int), lookupEnumIndex: make(map[string]int)}

	xmlFile, err := canboatxml.Open("canboat.xml")
	if err != nil {
		log.Fatal(err)
	}
	byteValue, _ := ioutil.ReadAll(xmlFile)

	// we initialize our Users array

	err = xml.Unmarshal(byteValue, c.pgnDefinitions)
	if err != nil {
		return nil, err
	}
	// create index
	for i, pgnInfo := range c.pgnDefinitions.PGNs.PGNInfo {
		c.pgnIndex[pgnInfo.PGN] = i
	}
	for i, lookupEnum := range c.pgnDefinitions.LookupEnumerations.LookupEnumeration {
		c.lookupEnumIndex[lookupEnum.Name] = i
	}
	defer xmlFile.Close()

	return &c, nil
}

func (c *Canboat) GetPGNInfo(pgn uint) (*PGNInfo, bool) {
	pgnIndex, ok := c.pgnIndex[int(pgn)]
	if !ok {
		return nil, false
	}
	pgnInfo := c.pgnDefinitions.PGNs.PGNInfo[pgnIndex]
	return &pgnInfo, true
}

func (c *Canboat) GetLookupEnumeration(name string) (*LookupEnumeration, bool) {
	lookupEnumIndex, ok := c.lookupEnumIndex[name]
	if ok {
		return &c.pgnDefinitions.LookupEnumerations.LookupEnumeration[lookupEnumIndex], true
	}
	return nil, false
}
