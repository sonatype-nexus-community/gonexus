package nexusiq

import (
	"bytes"
	"encoding/json"
	"fmt"
)

const restPolicyViolations = "api/v2/policyViolations"

// ApplicationViolation encapsulates the information about which violations an application has
type ApplicationViolation struct {
	Application      Application       `json:"application"`
	PolicyViolations []PolicyViolation `json:"policyViolations"`
}

// Equals performs a deep comparison between two ApplicationViolation objects
func (a *ApplicationViolation) Equals(b *ApplicationViolation) (_ bool) {
	if a == b {
		return true
	}

	if !a.Application.Equals(&b.Application) {
		return
	}

	if len(a.PolicyViolations) != len(b.PolicyViolations) {
		return
	}

	for i, v := range a.PolicyViolations {
		if !v.Equals(&b.PolicyViolations[i]) {
			return
		}
	}

	return true
}

type violationResponse struct {
	ApplicationViolations []ApplicationViolation `json:"applicationViolations"`
}

// GetAllPolicyViolations returns all policy violations
func GetAllPolicyViolations(iq IQ) ([]ApplicationViolation, error) {
	policyInfos, err := GetPolicies(iq)
	if err != nil {
		return nil, fmt.Errorf("could not get policies: %v", err)
	}

	var endpoint bytes.Buffer
	endpoint.WriteString(restPolicyViolations)
	for _, i := range policyInfos {
		endpoint.WriteString("&p=")
		endpoint.WriteString(i.ID)
	}

	body, _, err := iq.Get(endpoint.String())
	if err != nil {
		return nil, fmt.Errorf("could not get policy violations: %v", err)
	}

	var resp violationResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, fmt.Errorf("could not read policy violations response: %v", err)
	}

	return resp.ApplicationViolations, nil
}

// GetPolicyViolationsByName returns the policy violations by policy name
func GetPolicyViolationsByName(iq IQ, policyName string) ([]ApplicationViolation, error) {
	return nil, nil
}
