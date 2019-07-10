package nexusrm

import (
	"encoding/json"
	"net/http"
)

const restListRepositories = "service/rest/v1/repositories"

// GetRepositories returns a list of components in the indicated repository
func GetRepositories(rm RM) (repos []Repository, err error) {
	body, resp, err := rm.Get(restListRepositories)
	if err != nil || resp.StatusCode != http.StatusOK {
		return
	}

	err = json.Unmarshal(body, &repos)

	return
}
