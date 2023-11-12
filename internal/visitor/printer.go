package visitor

import (
	"fmt"
	"github.com/monopole/myrepos/internal/tree"
)

type Printer struct {
	Err error
}

func (v *Printer) VisitRootNode(n *tree.RootNode) {
	indent(0)
	fmt.Println(n.AbsPath())
}

func (v *Printer) VisitServerNode(n *tree.ServerNode) {
	indent(1)
	fmt.Println(n.Domain())
}

func (v *Printer) VisitOrgNode(n *tree.OrgNode) {
	indent(2)
	fmt.Printf("%s|%s|%s\n", n.NameDir(), n.NameOrigin(), n.NameUpstream())
}

func (v *Printer) VisitRepoNode(n *tree.RepoNode) {
	indent(3)
	fmt.Println(n.Name)
}
