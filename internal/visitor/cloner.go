package visitor

import (
	"fmt"
	"strconv"

	"github.com/TwiN/go-color"
	"github.com/monopole/myrepos/internal/runner"
	"github.com/monopole/myrepos/internal/tree"
)

const (
	repoNameFieldSize = 40
	spaces            = "                             "
)

var (
	fmtHeader = "%-" + strconv.Itoa(repoNameFieldSize) + "s"
	fmtReport = "%" + strconv.Itoa(repoNameFieldSize) + "s%30s %s\n"
)

func indent(i int) string {
	return spaces[:i*2]
}

func header(i int, arg string) string {
	return fmt.Sprintf(fmtHeader, indent(i)+arg)
}

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
	fmt.Println(color.InBlackOverGray(header(0, string(n.AbsPath()))))
}

func (v *Cloner) VisitServerNode(n *tree.ServerNode) {
	// TODO: Use some simple git command, e.g.
	//  git remote show origin
	// to comfirm that the server is reachable.
	if v.fatalErr != nil {
		return
	}
	if exists, isDir := n.AbsPath().Exists(); exists && !isDir {
		v.fatal(fmt.Errorf("%q exists but isn't a directory", n.AbsPath()))
		return
	}
	fmt.Println(color.InBlackOverGreen(header(1, string(n.Domain()))))
}

func (v *Cloner) VisitOrgNode(n *tree.OrgNode) {
	if v.fatalErr != nil {
		return
	}
	if exists, isDir := n.AbsPath().Exists(); exists && !isDir {
		v.fatal(fmt.Errorf("%q exists but isn't a directory", n.AbsPath()))
		return
	}
	fmt.Println(color.InBlackOverBlue(header(2, string(n.NameDir()))))
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
	fmt.Printf(fmtReport, n.Name, outcome, status)
}
