package tree

type Visitor interface {
	VisitRootNode(n *RootNode)
	VisitServerNode(n *ServerNode)
	VisitOrgNode(n *OrgNode)
	VisitRepoNode(n *RepoNode)
}
