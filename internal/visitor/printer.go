package visitor

import (
	"fmt"
	"github.com/monopole/myrepos/internal/tree"
)

type Printer struct {
	Err error
}

func (v *Printer) VisitRootNode(n *tree.RootNode) {
	fmt.Println(indent(0), n.AbsPath())
}

func (v *Printer) VisitServerNode(n *tree.ServerNode) {
	fmt.Println(indent(1), n.Domain())
}

func (v *Printer) VisitOrgNode(n *tree.OrgNode) {
	fmt.Print(indent(2))
	fmt.Printf("%s|%s|%s\n", n.NameDir(), n.NameOrigin(), n.NameUpstream())
}

func (v *Printer) VisitRepoNode(n *tree.RepoNode) {
	fmt.Println(indent(3), n.Name)
}
