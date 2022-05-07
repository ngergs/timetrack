package modes

type Mode int

const (
	Start Mode = iota
	Stop
	Status
	Sheet
)

var parseMap = map[string]Mode{
	"start":  Start,
	"stop":   Stop,
	"status": Status,
	"sheet":  Sheet,
}

// Parse maps the input string to a Mode enum. ok is false if no corresponding mapping has been found.
func Parse(input string) (value Mode, ok bool) {
	value, ok = parseMap[input]
	return
}
