package pkg

import (
	"fmt"
	"path"
	"strings"

	"github.com/monopole/myrepos/internal/file"
	"github.com/monopole/myrepos/internal/runner"
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

type ValidatedRepo struct {
	// domain is the git server domain.
	domain ServerDomain

	// Details about the git server (e.g. port number).
	serverSpec *ServerSpec

	// Organization maintaining the repo.
	origin OrgName

	// The name of the upstream organization from
	// which the repo was forked.  If this is empty,
	// or matches the origin, then the repo isn't a fork.
	upstream OrgName

	// The home of all repos on local disk.
	rootDir file.Path

	// The parent directory to the repository directory.
	dirName file.Path

	// Name of repository (also the name of the repository directory).
	name RepoName

	// Output of most recent git command.
	gr *runner.Runner
}

func (vr *ValidatedRepo) Print() {
	fmt.Printf("        domain: %s\n", vr.domain)
	fmt.Printf("    serverSpec: %s\n", vr.serverSpec)
	fmt.Printf("        origin: %s\n", vr.origin)
	fmt.Printf("      upstream: %s\n", vr.upstream)
	fmt.Printf("     parentDir: %s\n", vr.ParentDir())
	fmt.Printf("      repoName: %s\n", vr.name)
}

func (vr *ValidatedRepo) UrlOrigin() string {
	return vr.urlSpec(vr.origin)
}

func (vr *ValidatedRepo) UrlUpstream() string {
	return vr.urlSpec(vr.upstream)
}

func (vr *ValidatedRepo) urlSpec(o OrgName) string {
	p := path.Join(string(o), string(vr.name)) + ".git"
	if vr.serverSpec.scheme == SchemeHttps {
		// https://github.com/monopole/myrepos.git
		return "https://" + vr.domain.WithPort(vr.serverSpec.port) + "/" + p
	}
	// git@github.com:monopole/myrepos.git
	return "git@" + string(vr.domain) + ":" + p
}

func (vr *ValidatedRepo) Title() string {
	return fmt.Sprintf("%30s%22s%20s%30s : ", vr.rootDir, vr.domain, vr.dirName, vr.name)
}

func (vr *ValidatedRepo) ParentDir() file.Path {
	return vr.rootDir.Append(file.Path(vr.domain)).Append(vr.dirName)
}

func (vr *ValidatedRepo) FullDir() file.Path {
	return vr.ParentDir().Append(file.Path(vr.name))
}

func (vr *ValidatedRepo) IsAFork() bool {
	return vr.origin != vr.upstream
}

func (vr *ValidatedRepo) LastLog() (string, error) {
	if vr.gr == nil {
		return "", fmt.Errorf("no git process")
	}
	if err := vr.gr.Run(
		cmdLog,
		`--pretty=format:"%<(20)%ad%>(24)%an : %s"`,
		`--date=human`,
		`--abbrev=8`,
		`--max-count=1`,
	); err != nil {

		return "", err
	}
	return firstLine(vr.gr.GetOutput()), nil
}

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
func (vr *ValidatedRepo) Clone() (Outcome, error) {
	var err error
	if err = vr.ParentDir().MkDir(); err != nil {
		return Oops, fmt.Errorf("unable to make cloning location %w", err)
	}
	vr.gr, err = runner.NewRunner(gitProgram, vr.serverSpec.duration, commonErrs)
	if err != nil {
		return Oops, err
	}
	vr.gr.SetPwd(vr.ParentDir())
	if err = vr.gr.Run(cmdClone, vr.UrlOrigin()); err != nil {
		return Oops, err
	}
	// Dive in and check it.
	vr.gr.SetPwd(vr.FullDir())
	if vr.IsAFork() {
		fullSpec := vr.UrlUpstream()
		if err = vr.gr.Run(cmdRemote, "add", remoteUpstream, fullSpec); err != nil {
			return Oops, err
		}
		if err = vr.gr.Run(
			cmdRemote, "set-url", "--push", remoteUpstream,
			"disableFootGun_"+fullSpec); err != nil {
			return Oops, err
		}
	}
	return ClonedAt, vr.gr.Run(cmdRemote, "-v")
}

// Rebase switches to the main (or master) branch and rebases.
func (vr *ValidatedRepo) Rebase() (Outcome, error) {
	var (
		err        error
		mainBranch string
	)
	vr.gr, err = runner.NewRunner(gitProgram, vr.serverSpec.duration, commonErrs)
	if err != nil {
		return Oops, err
	}
	vr.gr.SetPwd(vr.FullDir())
	mainBranch, err = vr.determineMainBranch()
	if err != nil {
		return Oops, err
	}
	if err = vr.gr.Run(cmdCheckout, mainBranch); err != nil {
		return Oops, err
	}
	if vr.IsAFork() {
		if err = vr.gr.Run(cmdFetch, remoteUpstream); err != nil {
			return Oops, err
		}
		if err = vr.gr.Run(cmdDiff, path.Join(remoteUpstream, mainBranch)); err != nil {
			return Oops, err
		}
		if len(vr.gr.GetOutput()) == 0 {
			return NoUpdate, nil
		}
		if err = vr.gr.Run(cmdRebase, path.Join(remoteUpstream, mainBranch)); err != nil {
			return Oops, err
		}
		if err = vr.gr.Run(cmdPush, "-f", remoteOrigin, mainBranch); err != nil {
			return Oops, err
		}
		return RebasedTo, nil
	}
	if err = vr.gr.Run(cmdFetch, remoteOrigin); err != nil {
		return Oops, err
	}
	if err = vr.gr.Run(cmdDiff, path.Join(remoteOrigin, mainBranch)); err != nil {
		return Oops, err
	}
	if len(vr.gr.GetOutput()) == 0 {
		return NoUpdate, nil
	}
	if err = vr.gr.Run(cmdRebase, path.Join(remoteOrigin, mainBranch)); err != nil {
		return Oops, err
	}
	return RebasedTo, nil
}

func (vr *ValidatedRepo) determineMainBranch() (string, error) {
	if err := vr.gr.Run(cmdBranch, "--list", "main"); err != nil {
		return "", err
	}
	if len(vr.gr.GetOutput()) > 0 {
		return "main", nil
	}
	return "master", nil
}
