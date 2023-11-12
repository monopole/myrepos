package tree

import (
	"github.com/monopole/myrepos/internal/config"
	"github.com/monopole/myrepos/internal/file"
	"sort"
)

type OrgNode struct {
	parent       *ServerNode
	nameDir      file.Path
	nameOrigin   string
	nameUpstream string
	children     []*RepoNode
}

func (n *OrgNode) Accept(v Visitor) {
	v.VisitOrgNode(n)
	for _, c := range n.children {
		c.Accept(v)
	}
}

func (n *OrgNode) ServerSpec() *ServerSpec {
	return n.parent.ServerSpec()
}

func (n *OrgNode) NameDir() file.Path {
	return n.nameDir
}

func (n *OrgNode) NameOrigin() string {
	return n.nameOrigin
}

func (n *OrgNode) NameUpstream() string {
	return n.nameUpstream
}

func (n *OrgNode) AbsPath() file.Path {
	return n.parent.AbsPath().Append(n.nameDir)
}

func MakeOrgNode(p *ServerNode, orgName pkg.OrgName, names []pkg.RepoName) (on *OrgNode, err error) {
	dirName, origin, upstream := orgName.Parse()
	on = &OrgNode{
		parent:       p,
		nameDir:      dirName,
		nameOrigin:   string(origin),
		nameUpstream: string(upstream),
	}
	for _, repoName := range names {
		var rn *RepoNode
		rn, err = MakeRepoNode(on, repoName)
		if err != nil {
			return nil, err
		}
		on.children = append(on.children, rn)
	}
	sort.Slice(on.children, func(i, j int) bool {
		return on.children[i].Name < on.children[j].Name
	})
	return
}
