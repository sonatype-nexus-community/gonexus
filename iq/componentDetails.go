package nexusiq

import (
	"bytes"
	"encoding/json"
	"fmt"
)

const restComponentDetails = "api/v2/components/details"

type detailsResponse struct {
	ComponentDetails []ComponentDetail `json:"componentDetails"`
}

type detailsRequest struct {
	Components []Component `json:"components"`
}

// ComponentDetail lists information about a given component
type ComponentDetail struct {
	Component          Component   `json:"component"`
	MatchState         string      `json:"matchState"`
	CatalogDate        string      `json:"catalogDate"`
	RelativePopularity int64       `json:"relativePopularity,omitempty"`
	LicenseData        LicenseData `json:"licenseData"`
	SecurityData       struct {
		SecurityIssues []SecurityIssue `json:"securityIssues"`
	} `json:"securityData"`
}

// GetComponent returns information on a named component
func GetComponent(iq IQ, components []Component) ([]ComponentDetail, error) {
	req, err := json.Marshal(detailsRequest{components})
	if err != nil {
		return nil, fmt.Errorf("could not generate request: %v", err)
	}

	body, _, err := iq.Post(restComponentDetails, bytes.NewBuffer(req))
	if err != nil {
		return nil, fmt.Errorf("could not find component details: %v", err)
	}

	var resp detailsResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("could not process component details: %v", err)
	}

	return resp.ComponentDetails, nil
}
