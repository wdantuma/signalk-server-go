package socketcan

import (
	"log"

	"go.einride.tech/can"
	"go.einride.tech/can/pkg/candevice"
	"go.einride.tech/can/pkg/socketcan"
)

type canSource struct {
	source chan can.Frame
	device *candevice.Device
	label  string
}

func (cd *canSource) Source() chan can.Frame {
	return cd.source
}

func (cd *canSource) Label() string {
	return cd.label
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

func NewCanSource(canDevice string) (*canSource, error) {

	canSource := canSource{}
	err := canSource.setupDevice(canDevice)
	if err != nil {
		return nil, err
	}
	canSource.label = canDevice
	canSource.source = make(chan can.Frame)
	conn, err := socketcan.Dial("can", "can0")
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		recv := socketcan.NewReceiver(conn)
		for recv.Receive() {
			frame := recv.Frame()
			canSource.source <- frame
		}
	}()

	return &canSource, nil
}
