package tree

import (
	"github.com/monopole/myrepos/internal/config"
	"github.com/monopole/myrepos/internal/file"
	"path"
)

type RepoNode struct {
	parent *OrgNode
	Name   string
}

func (n *RepoNode) Accept(v Visitor) {
	v.VisitRepoNode(n)
}

func (n *RepoNode) AbsPath() file.Path {
	return n.AbsParent().Append(file.Path(n.Name))
}

func (n *RepoNode) AbsParent() file.Path {
	return n.parent.AbsPath()
}

func (n *RepoNode) ServerSpec() *ServerSpec {
	return n.parent.ServerSpec()
}

func (n *RepoNode) UrlOrigin() string {
	return n.urlSpec(n.parent.nameOrigin)
}

func (n *RepoNode) UrlUpstream() string {
	return n.urlSpec(n.parent.nameUpstream)
}

func (n *RepoNode) urlSpec(o string) string {
	p := path.Join(o, n.Name) + ".git"
	if n.ServerSpec().Scheme() == pkg.SchemeHttps {
		// https://github.com/monopole/myrepos.git
		return "https://" + n.parent.parent.Domain().WithPort(n.ServerSpec().Port()) + "/" + p
	}
	// git@github.com:monopole/myrepos.git
	return "git@" + string(n.parent.parent.Domain()) + ":" + p
}

func (n *RepoNode) IsAFork() bool {
	return n.parent.nameOrigin != n.parent.nameUpstream
}

func MakeRepoNode(p *OrgNode, name pkg.RepoName) (n *RepoNode, err error) {
	return &RepoNode{
		parent: p,
		Name:   string(name),
	}, nil
}
