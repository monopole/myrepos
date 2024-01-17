//go:build linux

package ssh

import (
	"fmt"
	"strings"
	"time"

	"github.com/monopole/myrepos/internal/runner"
)

const (
	noSshAgentErrLinux = `start an ssh-agent using: eval $(ssh-agent)`
)

func errIfNoSshAgent() error {
	r, err := runner.NewRunner("ps", 2*time.Second, nil)
	if err != nil {
		return err
	}
	if err = r.Run("-ef"); err != nil {
		return err
	}
	if !strings.Contains(r.GetOutput(), "ssh-agent") {
		return fmt.Errorf(noSshAgentErrLinux)
	}
	return nil
}

func errIfNoSshKeys() error {
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
