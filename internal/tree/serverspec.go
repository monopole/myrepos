package tree

import (
	"fmt"
	"github.com/monopole/myrepos/internal/config"
	"time"
)

type ServerSpec struct {
	port    int
	scheme  pkg.Scheme
	timeout time.Duration
}

func (s *ServerSpec) String() string {
	return fmt.Sprintf("schema=%s port=%d timeout=%s",
		s.scheme.String(), s.port, s.Timeout())
}

func (s *ServerSpec) Timeout() time.Duration {
	return s.timeout
}

func (s *ServerSpec) Scheme() pkg.Scheme {
	return s.scheme
}

func (s *ServerSpec) Port() int {
	return s.port
}

// MakeServerSpec returns a ServerSpec with default values.
func MakeServerSpec() *ServerSpec {
	return &ServerSpec{
		port:    0,
		scheme:  pkg.SchemeSsh,
		timeout: 4 * time.Minute,
	}
}

// FromServerOpts creates a ServerSpec from its serialized form.
func FromServerOpts(s *pkg.ServerOpts) (result *ServerSpec, err error) {
	result = MakeServerSpec()
	if s.Timeout != "" {
		result.timeout, err = time.ParseDuration(s.Timeout)
		if err != nil {
			return nil, fmt.Errorf("bad duration %q in git serverSpec; %w", s.Timeout, err)
		}
	}
	if s.Scheme != "" {
		switch s.Scheme {
		case pkg.SchemeHttps.String():
			result.scheme = pkg.SchemeHttps
		case pkg.SchemeSsh.String():
			result.scheme = pkg.SchemeSsh
		default:
			return nil, fmt.Errorf("unknown scheme %q in serverSpec opts", s.Scheme)
		}
	}
	return result, nil
}
