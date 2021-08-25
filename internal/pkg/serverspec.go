package pkg

import (
	"fmt"
	"time"
)

type ServerSpec struct {
	port     int
	scheme   Scheme
	duration time.Duration
}

func (s ServerSpec) String() string {
	return fmt.Sprintf("schema=%s port=%d timeout=%s",
		s.scheme.String(), s.port, s.duration.String())
}

// MakeServerSpec returns a ServerSpec with default values.
func MakeServerSpec() *ServerSpec {
	return &ServerSpec{
		port:     0,
		scheme:   SchemeSsh,
		duration: 4 * time.Minute,
	}
}
