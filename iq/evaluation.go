package nexusiq

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const restEvaluation = "api/v2/evaluation/applications/%s"

// ComponentIdentifier identifies the format and coordinates of a component
type ComponentIdentifier struct {
	Format      string `json:"format,omitempty"`
	Coordinates struct {
		ArtifactID string `json:"artifactId,omitempty"`
		GroupID    string `json:"groupId,omitempty"`
		Version    string `json:"version,omitempty"`
		Extension  string `json:"extension,omitempty"`
	} `json:"coordinates"`
}

// Component encapsulates the details of a component in IQ
type Component struct {
	Hash        string              `json:"hash,omitempty"`
	ComponentID ComponentIdentifier `json:"componentIdentifier,omitempty"`
	Proprietary bool                `json:"proprietary,omitempty"`
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

type iqEvaluationRequest struct {
	Components []Component `json:"components"`
}

// EvaluateComponents evaluates the list of components
func EvaluateComponents(iq *IQ, components []Component, applicationID string) (eval *Evaluation, err error) {
	request, err := json.Marshal(iqEvaluationRequest{Components: components})
	if err != nil {
		return
	}

	requestEndpoint := fmt.Sprintf(restEvaluation, applicationID)
	body, _, err := iq.Post(requestEndpoint, request)
	if err != nil {
		return
	}

	var results iqEvaluationRequestResponse
	if err = json.Unmarshal(body, &results); err != nil {
		return
	}

	ticker := time.NewTicker(5 * time.Second)
	done := make(chan bool, 1)
	go func() {
		getEvaluationResults := func() (*Evaluation, error) {
			body, resp, err := iq.Get(results.ResultsURL)
			if err != nil {
				return nil, err
			}

			if resp.StatusCode == http.StatusNotFound {
				return nil, nil
			}

			var eval Evaluation
			if err = json.Unmarshal(body, &eval); err != nil {
				return nil, err
			}

			return &eval, nil
		}

		for {
			select {
			case <-ticker.C:
				if eval, err = getEvaluationResults(); eval != nil {
					ticker.Stop()
					done <- true
				}
			case <-time.After(5 * time.Minute):
				ticker.Stop()
				err = errors.New("Timed out waiting for valid results")
				done <- true
			}
		}
	}()
	<-done

	return
}
