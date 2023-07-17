package source

type CanSource interface {
	Source() chan ExtendedFrame
	Label() string
}
