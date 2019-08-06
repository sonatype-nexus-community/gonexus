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
	endpoint.WriteString("?")
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
func GetPolicyViolationsByName(iq IQ, policyNames ...string) ([]ApplicationViolation, error) {
	policies, err := GetPolicies(iq)
	if err != nil {
		return nil, fmt.Errorf("did not find policy: %v", err)
	}

	var endpoint bytes.Buffer
	endpoint.WriteString(restPolicyViolations)
	endpoint.WriteString("?")

	for _, p := range policyNames {
		for _, policy := range policies {
			if p == policy.Name {
				endpoint.WriteString("&p=")
				endpoint.WriteString(policy.ID)
			}
		}
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
