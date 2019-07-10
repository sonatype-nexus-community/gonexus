package nexusrm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// http://localhost:8081/service/rest/v1/components?continuationToken=foo&repository=bar
const restListComponentsByRepo = "service/rest/v1/components?repository=%s"

// RepositoryItemAssets describes the assets associated with a component
type RepositoryItemAssets struct {
	DownloadURL string `json:"downloadUrl"`
	Path        string `json:"path"`
	ID          string `json:"id"`
	Repository  string `json:"repository"`
	Format      string `json:"format"`
	Checksum    struct {
		Sha1 string `json:"sha1"`
		Md5  string `json:"md5"`
	} `json:"checksum"`
}

// Equals compares two RepositoryItemAssets objects
func (a *RepositoryItemAssets) Equals(b *RepositoryItemAssets) (_ bool) {
	if a == b {
		return true
	}

	if a.DownloadURL != b.DownloadURL {
		return
	}

	if a.Path != b.Path {
		return
	}

	if a.ID != b.ID {
		return
	}

	if a.Repository != b.Repository {
		return
	}

	if a.Format != b.Format {
		return
	}

	if a.Checksum.Sha1 != b.Checksum.Sha1 {
		return
	}

	if a.Checksum.Md5 != b.Checksum.Md5 {
		return
	}

	return true
}

// RepositoryItem holds the data of a component in a repository
type RepositoryItem struct {
	ID         string                 `json:"id"`
	Repository string                 `json:"repository"`
	Format     string                 `json:"format"`
	Group      string                 `json:"group"`
	Name       string                 `json:"name"`
	Version    string                 `json:"version"`
	Assets     []RepositoryItemAssets `json:"assets"`
	Tags       []interface{}          `json:"tags"`
}

// Equals compares two RepositoryItem objects
func (a *RepositoryItem) Equals(b *RepositoryItem) (_ bool) {
	if a == b {
		return true
	}

	if a.ID != b.ID {
		return
	}

	if a.Repository != b.Repository {
		return
	}

	if a.Format != b.Format {
		return
	}

	if a.Group != b.Group {
		return
	}

	if a.Name != b.Name {
		return
	}

	if a.Version != b.Version {
		return
	}

	if len(a.Assets) != len(b.Assets) {
		return
	}

	for i, asset := range a.Assets {
		if !asset.Equals(&b.Assets[i]) {
			return
		}
	}

	return true
}

const hashPart = 20

// Hash is a hack which returns the most appopriate IQable hash of a repo item
func (a *RepositoryItem) Hash() string {
	var hash string

	sumByExt := func(exts []string) string {
		ext := exts[0]
		for _, ass := range a.Assets {
			if strings.HasSuffix(ass.Path, ext) {
				return ass.Checksum.Sha1
			}
		}
		return ""
	}

	switch a.Format {
	case "maven2":
		hash = sumByExt([]string{"jar"})
	case "rubygems":
		hash = sumByExt([]string{"gem"})
	case "npm":
		hash = sumByExt([]string{"tar.gz"})
	case "pipy":
		hash = sumByExt([]string{"whl", "tar.gz"})
	default:
		hash = ""
	}
	if len(hash) < hashPart {
		return hash
	}
	return hash[0:hashPart]
}

type listComponentsResponse struct {
	Items             []RepositoryItem `json:"items"`
	ContinuationToken string           `json:"continuationToken"`
}

// GetComponents returns a list of components in the indicated repository
func GetComponents(rm RM, repo string) (items []RepositoryItem, err error) {
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
