package ssh

import (
	"fmt"
	"strings"
	"time"

	"github.com/monopole/myrepos/internal/runner"
)

const (
	NoSshAgentErr = `start an ssh-agent using: eval $(ssh-agent)`
	NoSshKeysErr  = `add keys to your ssh-agent using: ssh-add`
)

func ErrIfNoSshAgent() error {
	r, err := runner.NewRunner("ps", 2*time.Second, nil)
	if err != nil {
		return err
	}
	if err = r.Run("-ef"); err != nil {
		return err
	}
	if !strings.Contains(r.GetOutput(), "ssh-agent") {
		return fmt.Errorf(NoSshAgentErr)
	}
	return nil
}

func ErrIfNoSshKeys() error {
	r, err := runner.NewRunner("ssh-add", 2*time.Second, nil)
	if err != nil {
		return err
	}
	err = r.Run("-l")
	if strings.Contains(r.GetOutput(), "The agent has no identities.") {
		return fmt.Errorf(NoSshKeysErr)
	}
	return err
}
