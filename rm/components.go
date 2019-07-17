package nexusrm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

const (
	restComponents           = "service/rest/v1/components"
	restListComponentsByRepo = "service/rest/v1/components?repository=%s"
)

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

// GetComponentByID returns a component by ID
func GetComponentByID(rm RM, id string) (item RepositoryItem, err error) {
	url := fmt.Sprintf("%s/%s", restComponents, id)

	body, resp, err := rm.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		return
	}

	err = json.Unmarshal(body, &item)

	return
}

// DeleteComponentByID uploads a component to repository manager
func DeleteComponentByID(rm RM, id string) error {
	url := fmt.Sprintf("%s/%s", restComponents, id)

	if resp, err := rm.Del(url); err != nil && resp.StatusCode != http.StatusNoContent {
		return err
	}

	return nil
}

// UploadComponent uploads a component to repository manager
func UploadComponent(rm RM, repo string, component uploadComponent) error {
	fields, files := component.formData()

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	defer w.Close()

	for k, v := range fields {
		if v != "" {
			w.WriteField(k, v)
		}
	}

	for k, v := range files {
		fw, err := w.CreateFormFile(k, v.Name())
		if err != nil {
			return err
		}
		// fw.Write(

		if _, err = io.Copy(fw, v); err != nil {
			return err
		}
	}

	if err := w.Close(); err != nil {
		return err
	}

	url := fmt.Sprintf(restListComponentsByRepo, repo)
	req, err := rm.NewRequest("POST", url, &b)
	req.Header.Set("Content-Type", w.FormDataContentType())
	if err != nil {
		return err
	}

	_, resp, err := rm.Do(req)
	// _, resp, err := rm.Post(url, &b)
	if err != nil && resp.StatusCode != http.StatusNoContent {
		return err
	}

	return nil
}
