package nexusrm

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hokiegeek/gonexus"
)

// http://localhost:8081/service/rest/v1/components?continuationToken=foo&repository=bar
const restListComponentsByRepo = "service/rest/v1/components?repository=%s"
const restListRepositories = "service/rest/v1/repositories"

const hashPart = 20

// RM holds basic and state info of the Repository Manager server we will connect to
type RM struct {
	nexus.Server
}

// GetComponents returns a list of components in the indicated repository
func (rm *RM) GetComponents(repo string) (items []RepositoryItem, err error) {
	continuation := ""

	getComponents := func() (listResp listComponentsResponse, err error) {
		url := fmt.Sprintf(restListComponentsByRepo, repo)

		if continuation != "" {
			url += "&continuationToken=" + continuation
		}

		body, resp, err := rm.Get(url)
		if err != nil || resp.StatusCode != http.StatusOK {
			return
		}

		err = json.Unmarshal(body, &listResp)

		return
	}

	for {
		resp, err := getComponents()
		if err != nil {
			return items, err
		}

		items = append(items, resp.Items...)

		if resp.ContinuationToken == "" {
			break
		}

		continuation = resp.ContinuationToken
	}

	return
}

// GetRepositories returns a list of components in the indicated repository
func (rm *RM) GetRepositories() (repos []Repository, err error) {
	body, resp, err := rm.Get(restListRepositories)
	if err != nil || resp.StatusCode != http.StatusOK {
		return
	}

	err = json.Unmarshal(body, &repos)

	return
}

// New creates a new Repository Manager instance
func New(host, username, password string) (rm *RM, err error) {
	rm = new(RM)
	rm.Host = host
	rm.Username = username
	rm.Password = password
	return
}
