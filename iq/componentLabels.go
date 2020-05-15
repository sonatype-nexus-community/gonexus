package nexusiq

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	restLabelComponent         = "api/v2/components/%s/labels/%s/applications/%s"
	restLabelComponentByOrg    = "api/v2/labels/organization/%s"
	restLabelComponentByOrgDel = "api/v2/labels/organization/%s/%s"
	restLabelComponentByApp    = "api/v2/labels/application/%s"
	restLabelComponentByAppDel = "api/v2/labels/application/%s/%s"
)

// IqComponentLabel describes a component label
type IqComponentLabel struct {
	ID             string `json:"id,omitempty"`
	OwnerID        string `json:"ownerId,omitempty"`
	Label          string `json:"label"`
	LabelLowercase string `json:"labelLowercase,omitempty"`
	Description    string `json:"description,omitempty"`
	Color          string `json:"color"`
}

// ComponentLabelApply adds an existing label to a component for a given application
func ComponentLabelApply(iq IQ, comp Component, appID, label string) error {
	app, err := GetApplicationByPublicID(iq, appID)
	if err != nil {
		return fmt.Errorf("could not retrieve application with ID %s: %v", appID, err)
	}

	endpoint := fmt.Sprintf(restLabelComponent, comp.Hash, url.PathEscape(label), app.ID)
	_, resp, err := iq.Post(endpoint, nil)
	if err != nil {
		if resp == nil || resp.StatusCode != http.StatusNoContent {
			return fmt.Errorf("could not apply label: %v", err)
		}
	}

	return nil
}

// ComponentLabelUnapply removes an existing association between a label and a component
func ComponentLabelUnapply(iq IQ, comp Component, appID, label string) error {
	app, err := GetApplicationByPublicID(iq, appID)
	if err != nil {
		return fmt.Errorf("could not retrieve application with ID %s: %v", appID, err)
	}

	endpoint := fmt.Sprintf(restLabelComponent, comp.Hash, url.PathEscape(label), app.ID)
	resp, err := iq.Del(endpoint)
	if err != nil {
		if resp == nil || resp.StatusCode != http.StatusNoContent {
			return fmt.Errorf("could not unapply label: %v", err)
		}
	}

	return nil
}

func getComponentLabels(iq IQ, endpoint string) ([]IqComponentLabel, error) {
	body, _, err := iq.Get(endpoint)
	if err != nil {
		return nil, err
	}

	var labels []IqComponentLabel
	err = json.Unmarshal(body, &labels)
	if err != nil {
		return nil, err
	}

	return labels, nil
}

// GetComponentLabelsByOrganization retrieves an array of an organization's component label
func GetComponentLabelsByOrganization(iq IQ, organization string) ([]IqComponentLabel, error) {
	endpoint := fmt.Sprintf(restLabelComponentByOrg, organization)
	return getComponentLabels(iq, endpoint)
}

// GetComponentLabelsByAppID retrieves an array of an organization's component label
func GetComponentLabelsByAppID(iq IQ, appID string) ([]IqComponentLabel, error) {
	endpoint := fmt.Sprintf(restLabelComponentByApp, appID)
	return getComponentLabels(iq, endpoint)
}

func createLabel(iq IQ, endpoint, label, description, color string) (IqComponentLabel, error) {
	var labelResponse IqComponentLabel
	request, err := json.Marshal(IqComponentLabel{Label: label, Description: description, Color: color})
	if err != nil {
		return labelResponse, fmt.Errorf("could not marshal label: %v", err)
	}

	body, resp, err := iq.Post(endpoint, bytes.NewBuffer(request))
	if resp.StatusCode != http.StatusOK {
		return labelResponse, fmt.Errorf("did not succeeed in creating label: %v", err)
	}
	defer resp.Body.Close()

	if err := json.Unmarshal(body, &labelResponse); err != nil {
		return labelResponse, fmt.Errorf("could not read json of new label: %v", err)
	}

	return labelResponse, nil
}

// CreateComponentLabelForOrganization creates a label for an organization
func CreateComponentLabelForOrganization(iq IQ, organization, label, description, color string) (IqComponentLabel, error) {
	endpoint := fmt.Sprintf(restLabelComponentByOrg, organization)
	return createLabel(iq, endpoint, label, description, color)
}

// CreateComponentLabelForApplication creates a label for an application
func CreateComponentLabelForApplication(iq IQ, appID, label, description, color string) (IqComponentLabel, error) {
	endpoint := fmt.Sprintf(restLabelComponentByApp, appID)
	return createLabel(iq, endpoint, label, description, color)
}

// DeleteComponentLabelForOrganization deletes a label from an organization
func DeleteComponentLabelForOrganization(iq IQ, organization, label string) error {
	endpoint := fmt.Sprintf(restLabelComponentByOrgDel, organization, label)
	resp, err := iq.Del(endpoint)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("did not succeeed in deleting label: %v", err)
	}
	defer resp.Body.Close()

	return nil
}

// DeleteComponentLabelForApplication deletes a label from an application
func DeleteComponentLabelForApplication(iq IQ, appID, label string) error {
	endpoint := fmt.Sprintf(restLabelComponentByAppDel, appID, label)
	resp, err := iq.Del(endpoint)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("did not succeeed in deleting label: %v", err)
	}
	defer resp.Body.Close()

	return nil
}
