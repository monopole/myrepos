package tree

import (
	"github.com/monopole/myrepos/internal/config"
	"github.com/monopole/myrepos/internal/file"
	"sort"
)

type RootNode struct {
	absPath  file.Path
	children []*ServerNode
}

func (n *RootNode) Accept(v Visitor) {
	v.VisitRootNode(n)
	for _, c := range n.children {
		c.Accept(v)
	}
}

func (n *RootNode) AbsPath() file.Path {
	return n.absPath
}

func MakeRootNode(c *pkg.Config) (rn *RootNode, err error) {
	rn = &RootNode{
		absPath: absRootDir(c),
	}
	serverSpecs := make(map[pkg.ServerDomain]*ServerSpec)
	for d, opts := range c.ServerOpts {
		serverSpecs[d], err = FromServerOpts(&opts)
		if err != nil {
			return nil, err
		}
	}
	for domain, orgMap := range c.Layout {
		if _, ok := serverSpecs[domain]; !ok {
			serverSpecs[domain] = MakeServerSpec()
		}
		var sn *ServerNode
		sn, err = MakeServerNode(rn, domain, serverSpecs[domain], orgMap)
		if err != nil {
			return nil, err
		}
		rn.children = append(rn.children, sn)
	}
	sort.Slice(rn.children, func(i, j int) bool {
		return rn.children[i].Domain() < rn.children[j].Domain()
	})
	return rn, nil
}

func absRootDir(c *pkg.Config) file.Path {
	if c.Path == "" {
		return file.Home()
	}
	if c.Path.IsAbs() {
		return c.Path
	}
	return file.Home().Append(c.Path)
}
