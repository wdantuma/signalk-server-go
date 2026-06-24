package actisense

import (
	"bufio"
	"os"
	"path"

	"github.com/wdantuma/signalk-server-go/converter/nmea2000"
	"github.com/wdantuma/signalk-server-go/signalk"
	"github.com/wdantuma/signalk-server-go/source/base"
	"go.einride.tech/can"
)

type canDumpSource struct {
	rawSource chan can.Frame
	source    <-chan signalk.DeltaJson
	file      *os.File
	running   bool
	label     string
}

func (cd *canDumpSource) Source() <-chan signalk.DeltaJson {
	return cd.source
}

func (cd *canDumpSource) Label() string {
	return cd.label
}

func (cd *canDumpSource) Start() {
	cd.running = true
	go func() {
		for {
			fileScanner := bufio.NewScanner(cd.file)
			fileScanner.Split(bufio.ScanLines)
			for fileScanner.Scan() {
				frames, err := Parse(fileScanner.Text())
				if err != nil {
					return
				}
				for _, f := range frames {
					cd.rawSource <- f
					//time.Sleep(10 * time.Millisecond)
				}
			}
			cd.file.Seek(0, 0)
		}
		//canSource.file.Close()
		//close(canSource.Source)
	}()

}

func (cd *canDumpSource) Stop() {
	cd.running = false
}

func NewActiSenseSource(file string, converter nmea2000.Nmea2000ToSignalk) (base.DeltaSource, error) {
	canSource := canDumpSource{}
	var err error
	canSource.file, err = os.Open(file)
	if err != nil {
		return nil, err
	}
	canSource.label = path.Base(file)
	canSource.rawSource = make(chan can.Frame)
	canSource.source = converter.Convert(canSource.label, canSource.rawSource)
	return &canSource, nil
}
