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

// Equals performs a deep comparison on two PolicyInfo objects
func (a *PolicyInfo) Equals(b *PolicyInfo) (_ bool) {
	if a == b {
		return true
	}

	if a.ID != b.ID {
		return
	}

	if a.Name != b.Name {
		return
	}

	if a.OwnerID != b.OwnerID {
		return
	}

	if a.OwnerType != b.OwnerType {
		return
	}

	if a.ThreatLevel != b.ThreatLevel {
		return
	}

	if a.PolicyType != b.PolicyType {
		return
	}

	return true
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
