package visitor

import (
	"fmt"
	"github.com/TwiN/go-color"
	"github.com/monopole/myrepos/internal/runner"
	"github.com/monopole/myrepos/internal/tree"
)

type Cloner struct {
	lastErr  error
	fatalErr error
	gr       *runner.Runner
}

func (v *Cloner) fatal(e error) {
	v.fatalErr, v.lastErr = e, e
}

func (v *Cloner) Err() error {
	return v.lastErr
}

func (v *Cloner) VisitRootNode(n *tree.RootNode) {
	if v.fatalErr != nil {
		return
	}
	exists, isDir := n.AbsPath().Exists()
	if !exists {
		// Make it for instead of complain?
		// Could trigger a bunch of work if it's just a typo.
		v.fatal(fmt.Errorf("if you want the root dir %q, make it first", n.AbsPath()))
		return
	}
	if !isDir {
		v.fatal(fmt.Errorf("%q exists but isn't a directory", n.AbsPath()))
		return
	}
	indent(0)
	fmt.Println(color.InBlackOverWhite(n.AbsPath()))
}

func (v *Cloner) VisitServerNode(n *tree.ServerNode) {
	if v.fatalErr != nil {
		return
	}
	if exists, isDir := n.AbsPath().Exists(); exists && !isDir {
		v.fatal(fmt.Errorf("%q exists but isn't a directory", n.AbsPath()))
		return
	}
	indent(1)
	fmt.Println(color.YellowBackground, n.Domain(), color.Reset)
}

func (v *Cloner) VisitOrgNode(n *tree.OrgNode) {
	if v.fatalErr != nil {
		return
	}
	if exists, isDir := n.AbsPath().Exists(); exists && !isDir {
		v.fatal(fmt.Errorf("%q exists but isn't a directory", n.AbsPath()))
		return
	}
	indent(2)
	fmt.Println(color.CyanBackground, n.NameDir(), color.Reset)
}

func (v *Cloner) VisitRepoNode(n *tree.RepoNode) {
	if v.fatalErr != nil {
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
}

func (v *Cloner) reportErr(n *tree.RepoNode, err error) {
	v.lastErr = err
	v.reportStatus(n, Oops, err.Error())
}

func (v *Cloner) reportStatus(n *tree.RepoNode, outcome Outcome, status string) {
	indent(3)
	fmt.Printf("%30s%30s %s\n", n.Name, outcome, status)
}
