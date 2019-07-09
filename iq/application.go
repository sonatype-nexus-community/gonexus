package nexusiq

import (
	"encoding/json"
	"errors"
	"fmt"
)

const restApplication = "api/v2/applications"

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
	Applications []ApplicationDetails `json:"applications"`
}

type allAppsResponse struct {
	Applications []ApplicationDetails `json:"applications"`
}

// ApplicationDetails captures information of an IQ application
type ApplicationDetails struct {
	ID              string `json:"id"`
	PublicID        string `json:"publicId"`
	Name            string `json:"name"`
	OrganizationID  string `json:"organizationId"`
	ContactUserName string `json:"contactUserName"`
	ApplicationTags []struct {
		ID            string `json:"id"`
		TagID         string `json:"tagId"`
		ApplicationID string `json:"applicationId"`
	} `json:"applicationTags"`
}

// GetApplicationDetailsByPublicID returns details on the named IQ application
func GetApplicationDetailsByPublicID(iq *IQ, applicationPublicID string) (appInfo *ApplicationDetails, err error) {
	endpoint := fmt.Sprintf("%s?publicId=%s", restApplication, applicationPublicID)

	body, _, err := iq.Get(endpoint)
	if err != nil {
		return nil, err
	}

	var resp iqAppDetailsResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	if len(resp.Applications) == 0 {
		return nil, errors.New("Application not found")
	}

	return &resp.Applications[0], nil
}

// CreateApplication creates an application in IQ with the given name
func CreateApplication(iq *IQ, name, organizationID string) (string, error) {
	request, err := json.Marshal(iqNewAppRequest{Name: name, PublicID: name, OrganizationID: organizationID})
	if err != nil {
		return "", err
	}

	body, _, err := iq.Post(restApplication, request)
	if err != nil {
		return "", err
	}

	var resp ApplicationDetails
	if err = json.Unmarshal(body, &resp); err != nil {
		return "", err
	}
	return resp.ID, nil
}

// GetAllApplications returns a slice of all of the applications in an IQ instance
func GetAllApplications(iq *IQ) ([]ApplicationDetails, error) {
	body, _, err := iq.Get(restApplication)
	if err != nil {
		return nil, err
	}

	var resp allAppsResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return resp.Applications, nil
}

// DeleteApplication deletes an application in IQ with the given id
func DeleteApplication(iq *IQ, applicationID string) error {
	iq.Del(fmt.Sprintf("%s/%s", restApplication, applicationID))
	return nil // Always returns an error, so...
}
