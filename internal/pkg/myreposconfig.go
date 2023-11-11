package pkg

import (
	"fmt"
	"sort"
	"strings"

	"github.com/monopole/myrepos/internal/file"
)

type MyReposConfig struct {
	// Path is the path below which all repos are found.
	// It's specified outside of Layout to avoid extra indents.
	// If the path lacks a leading '/', it will be interpreted
	// as a path relative to value of the HOME environment
	// variable. If HOME is undefined, '.' is used.
	Path file.Path

	// Layout is the directory layout below Path.
	//
	// Just below Path one finds directories named after git
	// server domains. Below these one finds directories that
	// are usually named after "organizations", and below each of
	// these are the repositories maintained by that organization.
	//
	// The exception to this is when one wants a directory
	// name that doesn't match an organization name, and/or when
	// one wants to indicate that the repository was forked
	// from another organization.
	//
	// To do this, add pipe characters to the OrgName field.
	// For example, if the OrgName is specified as
	//
	//   sigs.k8s.io|monopole|kubernetes-sigs:
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
	Layout map[ServerDomain]map[OrgName][]RepoName

	// ServerOpts is a mapping from a git server domain name
	// to optional details about the serverSpec, like the scheme to
	// use when cloning, what timeout to use, what port, etc.
	ServerOpts map[ServerDomain]ServerOpts `yaml:"serverOpts"`
}

type ServerDomain string

func (d ServerDomain) WithPort(p int) string {
	if p > 0 {
		return fmt.Sprintf("%s:%d", d, p)
	}
	return string(d)
}

type OrgName string

type RepoName string

func (mb *MyReposConfig) ToRepos() (result []*ValidatedRepo, err error) {
	rootDir := mb.absRootDir()
	if _, isDir := rootDir.Exists(); !isDir {
		// TODO: make the directory instead of complain that it's missing?
		// Could be a big mistake, as this would then clone all the repos into it.
		return nil, fmt.Errorf("repo root directory %q doesn't exist", rootDir)
	}

	serverSpec := make(map[ServerDomain]*ServerSpec)
	for d, s := range mb.ServerOpts {
		serverSpec[d], err = s.ToServerSpec()
		if err != nil {
			return
		}
	}
	for domain, orgMap := range mb.Layout {
		if _, ok := serverSpec[domain]; !ok {
			serverSpec[domain] = MakeServerSpec()
		}
		for orgName, repoList := range orgMap {
			dirName, origin, upstream := parseOrgName(orgName)
			for _, repoName := range repoList {
				result = append(result, &ValidatedRepo{
					domain:     domain,
					serverSpec: serverSpec[domain],
					origin:     origin,
					upstream:   upstream,
					rootDir:    rootDir,
					dirName:    dirName,
					name:       repoName,
				})
			}
		}
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].domain != result[j].domain {
			return result[i].domain < result[j].domain
		}
		if result[i].dirName != result[j].dirName {
			return result[i].dirName < result[j].dirName
		}
		return result[i].name < result[j].name
	})
	return
}

func (mb *MyReposConfig) absRootDir() file.Path {
	if mb.Path == "" {
		return file.Home()
	}
	if mb.Path.IsAbs() {
		return mb.Path
	}
	return file.Home().Append(mb.Path)
}

func parseOrgName(orgName OrgName) (file.Path, OrgName, OrgName) {
	n := strings.Split(string(orgName), "|")
	if len(n) > 2 {
		return file.Path(n[0]), OrgName(n[1]), OrgName(n[2])
	}
	if len(n) > 1 {
		return file.Path(n[0]), OrgName(n[1]), OrgName(n[1])
	}
	return file.Path(n[0]), OrgName(n[0]), OrgName(n[0])
}
