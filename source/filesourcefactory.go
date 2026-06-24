package source

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"strings"

	"github.com/wdantuma/signalk-server-go/converter/nmea0183"
	"github.com/wdantuma/signalk-server-go/converter/nmea2000"
	"github.com/wdantuma/signalk-server-go/source/base"
	"github.com/wdantuma/signalk-server-go/source/nmea0183/pcap"
	"github.com/wdantuma/signalk-server-go/source/nmea2000/actisense"
	"github.com/wdantuma/signalk-server-go/source/nmea2000/candump"
)

func CreateFileSource(filename string, n2kConverter nmea2000.Nmea2000ToSignalk, nmeaConverter nmea0183.Nmea0183ToSignalk) (base.DeltaSource, error) {

	// very simple factory impl for now

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		return nil, scanner.Err()
	}
	parts := strings.Split(scanner.Text(), "#")
	if len(parts) == 2 {
		return candump.NewCanDumpSource(filename, n2kConverter)
	}
	parts = strings.Split(scanner.Text(), ",")

	if len(parts) >= 14 {
		return actisense.NewActiSenseSource(filename, n2kConverter)
	}
	if len(scanner.Bytes()) > 4 {
		if bytes.Equal(scanner.Bytes()[:4], []byte{0xa1, 0xb2, 0xc3, 0xd4}) ||
			bytes.Equal(scanner.Bytes()[:4], []byte{0xd4, 0xc3, 0xb2, 0xa1}) ||
			bytes.Equal(scanner.Bytes()[:4], []byte{0xa1, 0xb2, 0x3c, 0x4d}) ||
			bytes.Equal(scanner.Bytes()[:4], []byte{0x4d, 0x3c, 0xb2, 0xa1}) {
			return pcap.NewPcapSource(filename, nmeaConverter)
		}
	}

	return nil, errors.New("No matching format found")
}
