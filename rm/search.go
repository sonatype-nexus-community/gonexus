package nexusrm

import (
	"bytes"
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
	sort      SearchSort
	direction SearchSortDirection
	criteria  map[string]string
}

func (b *SearchQueryBuilder) addCriteria(c, v string) *SearchQueryBuilder {
	b.criteria[c] = v
	return b
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

	for k, v := range b.criteria {
		buf.WriteString("&")
		buf.WriteString(k)
		buf.WriteString("=")
		buf.WriteString(v)
	}

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

// Q allows specifying a keyword search
func (b *SearchQueryBuilder) Q(v string) *SearchQueryBuilder {
	return b.addCriteria("q", v)
}

// Repository allows specifying the repository to search
func (b *SearchQueryBuilder) Repository(v string) *SearchQueryBuilder {
	return b.addCriteria("repository", v)
}

// Format allows specifiying the format to filter by
func (b *SearchQueryBuilder) Format(v string) *SearchQueryBuilder {
	return b.addCriteria("format", v)
}

// Tag allows specifiying a tag to filter by
func (b *SearchQueryBuilder) Tag(v string) *SearchQueryBuilder {
	return b.addCriteria("tag", v)
}

// Group allows specifiying a group to filter by
func (b *SearchQueryBuilder) Group(v string) *SearchQueryBuilder {
	return b.addCriteria("group", v)
}

// Name allows specifiying a name to filter by
func (b *SearchQueryBuilder) Name(v string) *SearchQueryBuilder {
	return b.addCriteria("name", v)
}

// Version allows specifiying a version to filter by
func (b *SearchQueryBuilder) Version(v string) *SearchQueryBuilder {
	return b.addCriteria("version", v)
}

// Md5 allows specifiying an md5 sum to filter by
func (b *SearchQueryBuilder) Md5(v string) *SearchQueryBuilder {
	return b.addCriteria("md5", v)
}

// Sha1 allows specifiying an sha1 sum to filter by
func (b *SearchQueryBuilder) Sha1(v string) *SearchQueryBuilder {
	return b.addCriteria("sha1", v)
}

// Sha256 allows specifiying an sha256 sum to filter by
func (b *SearchQueryBuilder) Sha256(v string) *SearchQueryBuilder {
	return b.addCriteria("sha256", v)
}

// Sha512 allows specifiying an sha512 sum to filter by
func (b *SearchQueryBuilder) Sha512(v string) *SearchQueryBuilder {
	return b.addCriteria("sha512", v)
}

// Prerelease allows specifiying a prerelease qualifier to filter by
func (b *SearchQueryBuilder) Prerelease(v string) *SearchQueryBuilder {
	return b.addCriteria("prerelease", v)
}

// DockerImageName allows specifiying the name of a docker image to filter by
func (b *SearchQueryBuilder) DockerImageName(v string) *SearchQueryBuilder {
	return b.addCriteria("docker.imageName", v)
}

// DockerImageTag allows specifiying the tag of a docker image to filter by
func (b *SearchQueryBuilder) DockerImageTag(v string) *SearchQueryBuilder {
	return b.addCriteria("docker.imageTag", v)
}

// DockerLayerID allows specifiying the ID of a docker image layer to filter by
func (b *SearchQueryBuilder) DockerLayerID(v string) *SearchQueryBuilder {
	return b.addCriteria("docker.layerId", v)
}

// DockerContentDigest allows specifiying the digest of docker layers to filter by
func (b *SearchQueryBuilder) DockerContentDigest(v string) *SearchQueryBuilder {
	return b.addCriteria("docker.contentDigest", v)
}

// MavenGroupID allows specifiying the group name/id of maven component to filter by
func (b *SearchQueryBuilder) MavenGroupID(v string) *SearchQueryBuilder {
	return b.addCriteria("maven.groupId", v)
}

// MavenArtifactID allows specifiying the artifact id of maven component to filter by
func (b *SearchQueryBuilder) MavenArtifactID(v string) *SearchQueryBuilder {
	return b.addCriteria("maven.artifactId", v)
}

// MavenBaseVersion allows specifiying the version of maven component to filter by
func (b *SearchQueryBuilder) MavenBaseVersion(v string) *SearchQueryBuilder {
	return b.addCriteria("maven.baseVersion", v)
}

// MavenExtension allows specifiying the extension of maven component to filter by
func (b *SearchQueryBuilder) MavenExtension(v string) *SearchQueryBuilder {
	return b.addCriteria("maven.extension", v)
}

// MavenClassifier allows specifiying the classifier of maven component to filter by
func (b *SearchQueryBuilder) MavenClassifier(v string) *SearchQueryBuilder {
	return b.addCriteria("maven.classifier", v)
}

// NpmScope allows specifiying the scope of an NPM component to filter by
func (b *SearchQueryBuilder) NpmScope(v string) *SearchQueryBuilder {
	return b.addCriteria("npm.scope", v)
}

// NugetID allows specifiying the ID/name of a Nuget component to filter by
func (b *SearchQueryBuilder) NugetID(v string) *SearchQueryBuilder {
	return b.addCriteria("nuget.id", v)
}

// NugetTags allows specifiying the tags of a Nuget component to filter by
func (b *SearchQueryBuilder) NugetTags(v string) *SearchQueryBuilder {
	return b.addCriteria("nuget.tags", v)
}

// PypiClassifiers allows specifiying the classifiers of a pypi component to filter by
func (b *SearchQueryBuilder) PypiClassifiers(v string) *SearchQueryBuilder {
	return b.addCriteria("pypi.classifiers", v)
}

// PypiDescription allows specifiying the description of a pypi component to filter by
func (b *SearchQueryBuilder) PypiDescription(v string) *SearchQueryBuilder {
	return b.addCriteria("pypi.description", v)
}

// PypiKeywords allows specifiying the keywords of a pypi component to filter by
func (b *SearchQueryBuilder) PypiKeywords(v string) *SearchQueryBuilder {
	return b.addCriteria("pypi.keywords", v)
}

// PypiSummary allows specifiying the summary of a pypi component to filter by
func (b *SearchQueryBuilder) PypiSummary(v string) *SearchQueryBuilder {
	return b.addCriteria("pypi.summary", v)
}

// RubygemsDescription allows specifiying the description of a ruby gem to filter by
func (b *SearchQueryBuilder) RubygemsDescription(v string) *SearchQueryBuilder {
	return b.addCriteria("rubygems.description", v)
}

// RubygemsPlatform allows specifiying the platform of a ruby gem to filter by
func (b *SearchQueryBuilder) RubygemsPlatform(v string) *SearchQueryBuilder {
	return b.addCriteria("rubygems.platform", v)
}

// RubygemsSummary allows specifiying the summary of a ruby gem to filter by
func (b *SearchQueryBuilder) RubygemsSummary(v string) *SearchQueryBuilder {
	return b.addCriteria("rubygems.summary", v)
}

// YumArchitecture allows specifiying the architecture of a Yum package to filter by
func (b *SearchQueryBuilder) YumArchitecture(v string) *SearchQueryBuilder {
	return b.addCriteria("yum.architecture", v)
}

// NewSearchQueryBuilder creates a new instance of SearchQueryBuilder
func NewSearchQueryBuilder() *SearchQueryBuilder {
	b := new(SearchQueryBuilder)
	b.criteria = make(map[string]string)
	return b
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
