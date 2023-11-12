package visitor

import (
	"fmt"
	"github.com/monopole/myrepos/internal/runner"
	"github.com/monopole/myrepos/internal/tree"
)

type Cloner struct {
	Err      error
	FatalErr error
	// Output of most recent git command.
	gr *runner.Runner
}

func (v *Cloner) fatal(e error) {
	v.FatalErr, v.Err = e, e
}

func (v *Cloner) VisitRootNode(n *tree.RootNode) {
	if v.FatalErr != nil {
		return
	}
	exists, isDir := n.AbsPath().Exists()
	if !exists {
		// Make it for them instead of complain?
		v.fatal(fmt.Errorf("if you want the root dir %q, make it first", n.AbsPath()))
		return
	}
	if !isDir {
		v.fatal(fmt.Errorf("%q exists but isn't a directory", n.AbsPath()))
		return
	}
	indent(0)
	fmt.Println(n.AbsPath())
}

func (v *Cloner) VisitServerNode(n *tree.ServerNode) {
	if v.FatalErr != nil {
		return
	}
	if exists, isDir := n.AbsPath().Exists(); exists && !isDir {
		v.fatal(fmt.Errorf("%q exists but isn't a directory", n.AbsPath()))
		return
	}
	indent(1)
	fmt.Println(n.Domain())
}

func (v *Cloner) VisitOrgNode(n *tree.OrgNode) {
	if v.FatalErr != nil {
		return
	}
	if exists, isDir := n.AbsPath().Exists(); exists && !isDir {
		v.fatal(fmt.Errorf("%q exists but isn't a directory", n.AbsPath()))
		return
	}
	indent(2)
	fmt.Printf("%s|%s|%s\n", n.NameDir(), n.NameOrigin(), n.NameUpstream())
}

func (v *Cloner) VisitRepoNode(n *tree.RepoNode) {
	if v.FatalErr != nil {
		return
	}
	exists, isDir := n.AbsPath().Exists()
	if exists && !isDir {
		v.fatal(fmt.Errorf("%q exists but isn't a directory", n.AbsPath()))
		return
	}
	var (
		outcome Outcome
		err     error
	)
	if exists {
		outcome, err = v.Rebase(n)
	} else {
		outcome, err = v.Clone(n)
	}
	if err != nil {
		v.reportErr(n, err)
		return
	}
	var status string
	status, err = v.LastLog()
	if err != nil {
		v.reportErr(n, err)
		return
	}
	v.reportStatus(n, outcome, status)

	fmt.Println(n.Name)
}

func (v *Cloner) reportErr(n *tree.RepoNode, err error) {
	v.Err = err
	v.reportStatus(n, Oops, err.Error())
}

func (v *Cloner) reportStatus(n *tree.RepoNode, outcome Outcome, status string) {
	indent(3)
	fmt.Printf("%20s %20s %s\n", n.Name, outcome, status)
}
