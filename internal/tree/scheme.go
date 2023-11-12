package tree

type Scheme int

const (
	SchemeUnknown Scheme = iota // unknown
	SchemeSsh                   // ssh
	SchemeHttps                 // https
)

func (s Scheme) String() string {
	return []string{"unknown", "ssh", "https"}[s]
}
