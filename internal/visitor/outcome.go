package visitor

import "github.com/TwiN/go-color"

type Outcome int

const (
	Oops Outcome = iota
	RebasedTo
	ClonedAt
	NoUpdate
)

func (o Outcome) String() string {
	return []string{
		color.Red + "error" + color.Reset,
		color.Blue + "rebased to" + color.Reset,
		color.Green + "cloned to latest at" + color.Reset,
		color.Green + "no change since" + color.Reset,
	}[o]
}
