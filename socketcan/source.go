package socketcan

import (
	"bufio"
	"os"
	"path"
	"time"

	"go.einride.tech/can"
)

type CanSource struct {
	Source chan can.Frame
	file   *os.File
	Label  string
}

func NewCanDumpSource(file string) (*CanSource, error) {
	canSource := CanSource{}
	var err error
	canSource.file, err = os.Open(file)
	if err != nil {
		return nil, err
	}
	canSource.Label = path.Base(file)
	canSource.Source = make(chan can.Frame)
	go func() {
		fileScanner := bufio.NewScanner(canSource.file)
		fileScanner.Split(bufio.ScanLines)
		for fileScanner.Scan() {
			f, err := Parse(fileScanner.Text())
			if err != nil {
				return
			}
			canSource.Source <- f
			time.Sleep(10 * time.Millisecond)
		}
		canSource.file.Close()
		close(canSource.Source)
	}()

	return &canSource, nil
}
