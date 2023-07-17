package candumpsource

import (
	"bufio"
	"os"
	"path"
	"time"

	"github.com/wdantuma/signalk-server-go/source"
)

type canDumpSource struct {
	source chan source.ExtendedFrame
	file   *os.File
	label  string
}

func (cd *canDumpSource) Source() chan source.ExtendedFrame {
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
	canSource.source = make(chan source.ExtendedFrame)
	go func() {
		for {
			fileScanner := bufio.NewScanner(canSource.file)
			fileScanner.Split(bufio.ScanLines)
			for fileScanner.Scan() {
				f, err := source.Parse(fileScanner.Text())
				if err != nil {
					return
				}
				canSource.source <- *source.NewExtendedFrame(&f)
				time.Sleep(10 * time.Millisecond)
			}
			canSource.file.Seek(0, 0)
		}
		//canSource.file.Close()
		//close(canSource.Source)
	}()

	return &canSource, nil
}
