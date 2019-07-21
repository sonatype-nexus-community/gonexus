package nexusrm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
)

const (
	restComponents           = "service/rest/v1/components"
	restListComponentsByRepo = "service/rest/v1/components?repository=%s"
)

type listComponentsResponse struct {
	Items             []RepositoryItem `json:"items"`
	ContinuationToken string           `json:"continuationToken"`
}

// GetComponents returns a list of components in the indicated repository
func GetComponents(rm RM, repo string) ([]RepositoryItem, error) {
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

	items := make([]RepositoryItem, 0)
	for {
		resp, err := getComponents()
		if err != nil {
			return items, fmt.Errorf("could not get components: %v", err)
		}

		items = append(items, resp.Items...)

		if resp.ContinuationToken == "" {
			break
		}

		continuation = resp.ContinuationToken
	}

	return items, nil
}

// GetComponentByID returns a component by ID
func GetComponentByID(rm RM, id string) (RepositoryItem, error) {
	doError := func(err error) error {
		return fmt.Errorf("no component with id '%s': %v", id, err)
	}

	var item RepositoryItem

	url := fmt.Sprintf("%s/%s", restComponents, id)
	body, resp, err := rm.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		return item, doError(err)
	}

	if err := json.Unmarshal(body, &item); err != nil {
		return item, doError(err)
	}

	return item, nil
}

// DeleteComponentByID deletes the indicated component
func DeleteComponentByID(rm RM, id string) error {
	url := fmt.Sprintf("%s/%s", restComponents, id)

	if resp, err := rm.Del(url); err != nil && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("component not deleted '%s': %v", id, err)
	}

	return nil
}

// UploadComponent uploads a component to repository manager
func UploadComponent(rm RM, repo string, component UploadComponentWriter) error {
	if _, err := GetRepositoryByName(rm, repo); err != nil {
		return fmt.Errorf("could not find repository: %v", err)
	}

	doError := func(err error) error {
		return fmt.Errorf("component not uploaded: %v", err)
	}

	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	component.write(w)

	if err := w.Close(); err != nil {
		return doError(err)
	}

	url := fmt.Sprintf(restListComponentsByRepo, repo)
	req, err := rm.NewRequest("POST", url, &b)
	req.Header.Set("Content-Type", w.FormDataContentType())
	if err != nil {
		return doError(err)
	}

	if _, resp, err := rm.Do(req); err != nil && resp.StatusCode != http.StatusNoContent {
		return doError(err)
	}

	return nil
}
