package canboat

import (
	"embed"
	"encoding/xml"
	"io/ioutil"
	"log"
)

var Version = "5.0.1"

type Canboat struct {
	pgnDefinitions           *PGNDefinitions
	pgnIndex                 map[int]int
	lookupEnumIndex          map[string]int
	lookupFieldTypeEnumIndex map[string]int
}

//go:embed canboat.xml
var canboatxml embed.FS

func NewCanboat() (*Canboat, error) {
	c := Canboat{pgnDefinitions: &PGNDefinitions{}, pgnIndex: make(map[int]int), lookupEnumIndex: make(map[string]int), lookupFieldTypeEnumIndex: make(map[string]int)}

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
	for i, lookupFieldTypeEnum := range c.pgnDefinitions.LookupFieldTypeEnumerations.LookupFieldTypeEnumeration {
		c.lookupFieldTypeEnumIndex[lookupFieldTypeEnum.Name] = i
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

func (c *Canboat) GetLookupEnumeration(name string, value float64) (string, bool) {
	lookupEnumIndex, ok := c.lookupEnumIndex[name]
	if ok {
		lookupEnumeration := c.pgnDefinitions.LookupEnumerations.LookupEnumeration[lookupEnumIndex]
		for _, v := range lookupEnumeration.EnumPair {
			if v.ValueAttr == uint(value) {
				return v.Name, true
			}
		}
	}

	return "", false
}

func (c *Canboat) GetLookupFieldTypeEnumeration(name string, value float64) (*EnumFieldType, bool) {
	lookupFieldTypeEnumIndex, ok := c.lookupFieldTypeEnumIndex[name]
	if ok {
		lookupFieldTypeEnumeration := c.pgnDefinitions.LookupFieldTypeEnumerations.LookupFieldTypeEnumeration[lookupFieldTypeEnumIndex]
		for _, v := range lookupFieldTypeEnumeration.EnumFieldType {
			if int(v.Value) == int(value) {
				return &v, true
			}
		}
	}
	return nil, false
}
