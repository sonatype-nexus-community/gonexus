package nexusiq

import (
	"fmt"
	"net/http"
	"net/url"
)

// api/v2/components/{componentHash}/labels/{labelName}/applications/{applicationId}
const restLabelComponent = "api/v2/components/%s/labels/%s/applications/%s"

// ComponentLabelApply adds an existing label to a component for a given application
func ComponentLabelApply(iq IQ, label string, comp Component, appID string) error {
	app, err := GetApplicationByPublicID(iq, appID)
	if err != nil {
		return fmt.Errorf("could not retrieve application with ID %s: %v", appID, err)
	}

	endpoint := fmt.Sprintf(restLabelComponent, comp.Hash, url.PathEscape(label), app.ID)
	_, resp, err := iq.Post(endpoint, nil)
	if err != nil && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("could not apply label: %v", err)
	}

	return nil
}

// ComponentLabelUnapply removes an existing association between a label and a component
func ComponentLabelUnapply(iq IQ, label string, comp Component, appID string) error {
	app, err := GetApplicationByPublicID(iq, appID)
	if err != nil {
		return fmt.Errorf("could not retrieve application with ID %s: %v", appID, err)
	}

	endpoint := fmt.Sprintf(restLabelComponent, comp.Hash, url.PathEscape(label), app.ID)
	resp, err := iq.Del(endpoint)
	if err != nil && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("could not unapply label: %v", err)
	}

	return nil
}
