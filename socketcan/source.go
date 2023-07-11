package socketcan

import (
	"bufio"
	"os"
	"path"
	"time"
)

type CanSource struct {
	Source chan ExtendedFrame
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
	canSource.Source = make(chan ExtendedFrame)
	go func() {
		for {
			fileScanner := bufio.NewScanner(canSource.file)
			fileScanner.Split(bufio.ScanLines)
			for fileScanner.Scan() {
				f, err := Parse(fileScanner.Text())
				if err != nil {
					return
				}
				canSource.Source <- *NewExtendedFrame(&f)
				time.Sleep(10 * time.Millisecond)
			}
			canSource.file.Seek(0, 0)
		}
		//canSource.file.Close()
		//close(canSource.Source)
	}()

	return &canSource, nil
}
