package nexusiq

import (
	"encoding/json"
)

const restOrganization = "api/v2/organizations"

type iqNewOrgRequest struct {
	Name string `json:"name"`
}

type iqNewOrgResponse struct {
	ID   string   `json:"id"`
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

// CreateOrganization creates an organization in IQ with the given name
func CreateOrganization(iq *IQ, name string) (string, error) {
	request, err := json.Marshal(iqNewOrgRequest{Name: name})
	if err != nil {
		return "", err
	}

	body, _, err := iq.Post(restOrganization, request)
	if err != nil {
		return "", err
	}

	var resp iqNewOrgResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		return "", err
	}

	return resp.ID, nil
}
