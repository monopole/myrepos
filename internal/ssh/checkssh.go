package ssh

import (
	"runtime"
)

const (
	NoSshKeysErr = `add keys to your ssh-agent using: ssh-add`
)

func isOpSysWindows() bool {
	return runtime.GOOS == "windows"
}

func ErrIfNoSshAgent() error {
	return errIfNoSshAgent()
}

func ErrIfNoSshKeys() error {
	return errIfNoSshKeys()
}
