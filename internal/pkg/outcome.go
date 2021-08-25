package pkg

type Outcome int

const (
	Oops Outcome = iota
	RebasedTo
	ClonedAt
	NoUpdate
)

func (o Outcome) String() string {
	return []string{"error", "rebased to", "cloned to latest at", "no change since"}[o]
}
