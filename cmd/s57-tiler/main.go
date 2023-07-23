package main

import (
	"fmt"
	"os"

	"github.com/lukeroth/gdal"
	"github.com/tburke/iso8211"
)

func main() {
	driver, err := gdal.GetDriverByName("S57")
	if err != nil {
		fmt.Println(err.Error())
	}
	driver.Register()

	datasource := gdal.OpenDataSource("../../charts/ENC_ROOT/1R/7/1R7EMS01/1R7EMS01.000", 0)

	for i := 0; i < datasource.LayerCount(); i++ {
		l := datasource.LayerByIndex(i)
		count, ok := l.FeatureCount(false)
		if ok {
			ext, err := l.Extent(true)
			if err == nil {

				fmt.Printf("%s (%d) %f %f %f %f\n", l.Name(), count, ext.MinX(), ext.MinY(), ext.MaxX(), ext.MaxY())

			}
		}

	}

	f, err := os.Open("../../charts/ENC_ROOT/catalog.031")
	if err != nil {
		fmt.Println("No file. ", err)
	}

	var l iso8211.LeadRecord
	l.Read(f)
	var d iso8211.DataRecord
	d.Lead = &l
	for d.Read(f) == nil {
		fmt.Printf("-----\n")
		for i := 0; i < len(d.Fields); i++ {
			for n := 0; n < len(d.Fields[i].SubFields); n++ {
				fmt.Printf("Field (%d) subfield(%d): %s\n", i, n, d.Fields[i].SubFields[n])
			}

		}
	}

	//defer datasource.
}
