package nexusrm

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const restRepositories = "service/rest/v1/repositories"

/*
// RepositoryType enumerates the types of repositories in RM
type repositoryType int

const (
	Hosted = iota
	Proxy
	Group
)

func (r repositoryType) String() string {
	switch r {
	case proxy:
		return "proxy"
	case hosted:
		return "hosted"
	case group:
		return "group"
	default:
		return ""
	}
}
*/

type repositoryFormat int

// Enumerates the formats which can be created as Repository Manager repositories
const (
	Unknown repositoryFormat = iota
	Maven
	Npm
	Nuget
	Apt
	Docker
	Golang
	Raw
	Rubygems
	Bower
	Pypi
	Yum
	GitLfs
)

// Repository collects the information returned by RM about a repository
type Repository struct {
	Name       string `json:"name"`
	Format     string `json:"format"`
	Type       string `json:"type"`
	URL        string `json:"url"`
	Attributes struct {
		Proxy struct {
			RemoteURL string `json:"remoteUrl"`
		} `json:"proxy"`
	} `json:"attributes,omitempty"`
}

// GetRepositories returns a list of components in the indicated repository
func GetRepositories(rm RM) ([]Repository, error) {
	doError := func(err error) error {
		return fmt.Errorf("could not find repositories: %v", err)
	}

	body, resp, err := rm.Get(restRepositories)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, doError(err)
	}

	repos := make([]Repository, 0)
	if err := json.Unmarshal(body, &repos); err != nil {
		return nil, doError(err)
	}

	return repos, nil
}

// GetRepositoryByName returns information on a named repository
func GetRepositoryByName(rm RM, name string) (repo Repository, err error) {
	repos, err := GetRepositories(rm)
	if err != nil {
		return repo, fmt.Errorf("could not get list of repositories: %v", err)
	}

	for _, repo = range repos {
		if repo.Name == name {
			return
		}
	}

	return repo, fmt.Errorf("did not find repository '%s': %v", name, err)
}
