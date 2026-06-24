package can

import (
	"log"
	"net"

	"github.com/wdantuma/signalk-server-go/converter/nmea2000"
	"github.com/wdantuma/signalk-server-go/signalk"
	"github.com/wdantuma/signalk-server-go/source/base"
	"go.einride.tech/can"
	"go.einride.tech/can/pkg/candevice"
	"go.einride.tech/can/pkg/socketcan"
)

type canSource struct {
	rawSource  chan can.Frame
	source     <-chan signalk.DeltaJson
	device     *candevice.Device
	label      string
	running    bool
	connection net.Conn
}

func (cd *canSource) Source() <-chan signalk.DeltaJson {
	return cd.source
}

func (cd *canSource) Label() string {
	return cd.label
}

func (cd *canSource) Start() {
	cd.running = true
	go func() {
		recv := socketcan.NewReceiver(cd.connection)
		for recv.Receive() && cd.running {
			frame := recv.Frame()
			cd.rawSource <- frame
		}
	}()
}

func (cd *canSource) Stop() {
	cd.running = false
}

func (d *canSource) setupDevice(canDevice string) error {
	device, err := candevice.New(canDevice)
	if err != nil {
		log.Fatal(err)
	}

	//         c := make(chan os.Signal, 1)
	//         signal.Notify(c, os.Interrupt)
	//         go func(){
	//           for _ = range c {
	//            d.SetDown()
	//             // sig is a ^C, handle it
	//           }
	//         }()

	err = device.SetBitrate(250000)
	if err != nil {
		log.Fatal(err)
	}
	err = device.SetUp()
	if err != nil {
		log.Fatal(err)
	}
	d.device = device
	return nil
}

func NewCanSource(canDevice string, converter nmea2000.Nmea2000ToSignalk) (base.DeltaSource, error) {

	canSource := canSource{}
	err := canSource.setupDevice(canDevice)
	if err != nil {
		return nil, err
	}
	canSource.label = canDevice
	canSource.rawSource = make(chan can.Frame)
	canSource.source = converter.Convert(canSource.label, canSource.rawSource)
	conn, err := socketcan.Dial("can", "can0")
	if err != nil {
		log.Fatal(err)
	} else {
		canSource.connection = conn
	}

	return &canSource, nil
}
