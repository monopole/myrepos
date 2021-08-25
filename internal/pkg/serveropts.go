package pkg

import (
	"fmt"
	"time"
)

// ServerOpts is a serialized form of ServerSpec.
type ServerOpts struct {
	Port    int
	Scheme  string
	Timeout string
}

// ToServerSpec creates a ServerSpec from its serialized form.
func (s *ServerOpts) ToServerSpec() (result *ServerSpec, err error) {
	result = MakeServerSpec()
	if s.Timeout != "" {
		result.duration, err = time.ParseDuration(s.Timeout)
		if err != nil {
			return nil, fmt.Errorf("bad duration %q in git serverSpec; %w", s.Timeout, err)
		}
	}
	if s.Scheme != "" {
		switch s.Scheme {
		case SchemeHttps.String():
			result.scheme = SchemeHttps
		case SchemeSsh.String():
			result.scheme = SchemeSsh
		default:
			return nil, fmt.Errorf("unknown scheme %q in serverSpec opts", s.Scheme)
		}
	}
	return result, nil
}
