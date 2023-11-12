package tree

import (
	"github.com/monopole/myrepos/internal/config"
	"github.com/monopole/myrepos/internal/file"
	"sort"
)

type ServerNode struct {
	parent   *RootNode
	spec     *ServerSpec
	domain   pkg.ServerDomain
	children []*OrgNode
}

func (n *ServerNode) Accept(v Visitor) {
	v.VisitServerNode(n)
	for _, c := range n.children {
		c.Accept(v)
	}
}

func (n *ServerNode) AbsPath() file.Path {
	return n.parent.AbsPath().Append(file.Path(n.domain))
}

func (n *ServerNode) ServerSpec() *ServerSpec {
	return n.spec
}

func (n *ServerNode) Domain() pkg.ServerDomain {
	return n.domain
}

func MakeServerNode(
	p *RootNode, n pkg.ServerDomain,
	spec *ServerSpec, orgMap map[pkg.OrgName][]pkg.RepoName) (sn *ServerNode, err error) {
	sn = &ServerNode{
		parent: p,
		spec:   spec,
		domain: n,
	}
	for orgName, repoList := range orgMap {
		var on *OrgNode
		on, err = MakeOrgNode(sn, orgName, repoList)
		if err != nil {
			return nil, err
		}
		sn.children = append(sn.children, on)
	}
	sort.Slice(sn.children, func(i, j int) bool {
		return sn.children[i].nameDir < sn.children[j].nameDir
	})
	return
}
