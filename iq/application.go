package nexusiq

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const (
	restApplication         = "api/v2/applications"
	restApplicationByPublic = "api/v2/applications?publicId=%s"
)

type iqNewAppRequest struct {
	PublicID        string `json:"publicId"`
	Name            string `json:"name"`
	OrganizationID  string `json:"organizationId"`
	ContactUserName string `json:"contactUserName,omitempty"`
	ApplicationTags []struct {
		TagID string `json:"tagId"`
	} `json:"applicationTags,omitempty"`
}

type iqAppDetailsResponse struct {
	Applications []Application `json:"applications"`
}

type allAppsResponse struct {
	Applications []Application `json:"applications"`
}

// Application captures information of an IQ application
type Application struct {
	ID              string `json:"id"`
	PublicID        string `json:"publicId"`
	Name            string `json:"name"`
	OrganizationID  string `json:"organizationId"`
	ContactUserName string `json:"contactUserName"`
	ApplicationTags []struct {
		ID            string `json:"id"`
		TagID         string `json:"tagId"`
		ApplicationID string `json:"applicationId"`
	} `json:"applicationTags,omitempty"`
}

// Equals compares two Application objects
func (a *Application) Equals(b *Application) (_ bool) {
	if a == b {
		return true
	}

	if a.ID != b.ID {
		return
	}

	if a.PublicID != b.PublicID {
		return
	}

	if a.Name != b.Name {
		return
	}

	if a.OrganizationID != b.OrganizationID {
		return
	}

	if a.ContactUserName != b.ContactUserName {
		return
	}

	if len(a.ApplicationTags) != len(b.ApplicationTags) {
		return
	}

	for i, v := range a.ApplicationTags {
		if v.ID != b.ApplicationTags[i].ID {
			return
		}
		if v.TagID != b.ApplicationTags[i].TagID {
			return
		}
		if v.ApplicationID != b.ApplicationTags[i].ApplicationID {
			return
		}
	}

	return true
}

// GetApplicationByPublicID returns details on the named IQ application
func GetApplicationByPublicID(iq IQ, applicationPublicID string) (*Application, error) {
	doError := func(err error) error {
		return fmt.Errorf("application '%s' not found: %v", applicationPublicID, err)
	}
	endpoint := fmt.Sprintf(restApplicationByPublic, applicationPublicID)
	body, _, err := iq.Get(endpoint)
	if err != nil {
		return nil, doError(err)
	}

	var resp iqAppDetailsResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, doError(err)
	}

	if len(resp.Applications) == 0 {
		return nil, errors.New("application not found")
	}

	return &resp.Applications[0], nil
}

// CreateApplication creates an application in IQ with the given name
func CreateApplication(iq IQ, name, organizationID string) (string, error) {
	doError := func(err error) (string, error) {
		return "", fmt.Errorf("application '%s' not created: %v", name, err)
	}
	request, err := json.Marshal(iqNewAppRequest{Name: name, PublicID: name, OrganizationID: organizationID})
	if err != nil {
		return doError(err)
	}

	body, _, err := iq.Post(restApplication, bytes.NewBuffer(request))
	if err != nil {
		return doError(err)
	}

	var resp Application
	if err = json.Unmarshal(body, &resp); err != nil {
		return doError(err)
	}

	return resp.ID, nil
}

// DeleteApplication deletes an application in IQ with the given id
func DeleteApplication(iq IQ, applicationID string) error {
	if resp, err := iq.Del(fmt.Sprintf("%s/%s", restApplication, applicationID)); err != nil && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("application '%s' not deleted: %v", applicationID, err)
	}
	return nil
}

// GetAllApplications returns a slice of all of the applications in an IQ instance
func GetAllApplications(iq IQ) ([]Application, error) {
	body, _, err := iq.Get(restApplication)
	if err != nil {
		return nil, fmt.Errorf("applications not found: %v", err)
	}

	var resp allAppsResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("applications not found: %v", err)
	}

	return resp.Applications, nil
}
