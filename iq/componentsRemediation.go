package nexusiq

import (
	"bytes"
	"encoding/json"
	"fmt"
)

const (
	restRemediationByApp = "api/v2/components/remediation/application/"
	restRemediationByOrg = "api/v2/components/remediation/organization/"

	remediationTypeNoViolations = "next-no-violations"
	remediationTypeNonFailing   = "next-non-failing"
)

type remediationData struct {
	Component Component `json:"component"`
}

type remediationVersionChange struct {
	Type string          `json:"type"`
	Data remediationData `json:"data"`
}

// Remediation encapsulates the remediation information for a component
type Remediation struct {
	Component      Component                  `json:"component,omitempty"`
	VersionChanges []remediationVersionChange `json:"versionChanges"`
}

type remediationResponse struct {
	Remediation Remediation `json:"remediation"`
}

func createRemediationEndpoint(base, id, stage string) string {
	var buf bytes.Buffer

	buf.WriteString(base)
	buf.WriteString(id)
	if stage != "" {
		buf.WriteString("?stageId=")
		buf.WriteString(stage)
	}

	return buf.String()
}

func getRemediation(iq IQ, component Component, endpoint string) (Remediation, error) {
	request, err := json.Marshal(component)
	if err != nil {
		return Remediation{}, fmt.Errorf("could not build the request: %v", err)
	}

	body, _, err := iq.Post(endpoint, bytes.NewBuffer(request))
	if err != nil {
		return Remediation{}, fmt.Errorf("could not get remediation: %v", err)
	}

	var results remediationResponse
	if err = json.Unmarshal(body, &results); err != nil {
		return Remediation{}, fmt.Errorf("could not parse remediation response: %v", err)
	}

	results.Remediation.Component = component
	return results.Remediation, nil
}

func getRemediationByAppInternalID(iq IQ, component Component, stage, appInternalID string) (Remediation, error) {
	return getRemediation(iq, component, createRemediationEndpoint(restRemediationByApp, appInternalID, stage))
}

// GetRemediationByApp retrieves the remediation information on a component based on an application's policies
func GetRemediationByApp(iq IQ, component Component, stage, applicationID string) (Remediation, error) {
	app, err := GetApplicationByPublicID(iq, applicationID)
	if err != nil {
		return Remediation{}, fmt.Errorf("could not get application: %v", err)
	}

	return getRemediationByAppInternalID(iq, component, stage, app.ID)
}

// GetRemediationByOrg retrieves the remediation information on a component based on an organization's policies
func GetRemediationByOrg(iq IQ, component Component, stage, organizationName string) (Remediation, error) {
	org, err := GetOrganizationByName(iq, organizationName)
	if err != nil {
		return Remediation{}, fmt.Errorf("could not get organization: %v", err)
	}

	endpoint := createRemediationEndpoint(restRemediationByOrg, org.ID, stage)

	return getRemediation(iq, component, endpoint)
}

// GetRemediationsByAppReport retrieves the remediation information on each component of a report
func GetRemediationsByAppReport(iq IQ, applicationID, reportID string) (remediations []Remediation, err error) {
	report, err := GetRawReportByAppReportID(iq, applicationID, reportID)
	if err != nil {
		return nil, fmt.Errorf("could not get report %s for app %s: %v", reportID, applicationID, err)
	}

	app, err := GetApplicationByPublicID(iq, applicationID)
	if err != nil {
		return nil, fmt.Errorf("could not get application: %v", err)
	}

	for _, c := range report.Components {
		if err != nil {
			break
		}
		var remediation Remediation
		remediation, err = getRemediationByAppInternalID(iq, c.Component, "", app.ID)
		remediations = append(remediations, remediation)
	}

	return
}
