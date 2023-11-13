package config

import (
	"fmt"
	"strings"

	"github.com/monopole/myrepos/internal/file"
)

type Config struct {
	// Path is the root path to all local storage repos.
	// It's specified outside of Layout to avoid extra indents.
	// If the path lacks a leading '/', it will be interpreted
	// as a path relative to value of the HOME environment
	// variable. If HOME is undefined, '.' is used.
	Path file.Path

	// Layout is the directory layout below Path.
	//
	// Just below Path one finds directories named after git
	// server domains. Below these one finds directories that
	// are usually named after "organizations", e.g google,
	// kubernetes, internal departments, etc. Below each of
	// these are the repositories maintained by that organization.
	//
	// Sometimes one wants a directory name that doesn't match
	// an organization name, and/or one wants to indicate that
	// the repository was forked from another organization.
	// See the OrgName field description for notes on doing this.
	Layout map[ServerDomain]map[OrgName][]RepoName

	// ServerOpts is a mapping from a git server domain name
	// to optional details about the git server, like the scheme to
	// use when cloning, what timeout to use, what port, etc.
	ServerOpts map[ServerDomain]ServerOpts `yaml:"serverOpts"`
}

// ServerDomain is the domain of the git server (e.g. github.com).
type ServerDomain string

// RepoName is the name of a repository, e.g. kubectl
type RepoName string

// OrgName names the git "organization".
//
// Sometimes one wants a directory name that doesn't match
// an organization name, and/or one wants to indicate that
// the repository was forked from another organization.
//
// To do this, add pipe characters ('|') to the OrgName field.
// For example, if the OrgName is specified as
//
//	sigs.k8s.io|monopole|kubernetes-sigs:
//
// then the repo will be cloned into a directory named
// 'sigs.k8s.io', its 'origin' remote will be set to
// 'monopole', and its 'upstream' remote will be set
// to 'kubernetes-sigs'.
//
// This flexibility is required for working with repositories
// that contain Go modules, because the leading part of the
// Go module name can differ from the name of the GitHub
// organization that maintains the module. The organization
// that maintains a Go module might change multiple times for
// whatever reason, but the Go module path must be sticky
// over time to avoid breaking Go import statements.
// Redirection services mitigate organization name changes.
type OrgName string

// ServerOpts provides details about using the git server
type ServerOpts struct {
	// What port (if not the default port) to use in the git clone url.
	Port int
	// Specify "https" or "ssh".
	Scheme string
	// How long to wait for a git operation?  Use time.Duration
	// format, e.g. '80s' or '10m'.
	Timeout string
}

func (on OrgName) Parse() (file.Path, OrgName, OrgName) {
	n := strings.Split(string(on), "|")
	if len(n) > 2 {
		return file.Path(n[0]), OrgName(n[1]), OrgName(n[2])
	}
	if len(n) > 1 {
		return file.Path(n[0]), OrgName(n[1]), OrgName(n[1])
	}
	return file.Path(n[0]), OrgName(n[0]), OrgName(n[0])
}

func (d ServerDomain) WithPort(p int) string {
	if p > 0 {
		return fmt.Sprintf("%s:%d", d, p)
	}
	return string(d)
}
