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

type componentRequested struct {
	Hash        string               `json:"hash,omitempty"`
	ComponentID *ComponentIdentifier `json:"componentIdentifier,omitempty"`
	PackageURL  string               `json:"packageUrl,omitempty"`
}

func componentRequestedFromComponent(c Component) componentRequested {
	return componentRequested{
		Hash:        c.Hash,
		ComponentID: c.ComponentID,
		PackageURL:  c.PackageURL,
	}
}

type detailsRequest struct {
	Components []componentRequested `json:"components"`
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
func GetComponent(iq IQ, component Component) (ComponentDetail, error) {
	deets, err := GetComponents(iq, []Component{component})
	if deets == nil || len(deets) == 0 {
		return ComponentDetail{}, err
	}
	return deets[0], err
}

// GetComponents returns information on the named components
func GetComponents(iq IQ, components []Component) ([]ComponentDetail, error) {
	reqComponents := detailsRequest{Components: make([]componentRequested, len(components))}
	for i, c := range components {
		reqComponents.Components[i] = componentRequestedFromComponent(c)
	}

	req, err := json.MarshalIndent(reqComponents, "", " ")
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

// GetComponentsByApplication returns an array with all components along with their
func GetComponentsByApplication(iq IQ, appPublicID string) ([]ComponentDetail, error) {
	componentHashes := make(map[string]struct{})
	components := make([]Component, 0)
	stages := []Stage{StageBuild, StageStageRelease, StageRelease, StageOperate}
	for _, stage := range stages {
		if report, err := GetRawReportByAppID(iq, appPublicID, string(stage)); err == nil {
			for _, c := range report.Components {
				if _, ok := componentHashes[c.Hash]; !ok {
					componentHashes[c.Hash] = struct{}{}
					components = append(components, c.Component)
				}
			}
		}
	}

	return GetComponents(iq, components)
}

// GetAllComponents returns an array with all components along with their
func GetAllComponents(iq IQ) ([]ComponentDetail, error) {
	apps, err := GetAllApplications(iq)
	if err != nil {
		return nil, err
	}

	componentHashes := make(map[string]struct{})
	components := make([]ComponentDetail, 0)

	for _, app := range apps {
		appComponents, err := GetComponentsByApplication(iq, app.PublicID)
		// TODO: catcher
		if err != nil {
			return nil, err
		}

		for _, c := range appComponents {
			if _, ok := componentHashes[c.Component.Hash]; !ok {
				componentHashes[c.Component.Hash] = struct{}{}
				components = append(components, c)
			}
		}
	}

	return components, nil
}
