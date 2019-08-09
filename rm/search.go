package nexusrm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	nexus "github.com/sonatype-nexus-community/gonexus"
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

// SearchSort enumerates the sort options allowed
type SearchSort int

// Available sort options
const (
	None SearchSort = iota
	Group
	Name
	Version
	Repo
)

// SearchSortDirection enumerates the direction of the sort
type SearchSortDirection int

// Can be ascending (Asc) or descending (Desc)
const (
	Asc SearchSortDirection = iota
	Desc
)

// SearchQueryBuilder allows you to build a search query
type SearchQueryBuilder struct {
	QueryBuilder
	sort      SearchSort
	direction SearchSortDirection
}

// Build will build the assembled search query
func (b *SearchQueryBuilder) Build() string {
	var buf bytes.Buffer

	if b.sort != None {
		buf.WriteString("sort=")
		switch b.sort {
		case Group:
			buf.WriteString("group")
		case Name:
			buf.WriteString("name")
		case Version:
			buf.WriteString("version")
		case Repo:
			buf.WriteString("repository")
		}
		buf.WriteString("&")
	}

	buf.WriteString("direction=")
	switch b.direction {
	case Asc:
		buf.WriteString("asc")
	case Desc:
		buf.WriteString("desc")
	}

	b.buildCriteria(&buf)

	return buf.String()
}

// Sort allows specifying how to sort the data (Defaults to "magic")
func (b *SearchQueryBuilder) Sort(v SearchSort) *SearchQueryBuilder {
	b.sort = v
	return b
}

// Direction allows specifying the direction to sort (Defaults to Asc)
func (b *SearchQueryBuilder) Direction(v SearchSortDirection) *SearchQueryBuilder {
	b.direction = v
	return b
}

// NewSearchQueryBuilder creates a new instance of SearchQueryBuilder
func NewSearchQueryBuilder() *SearchQueryBuilder {
	b := new(SearchQueryBuilder)
	b.criteria = make(map[string]string)
	return b
}

func search(rm RM, endpoint string, queryBuilder nexus.SearchQueryBuilder, responseHandler func([]byte) (string, error)) error {
	continuation := ""
	queryEndpoint := fmt.Sprintf("%s?%s", endpoint, queryBuilder.Build())

	get := func() (body []byte, err error) {
		url := queryEndpoint

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
	fmt.Println("query", query.Build())
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
