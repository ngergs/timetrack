package states

type State int

const (
	Open State = iota
	Closed
)

func (state State) String() string {
	switch state {
	case Open:
		return "Open"
	case Closed:
		return "Closed"
	default:
		return "Unknown"
	}
}
