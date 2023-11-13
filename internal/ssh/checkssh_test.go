package ssh_test

import (
	. "github.com/monopole/myrepos/internal/ssh"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheckForSshAgent(t *testing.T) {
	assert.NoError(t, ErrIfNoSshAgent())
}

func TestCheckForSshIdentities(t *testing.T) {
	assert.NoError(t, ErrIfNoSshKeys())
}
