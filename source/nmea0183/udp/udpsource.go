package udp

import (
	"net"

	"github.com/adrianmo/go-nmea"
	"github.com/wdantuma/signalk-server-go/converter/nmea0183"
	"github.com/wdantuma/signalk-server-go/signalk"
	"github.com/wdantuma/signalk-server-go/source/base"
)

const (
	MAX_DGRAM_SIZE = 64 * 1024
)

type udpSource struct {
	rawSource     chan nmea.Sentence
	source        <-chan signalk.DeltaJson
	label         string
	iface         *net.Interface
	listenAddress string
	running       bool
	address       *net.UDPAddr
}

func NewUdpSource(address string, converter nmea0183.Nmea0183ToSignalk) (base.DeltaSource, error) {
	udpSource := udpSource{}
	udpSource.label = address
	udpSource.rawSource = make(chan nmea.Sentence)
	udpSource.source = converter.Convert(udpSource.label, udpSource.rawSource)
	udpSource.running = true
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil, err
	}
	udpSource.address = addr
	return &udpSource, nil
}

func (cd *udpSource) Source() <-chan signalk.DeltaJson {
	return cd.source
}

func (fs *udpSource) Start() {
	fs.running = true
	go func() {

		l, err := net.ListenUDP("udp", fs.address)
		if err != nil {
			panic(err)
		}
		l.SetReadBuffer(MAX_DGRAM_SIZE)
		defer l.Close()
		buf := make([]byte, MAX_DGRAM_SIZE)
		for fs.running {
			n, _, _ := l.ReadFromUDP(buf)
			if fs.running {
				sentence, err := nmea.Parse(string(buf[:n]))
				if err == nil {
					fs.rawSource <- sentence
				}
			}
		}
	}()
}

func (fs *udpSource) Stop() {
	fs.running = false
	close(fs.rawSource)
}

func (p *udpSource) Label() string {
	return "UDP source"
}
