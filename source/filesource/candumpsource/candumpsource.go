package candumpsource

import (
	"bufio"
	"os"
	"path"
	"time"

	"go.einride.tech/can"
)

type canDumpSource struct {
	source chan can.Frame
	file   *os.File
	label  string
}

func (cd *canDumpSource) Source() chan can.Frame {
	return cd.source
}

func (cd *canDumpSource) Label() string {
	return cd.label
}

func NewCanDumpSource(file string) (*canDumpSource, error) {
	canSource := canDumpSource{}
	var err error
	canSource.file, err = os.Open(file)
	if err != nil {
		return nil, err
	}
	canSource.label = path.Base(file)
	canSource.source = make(chan can.Frame)
	go func() {
		for {
			fileScanner := bufio.NewScanner(canSource.file)
			fileScanner.Split(bufio.ScanLines)
			for fileScanner.Scan() {
				f, err := Parse(fileScanner.Text())
				if err != nil {
					return
				}
				canSource.source <- f
				time.Sleep(10 * time.Millisecond)
			}
			canSource.file.Seek(0, 0)
		}
		//canSource.file.Close()
		//close(canSource.Source)
	}()

	return &canSource, nil
}
