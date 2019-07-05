package nexusrm

import (
	"strings"
)

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

// RepositoryType enumerates the types of repositories in RM
/*
type RepositoryType int

const (
	proxy = iota
	hosted
	group
)

func (r RepositoryType) String() string {
	switch r {
	case proxy:
		return "proxy"
	case hosted:
		return "hosted"
	case group:
		return "group"
	default:
		return ""
	}
}
*/

// Repository collects the information returned by RM about a repository
type Repository struct {
	Name       string `json:"name"`
	Format     string `json:"format"`
	Type       string `json:"type"`
	URL        string `json:"url"`
	Attributes struct {
		Proxy struct {
			RemoteURL string `json:"remoteUrl"`
		} `json:"proxy"`
	} `json:"attributes,omitempty"`
}
