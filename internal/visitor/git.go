package visitor

import (
	"fmt"
	"github.com/monopole/myrepos/internal/runner"
	"github.com/monopole/myrepos/internal/tree"
	"path"
	"strings"
)

const gitProgram = "git"

// commonErrs maps common error substrings to one line summaries.
var commonErrs = map[string]string{
	"Could not resolve hostname": "cannot reach host",
	"Connection refused":         "cannot reach host",
	"You have unstaged changes":  "unstaged changes - commit or stash first",
	"Operation timed out":        "timed out - is repo accessible?",
}

const (
	remoteUpstream = "upstream"
	remoteOrigin   = "origin"
	cmdDiff        = "diff"
	cmdFetch       = "fetch"
	cmdRebase      = "rebase"
	cmdClone       = "clone"
	cmdRemote      = "remote"
	cmdLog         = "log"
	cmdBranch      = "branch"
	cmdPush        = "push"
	cmdCheckout    = "checkout"
)

func deQuote(arg string) string {
	return strings.Trim(arg, `'"`)
}

func firstLine(arg string) string {
	result := strings.Split(arg, "\n")
	if len(result) > 0 {
		return deQuote(result[0])
	}
	return deQuote(arg)
}

// Clone attempts to clone the repo.
func (v *Cloner) Clone(n *tree.RepoNode) (Outcome, error) {
	var err error
	if err = n.AbsPath().MkDir(); err != nil {
		return Oops, fmt.Errorf("unable to make cloning location %w", err)
	}
	v.gr, err = runner.NewRunner(gitProgram, n.ServerSpec().Timeout(), commonErrs)
	if err != nil {
		return Oops, err
	}
	v.gr.SetPwd(n.AbsParent())
	if err = v.gr.Run(cmdClone, n.UrlOrigin()); err != nil {
		return Oops, err
	}
	// Dive in and check it.
	v.gr.SetPwd(n.AbsPath())
	if n.IsAFork() {
		fullSpec := n.UrlUpstream()
		if err = v.gr.Run(cmdRemote, "add", remoteUpstream, fullSpec); err != nil {
			return Oops, err
		}
		if err = v.gr.Run(
			cmdRemote, "set-url", "--push", remoteUpstream,
			"disableFootGun_"+fullSpec); err != nil {
			return Oops, err
		}
	}
	return ClonedAt, v.gr.Run(cmdRemote, "-v")
}

// Rebase switches to the main (or master) branch and rebases.
func (v *Cloner) Rebase(n *tree.RepoNode) (Outcome, error) {
	var (
		err        error
		mainBranch string
	)
	v.gr, err = runner.NewRunner(gitProgram, n.ServerSpec().Timeout(), commonErrs)
	if err != nil {
		return Oops, err
	}
	v.gr.SetPwd(n.AbsPath())
	mainBranch, err = v.determineMainBranch()
	if err != nil {
		return Oops, err
	}
	if err = v.gr.Run(cmdCheckout, mainBranch); err != nil {
		return Oops, err
	}
	if n.IsAFork() {
		if err = v.gr.Run(cmdFetch, remoteUpstream); err != nil {
			return Oops, err
		}
		if err = v.gr.Run(cmdDiff, path.Join(remoteUpstream, mainBranch)); err != nil {
			return Oops, err
		}
		if len(v.gr.GetOutput()) == 0 {
			return NoUpdate, nil
		}
		if err = v.gr.Run(cmdRebase, path.Join(remoteUpstream, mainBranch)); err != nil {
			return Oops, err
		}
		if err = v.gr.Run(cmdPush, "-f", remoteOrigin, mainBranch); err != nil {
			return Oops, err
		}
		return RebasedTo, nil
	}
	if err = v.gr.Run(cmdFetch, remoteOrigin); err != nil {
		return Oops, err
	}
	if err = v.gr.Run(cmdDiff, path.Join(remoteOrigin, mainBranch)); err != nil {
		return Oops, err
	}
	if len(v.gr.GetOutput()) == 0 {
		return NoUpdate, nil
	}
	if err = v.gr.Run(cmdRebase, path.Join(remoteOrigin, mainBranch)); err != nil {
		return Oops, err
	}
	return RebasedTo, nil
}

func (v *Cloner) determineMainBranch() (string, error) {
	if err := v.gr.Run(cmdBranch, "--list", "main"); err != nil {
		return "", err
	}
	if len(v.gr.GetOutput()) > 0 {
		return "main", nil
	}
	return "master", nil
}

func (v *Cloner) LastLog() (string, error) {
	if v.gr == nil {
		return "", fmt.Errorf("no git process")
	}
	if err := v.gr.Run(
		cmdLog,
		`--pretty=format:"%<(26)%ad%>(30)%an : %s"`,
		`--date=human`,
		`--abbrev=8`,
		`--max-count=1`,
	); err != nil {

		return "", err
	}
	return firstLine(v.gr.GetOutput()), nil
}
