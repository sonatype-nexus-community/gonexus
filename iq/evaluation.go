package nexusiq

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const restEvaluation = "api/v2/evaluation/applications/%s"

type componentCoordinates struct {
	ArtifactID string `json:"artifactId,omitempty"`
	GroupID    string `json:"groupId,omitempty"`
	Version    string `json:"version,omitempty"`
	Extension  string `json:"extension,omitempty"`
}

// ComponentIdentifier identifies the format and coordinates of a component
type ComponentIdentifier struct {
	Format      string               `json:"format,omitempty"`
	Coordinates componentCoordinates `json:"coordinates,omitempty"`
}

// Equals compares two ComponentIdentifier objects
func (a *ComponentIdentifier) Equals(b *ComponentIdentifier) (_ bool) {
	if a == b {
		return true
	}

	if a.Format != b.Format {
		return
	}

	if a.Coordinates.ArtifactID != b.Coordinates.ArtifactID {
		return
	}

	if a.Coordinates.GroupID != b.Coordinates.GroupID {
		return
	}

	if a.Coordinates.Version != b.Coordinates.Version {
		return
	}

	if a.Coordinates.Extension != b.Coordinates.Extension {
		return
	}

	return true
}

// Component encapsulates the details of a component in IQ
type Component struct {
	Hash        string               `json:"hash,omitempty"`
	ComponentID *ComponentIdentifier `json:"componentIdentifier,omitempty"`
	Proprietary bool                 `json:"proprietary,omitempty"`
	PackageURL  string               `json:"packageUrl,omitempty"`
	MatchState  string               `json:"matchState,omitempty"`
	Pathnames   []string             `json:"pathnames,omitempty"`
}

// Equals compares two Component objects
func (a *Component) Equals(b *Component) (_ bool) {
	if a == b {
		return true
	}

	if a.Hash != b.Hash {
		return
	}

	if !a.ComponentID.Equals(b.ComponentID) {
		return
	}

	if a.Proprietary != b.Proprietary {
		return
	}

	if a.PackageURL != b.PackageURL {
		return
	}

	if a.MatchState != b.MatchState {
		return
	}

	if len(a.Pathnames) != len(b.Pathnames) {
		return
	}

	for i, p := range a.Pathnames {
		if p != b.Pathnames[i] {
			return
		}
	}

	return true
}

type versionDetailsRubyGem struct {
	Authors        string      `json:"authors"`
	BuiltAt        time.Time   `json:"built_at"`
	CreatedAt      time.Time   `json:"created_at"`
	Description    string      `json:"description"`
	DownloadsCount int         `json:"downloads_count"`
	Metadata       struct{}    `json:"metadata"`
	Number         string      `json:"number"`
	Summary        string      `json:"summary"`
	Platform       string      `json:"platform"`
	RubyVersion    interface{} `json:"ruby_version"`
	Prerelease     bool        `json:"prerelease"`
	Licenses       interface{} `json:"licenses"`
	Requirements   interface{} `json:"requirements"`
	Sha            string      `json:"sha"`
}

// NewComponentFromString creates a new Component object by parsing
// a string in the expected format; format:group:name:version:ext
func NewComponentFromString(str string) (*Component, error) {
	split := strings.Split(str, ":")

	c := new(Component)
	if len(split) == 1 {
		c.Hash = str
	} else {
		switch split[0] {
		case "maven":
			// c.PackageURL = fmt.Sprintf("pkg:maven/%s/%s@%s?type=%s", split[1], split[2], split[3], split[4])
			c.ComponentID.Format = split[0]
			c.ComponentID.Coordinates.GroupID = split[1]
			c.ComponentID.Coordinates.ArtifactID = split[2]
			c.ComponentID.Coordinates.Version = split[3]
			c.ComponentID.Coordinates.Extension = split[4]
		case "gem":
			c.PackageURL = fmt.Sprintf("pkg:gem/%s@%s?platform=ruby", split[1], split[2])
		case "npm":
			c.PackageURL = fmt.Sprintf("pkg:npm/%s@%s", split[1], split[2])
		case "pypi":
			c.PackageURL = fmt.Sprintf("pkg:pypi/%s@%s?extension=%s", split[1], split[2], split[3])
		case "nuget":
			c.PackageURL = fmt.Sprintf("pkg:nuget/%s@%s", split[1], split[2])
		// case "go":
		default:
			return c, fmt.Errorf("format %s unsupported", split[0])
		}
	}

	return c, nil
}

// PolicyViolation is the policies violated by a component
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

// License identifier an OSS license recognized by Sonatype
type License struct {
	LicenseID   string `json:"licenseId"`
	LicenseName string `json:"licenseName"`
}

// Equals compares two License objects
func (a *License) Equals(b *License) (_ bool) {
	if a == b {
		return true
	}

	if a.LicenseID != b.LicenseID {
		return
	}

	if a.LicenseName != b.LicenseName {
		return
	}

	return true
}

// LicenseData encapsulates the information on the different licenses of a component
type LicenseData struct {
	Status                  string    `json:"status"`
	DeclaredLicenses        []License `json:"declaredLicenses"`
	ObservedLicenses        []License `json:"observedLicenses"`
	OverriddenLicenses      []License `json:"overriddenLicenses"`
	EffectiveLicenseThreats []struct {
		LicenseThreatGroupCategory string `json:"licenseThreatGroupCategory"`
		LicenseThreatGroupLevel    int64  `json:"licenseThreatGroupLevel"`
		LicenseThreatGroupName     string `json:"licenseThreatGroupName"`
	} `json:"effectiveLicenseThreats,omitempty"`
}

// Equals compares two LicenseData objects
func (a *LicenseData) Equals(b *LicenseData) (_ bool) {
	if a == b {
		return true
	}

	if a.Status != b.Status {
		return
	}

	if len(a.DeclaredLicenses) != len(b.DeclaredLicenses) {
		return
	}

	for i, l := range a.DeclaredLicenses {
		if !l.Equals(&b.DeclaredLicenses[i]) {
			return
		}
	}

	if len(a.ObservedLicenses) != len(b.ObservedLicenses) {
		return
	}

	for i, l := range a.ObservedLicenses {
		if !l.Equals(&b.ObservedLicenses[i]) {
			return
		}
	}

	if len(a.OverriddenLicenses) != len(b.OverriddenLicenses) {
		return
	}

	for i, l := range a.OverriddenLicenses {
		if !l.Equals(&b.OverriddenLicenses[i]) {
			return
		}
	}

	if len(a.EffectiveLicenseThreats) != len(b.EffectiveLicenseThreats) {
		return
	}

	for i, l := range a.EffectiveLicenseThreats {
		if l.LicenseThreatGroupCategory != b.EffectiveLicenseThreats[i].LicenseThreatGroupCategory {
			return
		}

		if l.LicenseThreatGroupLevel != b.EffectiveLicenseThreats[i].LicenseThreatGroupLevel {
			return
		}

		if l.LicenseThreatGroupName != b.EffectiveLicenseThreats[i].LicenseThreatGroupName {
			return
		}

	}

	return true
}

// SecurityIssue encapsulates a security issue in the Sonatype database
type SecurityIssue struct {
	Source         string  `json:"source"`
	Reference      string  `json:"reference"`
	Severity       float64 `json:"severity"`
	Status         string  `json:"status"`
	URL            string  `json:"url"`
	ThreatCategory string  `json:"threatCategory"`
}

// Equals compares two SecurityIssue objects
func (a *SecurityIssue) Equals(b *SecurityIssue) (_ bool) {
	if a == b {
		return true
	}

	if a.Source != b.Source {
		return
	}

	if a.Reference != b.Reference {
		return
	}

	if a.Severity != b.Severity {
		return
	}

	if a.Status != b.Status {
		return
	}

	if a.URL != b.URL {
		return
	}

	if a.ThreatCategory != b.ThreatCategory {
		return
	}

	return true
}

// ComponentEvaluationResult holds the results of a component evaluation
type ComponentEvaluationResult struct {
	Component    Component   `json:"component"`
	MatchState   string      `json:"matchState"`
	CatalogDate  string      `json:"catalogDate"`
	LicensesData LicenseData `json:"licenseData"`
	SecurityData struct {
		SecurityIssues []SecurityIssue `json:"securityIssues"`
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
func EvaluateComponents(iq IQ, components []Component, applicationID string) (*Evaluation, error) {
	request, err := json.Marshal(iqEvaluationRequest{Components: components})
	if err != nil {
		return nil, fmt.Errorf("could not build the request: %v", err)
	}

	requestEndpoint := fmt.Sprintf(restEvaluation, applicationID)
	body, _, err := iq.Post(requestEndpoint, bytes.NewBuffer(request))
	if err != nil {
		return nil, fmt.Errorf("components not evaluated: %v", err)
	}

	var results iqEvaluationRequestResponse
	if err = json.Unmarshal(body, &results); err != nil {
		return nil, fmt.Errorf("could not parse evaluation response: %v", err)
	}

	getEvaluationResults := func() (*Evaluation, error) {
		body, resp, e := iq.Get(results.ResultsURL)
		if e != nil {
			if resp.StatusCode != http.StatusNotFound {
				return nil, fmt.Errorf("could not retrieve evaluation results: %v", err)
			}
			return nil, nil
		}

		var ev Evaluation
		if err = json.Unmarshal(body, &ev); err != nil {
			return nil, err
		}

		return &ev, nil
	}

	var eval *Evaluation
	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-ticker.C:
			if eval, err = getEvaluationResults(); eval != nil || err != nil {
				ticker.Stop()
				return eval, err
			}
		case <-time.After(5 * time.Minute):
			ticker.Stop()
			return nil, errors.New("timed out waiting for valid evaluation results")
		}
	}
}
