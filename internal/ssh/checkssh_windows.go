//go:build windows

package ssh

const (
	noSshAgentErrWindows = `start an ssh-agent using: eval $(ssh-agent)`
)

func errIfNoSshAgent() error {
	// TODO: implement
	// Maybe use https://github.com/mitchellh/go-ps
	// User must have ssh-add installed -
	// see https://interworks.com/blog/2021/09/08/how-to-enable-ssh-commands-in-windows/
	// On windows, start the agent with start-ssh-agent.cmd ?
	return fmt.Errorf(noSshAgentErrWindows)
}

func errIfNoSshKeys() error {
	// TODO: implement - maybe same as linux
	return fmt.Errorf(NoSshKeysErr)
}
