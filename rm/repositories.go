package nexusrm

import (
	"encoding/json"
	"net/http"

	"github.com/hokiegeek/gonexus"
)

const restListRepositories = "service/rest/v1/repositories"

// GetRepositories returns a list of components in the indicated repository
func GetRepositories(rm nexus.Server) (repos []Repository, err error) {
	body, resp, err := rm.Get(restListRepositories)
	if err != nil || resp.StatusCode != http.StatusOK {
		return
	}

	err = json.Unmarshal(body, &repos)

	return
}
