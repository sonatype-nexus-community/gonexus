package nexusrm

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	restAssets           = "service/rest/v1/assets"
	restListAssetsByRepo = "service/rest/v1/assets?repository=%s"
)

type repositoryItemAssetsChecksum struct {
	Sha1 string `json:"sha1"`
	Md5  string `json:"md5"`
}

// RepositoryItemAsset describes the assets associated with a component
type RepositoryItemAsset struct {
	DownloadURL string                       `json:"downloadUrl"`
	Path        string                       `json:"path"`
	ID          string                       `json:"id"`
	Repository  string                       `json:"repository"`
	Format      string                       `json:"format"`
	Checksum    repositoryItemAssetsChecksum `json:"checksum"`
}

type listAssetsResponse struct {
	Items             []RepositoryItemAsset `json:"items"`
	ContinuationToken string                `json:"continuationToken"`
}

// GetAssets returns a list of assets in the indicated repository
func GetAssets(rm RM, repo string) (items []RepositoryItemAsset, err error) {
	continuation := ""

	get := func() (listResp listAssetsResponse, err error) {
		url := fmt.Sprintf(restListAssetsByRepo, repo)

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

	items = make([]RepositoryItemAsset, 0)
	for {
		resp, err := get()
		if err != nil {
			return items, fmt.Errorf("could not get assets: %v", err)
		}

		items = append(items, resp.Items...)

		if resp.ContinuationToken == "" {
			break
		}

		continuation = resp.ContinuationToken
	}

	return items, nil
}

// GetAssetByID returns an asset by ID
func GetAssetByID(rm RM, id string) (items RepositoryItemAsset, err error) {
	doError := func(err error) error {
		return fmt.Errorf("no asset with id '%s': %v", id, err)
	}

	var item RepositoryItemAsset

	url := fmt.Sprintf("%s/%s", restAssets, id)
	body, resp, err := rm.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		return item, doError(err)
	}

	if err := json.Unmarshal(body, &item); err != nil {
		return item, doError(err)
	}

	return item, nil
}

// DeleteAssetByID deletes the asset indicated by ID
func DeleteAssetByID(rm RM, id string) error {
	url := fmt.Sprintf("%s/%s", restAssets, id)

	if resp, err := rm.Del(url); err != nil && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("asset not deleted '%s': %v", id, err)
	}

	return nil
}
