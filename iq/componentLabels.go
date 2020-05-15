package nexusiq

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	restLabelComponent      = "api/v2/components/%s/labels/%s/applications/%s"
	restLabelComponentByOrg = "api/v2/labels/organization/%s"
	restLabelComponentByApp = "api/v2/labels/application/%s"
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
