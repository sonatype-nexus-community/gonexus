package nexusiq

import (
	"encoding/json"
	"fmt"
)

const restOrganization = "api/v2/organizations"

type iqNewOrgRequest struct {
	Name string `json:"name"`
}

type allOrgsResponse struct {
	Organizations []Organization `json:"organizations"`
}

// Organization describes the data in IQ about a given organization
type Organization struct {
	ID   string   `json:"id"`
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

func (a *Organization) Equals(b *Organization) (_ bool) {
	if a == b {
		return true
	}

	if a.ID != b.ID {
		return
	}

	if a.Name != b.Name {
		return
	}

	if len(a.Tags) != len(b.Tags) {
		return
	}

	for i, t := range a.Tags {
		if t != b.Tags[i] {
			return
		}
	}

	return true
}

// GetOrganizationByName returns details on the named IQ organization
func GetOrganizationByName(iq IQ, organizationName string) (*Organization, error) {
	orgs, err := GetAllOrganizations(iq)
	if err != nil {
		return nil, err
	}
	for _, org := range orgs {
		if org.Name == organizationName {
			return &org, nil
		}
	}

	return nil, fmt.Errorf("Did not find organization with name %s", organizationName)
}

// CreateOrganization creates an organization in IQ with the given name
func CreateOrganization(iq IQ, name string) (string, error) {
	request, err := json.Marshal(iqNewOrgRequest{Name: name})
	if err != nil {
		return "", err
	}

	body, _, err := iq.Post(restOrganization, request)
	if err != nil {
		return "", err
	}

	var org Organization
	if err = json.Unmarshal(body, &org); err != nil {
		return "", err
	}

	return org.ID, nil
}

// GetAllOrganizations returns a slice of all of the organizations in an IQ instance
func GetAllOrganizations(iq IQ) ([]Organization, error) {
	body, _, err := iq.Get(restOrganization)
	if err != nil {
		return nil, err
	}

	var resp allOrgsResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return resp.Organizations, nil
}
