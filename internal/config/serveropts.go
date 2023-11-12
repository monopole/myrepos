package pkg

// ServerOpts is a serialized form of ServerSpec.
type ServerOpts struct {
	Port    int
	Scheme  string
	Timeout string
}
