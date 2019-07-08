package nexusrm

import (
	"encoding/json"
	"fmt"
	"net/http"

	nexus "github.com/hokiegeek/gonexus"
)

// http://localhost:8081/service/rest/v1/components?continuationToken=foo&repository=bar
const restListComponentsByRepo = "service/rest/v1/components?repository=%s"

// GetComponents returns a list of components in the indicated repository
func GetComponents(rm nexus.Server, repo string) (items []RepositoryItem, err error) {
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
