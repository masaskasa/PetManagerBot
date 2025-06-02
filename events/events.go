package events

type IFetcher interface {
	Fetch(limit int) ([]Event, error)
}

type IProcessor interface {
	Process(event Event) error
}

type Type int

const (
	Unknown Type = iota + 1
	Message
)

type Event struct {
	Type Type
	Text string
	Meta interface{}
}
