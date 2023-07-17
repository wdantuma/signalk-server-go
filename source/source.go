package source

import "go.einride.tech/can"

type CanSource interface {
	Source() chan can.Frame
	Label() string
}
