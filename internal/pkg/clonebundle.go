package pkg

type CloneBundle struct {
	Server ServerDomain
	Repos  map[OrgName][]RepoName
}
