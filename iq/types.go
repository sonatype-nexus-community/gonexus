package nexusiq

// Component identifies a component within IQ
type Component struct {
	Hash                string `json:"hash,omitempty"`
	ComponentIdentifier struct {
		Format      string `json:"format,omitempty"`
		Coordinates struct {
			ArtifactID string `json:"artifactId,omitempty"`
			GroupID    string `json:"groupId,omitempty"`
			Version    string `json:"version,omitempty"`
			Extension  string `json:"extension,omitempty"`
		} `json:"coordinates"`
	} `json:"componentIdentifier,omitempty"`
	Proprietary bool `json:"proprietary,omitempty"`
}

// PolicyViolation is a struct
type PolicyViolation struct {
	PolicyID             string `json:"policyId"`
	PolicyName           string `json:"policyName"`
	ThreatLevel          int    `json:"threatLevel"`
	ConstraintViolations []struct {
		ConstraintID   string `json:"constraintId"`
		ConstraintName string `json:"constraintName"`
		Reasons        []struct {
			Reason string `json:"reason"`
		} `json:"reasons"`
	} `json:"constraintViolations"`
}

// ComponentEvaluationResult is also a struct
type ComponentEvaluationResult struct {
	Component   Component `json:"component"`
	MatchState  string    `json:"matchState"`
	CatalogDate string    `json:"catalogDate"`
	LicenseData struct {
		DeclaredLicenses []struct {
			LicenseID   string `json:"licenseId"`
			LicenseName string `json:"licenseName"`
		} `json:"declaredLicenses"`
		ObservedLicenses []struct {
			LicenseID   string `json:"licenseId"`
			LicenseName string `json:"licenseName"`
		} `json:"observedLicenses"`
		OverriddenLicenses []interface{} `json:"overriddenLicenses"`
		Status             string        `json:"status"`
	} `json:"licenseData"`
	SecurityData struct {
		SecurityIssues []struct {
			Source         string  `json:"source"`
			Reference      string  `json:"reference"`
			Severity       float64 `json:"severity"`
			Status         string  `json:"status"`
			URL            string  `json:"url"`
			ThreatCategory string  `json:"threatCategory"`
		} `json:"securityIssues"`
	} `json:"securityData"`
	PolicyData struct {
		PolicyViolations []PolicyViolation `json:"policyViolations"`
	} `json:"policyData"`
}

// HighestThreatPolicy returns the policy with the highest threat value
func (c *ComponentEvaluationResult) HighestThreatPolicy() *PolicyViolation {
	max, maxVal := -1, -1

	for i, p := range c.PolicyData.PolicyViolations {
		if p.ThreatLevel > maxVal {
			max = i
			maxVal = p.ThreatLevel
		}
	}

	if max < 0 {
		return nil
	}

	return &c.PolicyData.PolicyViolations[max]
}

// Evaluation response thingy
type Evaluation struct {
	SubmittedDate  string                      `json:"submittedDate"`
	EvaluationDate string                      `json:"evaluationDate"`
	ApplicationID  string                      `json:"applicationId"`
	Results        []ComponentEvaluationResult `json:"results"`
	IsError        bool                        `json:"isError"`
	ErrorMessage   interface{}                 `json:"errorMessage"`
}

type iqEvaluationRequestResponse struct {
	ResultID      string `json:"resultId"`
	SubmittedDate string `json:"submittedDate"`
	ApplicationID string `json:"applicationId"`
	ResultsURL    string `json:"resultsUrl"`
}

type iqNewOrgRequest struct {
	Name string `json:"name"`
}

type iqNewOrgResponse struct {
	ID   string   `json:"id"`
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

type iqNewAppRequest struct {
	PublicID        string `json:"publicId"`
	Name            string `json:"name"`
	OrganizationID  string `json:"organizationId"`
	ContactUserName string `json:"contactUserName,omitempty"`
	ApplicationTags []struct {
		TagID string `json:"tagId"`
	} `json:"applicationTags,omitempty"`
}

type iqAppInfoResponse struct {
	Applications []iqAppInfo `json:"applications"`
}

type iqAppInfo struct {
	ID              string `json:"id"`
	PublicID        string `json:"publicId"`
	Name            string `json:"name"`
	OrganizationID  string `json:"organizationId"`
	ContactUserName string `json:"contactUserName"`
	ApplicationTags []struct {
		ID            string `json:"id"`
		TagID         string `json:"tagId"`
		ApplicationID string `json:"applicationId"`
	} `json:"applicationTags"`
}

type iqEvaluationRequest struct {
	Components []Component `json:"components"`
}
