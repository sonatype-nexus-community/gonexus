package nexusiq

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	nexus "github.com/sonatype-nexus-community/gonexus"
)

const restSearchComponent = "api/v2/search/component"

type searchResponse struct {
	Criteria criteria       `json:"criteria"`
	Results  []SearchResult `json:"results"`
}

type criteria struct {
	StageID             string              `json:"stageId"`
	Hash                string              `json:"hash"`
	PackageURL          string              `json:"packageUrl"`
	ComponentIdentifier ComponentIdentifier `json:"componentIdentifier"`
}

// SearchResult describes a component found based on search criteria
type SearchResult struct {
	ApplicationID       string              `json:"applicationId"`
	ApplicationName     string              `json:"applicationName"`
	ReportURL           string              `json:"reportUrl"`
	Hash                string              `json:"hash"`
	PackageURL          string              `json:"packageUrl"`
	ComponentIdentifier ComponentIdentifier `json:"componentIdentifier"`
}

// SearchQueryBuilder allows you to build a search query
type SearchQueryBuilder struct {
	criteria map[string]string
}

func (b *SearchQueryBuilder) addCriteria(c, v string) *SearchQueryBuilder {
	b.criteria[c] = v
	return b
}

func (b *SearchQueryBuilder) addCriteriaEncoded(c, v string) *SearchQueryBuilder {
	return b.addCriteria(c, url.QueryEscape(v))
}

// Build will build the assembled search query
func (b *SearchQueryBuilder) Build() string {
	var (
		buf      bytes.Buffer
		hasStage bool
	)

	for k, v := range b.criteria {
		if k == "stageId" {
			hasStage = true
		}
		buf.WriteString("&")
		buf.WriteString(k)
		buf.WriteString("=")
		buf.WriteString(v)
	}

	if !hasStage {
		buf.WriteString("&stageId=build")
	}

	return buf.String()
}

// Hash allows specifiying a sha1 hash to filter by
func (b *SearchQueryBuilder) Hash(v string) *SearchQueryBuilder {
	return b.addCriteria("hash", v)
}

// Format allows specifiying a format to filter by
func (b *SearchQueryBuilder) Format(v string) *SearchQueryBuilder {
	return b.addCriteria("format", v)
}

// ComponentIdentifier allows specifiying a component identifier to filter by
func (b *SearchQueryBuilder) ComponentIdentifier(c ComponentIdentifier) *SearchQueryBuilder {
	v, err := json.Marshal(c)
	if err != nil {
		return b
	}
	return b.addCriteriaEncoded("componentIdentifier", string(v))
}

// PackageURL allows specifiying a purl to filter by
func (b *SearchQueryBuilder) PackageURL(v string) *SearchQueryBuilder {
	return b.addCriteriaEncoded("packageUrl", v)
}

// Coordinates allows specifiying component coordinates to filter by
func (b *SearchQueryBuilder) Coordinates(c Coordinates) *SearchQueryBuilder {
	v, err := json.Marshal(c)
	if err != nil {
		return b
	}
	return b.addCriteriaEncoded("coordinates", string(v))
}

// Stage allows specifiying a stage to filter by
func (b *SearchQueryBuilder) Stage(v string) *SearchQueryBuilder {
	return b.addCriteria("stageId", v)
}

// NewSearchQueryBuilder creates a new instance of SearchQueryBuilder
func NewSearchQueryBuilder() *SearchQueryBuilder {
	b := new(SearchQueryBuilder)
	b.criteria = make(map[string]string)
	return b
}

// SearchComponents allows searching the indicated IQ instance for specific components
func SearchComponents(iq IQ, query nexus.SearchQueryBuilder) ([]SearchResult, error) {
	endpoint := restSearchComponent + "?" + query.Build()
	body, resp, err := iq.Get(endpoint)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("could not find component: %v", err)
	}

	var searchResp searchResponse
	if err := json.Unmarshal(body, &searchResp); err != nil {
		return nil, fmt.Errorf("could not process response: %v", err)
	}

	return searchResp.Results, nil
}
