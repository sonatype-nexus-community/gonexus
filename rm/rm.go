package nexusrm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hokiegeek/gonexus"
)

// http://localhost:8081/service/rest/v1/components?continuationToken=foo&repository=bar
const restListComponentsByRepo = "service/rest/v1/components?repository=%s"

const hashPart = 20

// RepositoryItem holds the data of a component in a repository
type RepositoryItem struct {
	ID         string `json:"id"`
	Repository string `json:"repository"`
	Format     string `json:"format"`
	Group      string `json:"group"`
	Name       string `json:"name"`
	Version    string `json:"version"`
	Assets     []struct {
		DownloadURL string `json:"downloadUrl"`
		Path        string `json:"path"`
		ID          string `json:"id"`
		Repository  string `json:"repository"`
		Format      string `json:"format"`
		Checksum    struct {
			Sha1 string `json:"sha1"`
			Md5  string `json:"md5"`
		} `json:"checksum"`
	} `json:"assets"`
	Tags []interface{} `json:"tags"`
}

func (i *RepositoryItem) hash() string {
	var hash string

	sumByExt := func(exts []string) string {
		ext := exts[0]
		for _, ass := range i.Assets {
			if strings.HasSuffix(ass.Path, ext) {
				return ass.Checksum.Sha1
			}
		}
		return ""
	}

	switch i.Format {
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

// RM holds basic and state info of the Repository Manager server we will connect to
type RM struct {
	nexus.Server
}

// ListComponents returns a list of components in the indicated repository
func (rm *RM) ListComponents(repo string) (items []RepositoryItem, err error) {
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

		if err = json.Unmarshal(body, &listResp); err != nil {
			return
		}

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

// New creates a new Repository Manager instance
func New(host, username, password string) (rm *RM, err error) {
	rm = new(RM)
	rm.Host = host
	rm.Username = username
	rm.Password = password
	return
}
