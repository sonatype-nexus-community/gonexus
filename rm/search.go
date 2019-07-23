package nexusrm

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sonatype-nexus-community/gonexus"
)

const (
	restSearchComponents = "service/rest/v1/search"
	restSearchAssets     = "service/rest/v1/search/assets"
	// restSearchAssetsDownload = "service/rest/v1/search/assets/download"
)

type searchComponentsResponse struct {
	Items             []RepositoryItem `json:"items"`
	ContinuationToken string           `json:"continuationToken"`
}

type searchAssetsResponse struct {
	Items             []RepositoryItemAsset `json:"items"`
	ContinuationToken string                `json:"continuationToken"`
}

func search(rm RM, endpoint string, queryBuilder nexus.SearchQueryBuilder, responseHandler func([]byte) (string, error)) error {
	continuation := ""

	get := func() (body []byte, err error) {
		url := endpoint + "?" + queryBuilder.Build()

		if continuation != "" {
			url += "&continuationToken=" + continuation
		}

		body, resp, err := rm.Get(url)
		if err != nil || resp.StatusCode != http.StatusOK {
			return
		}

		return
	}

	for {
		resp, err := get()
		if err != nil {
			return fmt.Errorf("could not find search items: %v", err)
		}

		continuation, err = responseHandler(resp)
		if err != nil {
			return fmt.Errorf("could not processes search items: %v", err)
		}

		if continuation == "" {
			break
		}
	}

	return nil
}

// SearchComponents allows searching the indicated RM instance for specific components
func SearchComponents(rm RM, query nexus.SearchQueryBuilder) ([]RepositoryItem, error) {
	items := make([]RepositoryItem, 0)

	err := search(rm, restSearchComponents, query, func(body []byte) (string, error) {
		var resp searchComponentsResponse
		if er := json.Unmarshal(body, &resp); er != nil {
			return "", er
		}

		items = append(items, resp.Items...)

		return resp.ContinuationToken, nil
	})

	return items, err
}

// SearchAssets allows searching the indicated RM instance for specific assets
func SearchAssets(rm RM, query nexus.SearchQueryBuilder) ([]RepositoryItemAsset, error) {
	items := make([]RepositoryItemAsset, 0)

	err := search(rm, restSearchAssets, query, func(body []byte) (string, error) {
		var resp searchAssetsResponse
		if er := json.Unmarshal(body, &resp); er != nil {
			return "", er
		}

		items = append(items, resp.Items...)

		return resp.ContinuationToken, nil
	})

	return items, err
}
