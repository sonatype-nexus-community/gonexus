package nexusiq

import (
	"encoding/json"
	"fmt"
)

const restPolicies = "api/v2/policies"

// PolicyInfo encapsulates the identifying information of an individual IQ policy
type PolicyInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	OwnerID     string `json:"ownerId"`
	OwnerType   string `json:"ownerType"`
	ThreatLevel int    `json:"threatLevel"`
	PolicyType  string `json:"policyType"`
}

type policiesList struct {
	Policies []PolicyInfo `json:"policies"`
}

// GetPolicies returns a list of all of the policies in IQ
func GetPolicies(iq IQ) ([]PolicyInfo, error) {
	body, _, err := iq.Get(restPolicies)
	if err != nil {
		return nil, fmt.Errorf("could not get list of policies: %v", err)
	}

	var resp policiesList
	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("could not read endpoint response: %v", err)
	}

	return resp.Policies, nil
}

// GetPolicyInfoByName returns an information object for the named policy
func GetPolicyInfoByName(iq IQ, policyName string) (PolicyInfo, error) {
	policies, err := GetPolicies(iq)
	if err != nil {
		return PolicyInfo{}, fmt.Errorf("did not find policy with name %s: %v", policyName, err)
	}

	for _, p := range policies {
		if p.Name == policyName {
			return p, nil
		}
	}

	return PolicyInfo{}, fmt.Errorf("did not find policy with name %s", policyName)
}
