package pcap

import (
	"fmt"
	"log"
	"path"
	"time"

	"github.com/adrianmo/go-nmea"
	"github.com/google/gopacket"
	"github.com/google/gopacket/ip4defrag"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/wdantuma/signalk-server-go/converter/nmea0183"
	"github.com/wdantuma/signalk-server-go/signalk"
	"github.com/wdantuma/signalk-server-go/source/base"
)

type pcapSource struct {
	rawSource chan nmea.Sentence
	source    <-chan signalk.DeltaJson
	handle    *pcap.Handle
	label     string
	running   bool
	loop      bool
}

func NewPcapSource(file string, converter nmea0183.Nmea0183ToSignalk) (base.DeltaSource, error) {
	pcapSource := pcapSource{}
	if handle, err := pcap.OpenOffline(file); err != nil {
		return nil, err
	} else {
		pcapSource.handle = handle
		pcapSource.label = path.Base(file)
		pcapSource.rawSource = make(chan nmea.Sentence)
		pcapSource.source = converter.Convert(pcapSource.label, pcapSource.rawSource)
		return &pcapSource, nil
	}
}

func (cd *pcapSource) Source() <-chan signalk.DeltaJson {
	return cd.source
}

func (cd *pcapSource) Label() string {
	return cd.label
}

func (p *pcapSource) Start() {
	p.running = true
	go func() {
		if p.loop {
			for p.running && p.loop {
				p.processFile()
			}
		} else {
			p.processFile()
		}
		p.Stop()
	}()
}

func (p *pcapSource) Stop() {
	p.running = false
}

func (p *pcapSource) processFile() {
	n := 0
	packetSource := gopacket.NewPacketSource(p.handle, p.handle.LinkType())
	defragger := ip4defrag.NewIPv4Defragmenter()
	var prevTimestamp time.Time = time.Time{}
	var startTimestamp = time.Now()
	for packet := range packetSource.Packets() {
		n++
		ip4Layer := packet.Layer(layers.LayerTypeIPv4)
		if ip4Layer == nil {
			continue
		}
		ip4 := ip4Layer.(*layers.IPv4)
		l := ip4.Length

		newip4, err := defragger.DefragIPv4(ip4)
		if err != nil {
			log.Fatalln("Error while de-fragmenting", err)
		} else if newip4 == nil {
			continue // packet fragment, we don't have whole packet yet.
		}
		if newip4.Length != l {
			pb, ok := packet.(gopacket.PacketBuilder)
			if !ok {
				panic("Not a PacketBuilder")
			}
			nextDecoder := newip4.NextLayerType()
			nextDecoder.Decode(newip4.Payload, pb)
		}
		if !prevTimestamp.Equal(time.Time{}) {
			duration := packet.Metadata().Timestamp.Sub(prevTimestamp)
			time.Sleep(duration)
		}
		prevTimestamp = packet.Metadata().Timestamp
		if !p.running {
			break
		}
		p.handlePacket(packet)
	}
	totalDuration := time.Since(startTimestamp)
	fmt.Printf("Processed :%d packets in %s\n", n, totalDuration)
}

func (p *pcapSource) handlePacket(packet gopacket.Packet) {

	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	udpLayer := packet.Layer(layers.LayerTypeUDP)

	if ipLayer != nil && udpLayer != nil {
		dstPort := udpLayer.(*layers.UDP).DstPort

		if dstPort == 10110 {
			sentence, err := nmea.Parse(string(udpLayer.LayerPayload()))
			if err == nil {
				p.rawSource <- sentence
			}
		}
	}
}
