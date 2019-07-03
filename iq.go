package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	// "log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
	// "net/http/httputil"
)

const iqSessionCookieName = "CLM-CSRF-TOKEN"

const iqRestSessionPrivate = "%s/rest/user/session"
const iqRestOrganization = "%s/api/v2/organizations"
const iqRestApplication = "%s/api/v2/applications"
const iqRestOrganizationPrivate = "%s/rest/organization/%s"
const iqRestEvaluation = "%s/api/v2/evaluation/applications/%s"
const iqRestEvaluationResults = "%s/api/v2/evaluation/applications/%s/results/%s"

// IQComponent identifies a component within IQ
type IQComponent struct {
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

// IQPolicyViolation is a struct
type IQPolicyViolation struct {
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

// IQComponentEvaluationResult is also a struct
type IQComponentEvaluationResult struct {
	Component   IQComponent `json:"component"`
	MatchState  string      `json:"matchState"`
	CatalogDate string      `json:"catalogDate"`
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
		PolicyViolations []IQPolicyViolation `json:"policyViolations"`
	} `json:"policyData"`
}

// HighestThreatPolicy returns the policy with the highest threat value
func (c *IQComponentEvaluationResult) HighestThreatPolicy() *IQPolicyViolation {
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

// IQEvaluation response thingy
type IQEvaluation struct {
	SubmittedDate  string                        `json:"submittedDate"`
	EvaluationDate string                        `json:"evaluationDate"`
	ApplicationID  string                        `json:"applicationId"`
	Results        []IQComponentEvaluationResult `json:"results"`
	IsError        bool                          `json:"isError"`
	ErrorMessage   interface{}                   `json:"errorMessage"`
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
	Components []IQComponent `json:"components"`
}

// IQ holds basic and state info on the IQ Server we will connect to
type IQ struct {
	host, username, password string
}

func (iq *IQ) getApplicationInfoByName(applicationName string) (appInfo *iqAppInfo, err error) {
	endpoint := fmt.Sprintf("%s?publicId=%s", iqRestApplication, applicationName)

	body, _, err := iq.get(endpoint)
	if err != nil {
		return nil, err
	}

	var resp iqAppInfoResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	if len(resp.Applications) > 0 {
		return &resp.Applications[0], nil
	}

	return nil, errors.New("Application not found")
}

func (iq *IQ) createTempApplication() (orgID string, appName string, appID string, err error) {
	rand.Seed(time.Now().UnixNano())
	name := strconv.Itoa(rand.Int())

	orgID, err = iq.CreateOrganization(name)
	if err != nil {
		return
	}

	appName = fmt.Sprintf("%s_app", name)

	appID, err = iq.CreateApplication(appName, orgID)
	if err != nil {
		return
	}

	return
}

func (iq *IQ) deleteTempApplication(applicationName string) error {
	appInfo, err := iq.getApplicationInfoByName(applicationName)
	if err != nil {
		return err
	}

	if err := iq.DeleteApplication(appInfo.ID); err != nil {
		return err
	}

	if err := iq.DeleteOrganization(appInfo.OrganizationID); err != nil {
		return err
	}

	return nil
}

func (iq *IQ) getSessionCookies() ([]*http.Cookie, error) {
	_, resp, err := iq.get(iqRestSessionPrivate)
	if err != nil {
		return nil, err
	}

	return resp.Cookies(), nil
}

func (iq *IQ) http(method, endpoint string, payload io.Reader) ([]byte, *http.Response, error) {
	url := fmt.Sprintf(endpoint, iq.host)
	request, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, nil, err
	}

	request.SetBasicAuth(iq.username, iq.password)
	if payload != nil {
		request.Header.Set("Content-Type", "application/json")
	}

	// dump, _ := httputil.DumpRequest(request, true)
	// fmt.Printf("%q\n", dump)

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, nil, err
	}

	if resp.StatusCode == http.StatusOK {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		return body, resp, err
	}

	return nil, resp, errors.New(resp.Status)
}

func (iq *IQ) get(endpoint string) ([]byte, *http.Response, error) {
	return iq.http("GET", endpoint, nil)
}

func (iq *IQ) post(endpoint string, payload []byte) ([]byte, *http.Response, error) {
	return iq.http("POST", endpoint, bytes.NewBuffer(payload))
}

func (iq *IQ) del(endpoint string) error {
	_, _, err := iq.http("DELETE", endpoint, nil)
	return err
}

// CreateOrganization creates an organization in IQ with the given name
func (iq *IQ) CreateOrganization(name string) (string, error) {
	request, err := json.Marshal(iqNewOrgRequest{Name: name})
	if err != nil {
		return "", err
	}

	body, _, err := iq.post(iqRestOrganization, request)
	if err != nil {
		return "", err
	}

	var resp iqNewOrgResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		return "", err
	}

	return resp.ID, nil
}

// CreateApplication creates an application in IQ with the given name
func (iq *IQ) CreateApplication(name, organizationID string) (string, error) {
	request, err := json.Marshal(iqNewAppRequest{Name: name, PublicID: name, OrganizationID: organizationID})
	if err != nil {
		return "", err
	}

	body, _, err := iq.post(iqRestApplication, request)
	if err != nil {
		return "", err
	}

	var resp iqAppInfo
	if err = json.Unmarshal(body, &resp); err != nil {
		return "", err
	}
	return resp.ID, nil
}

// DeleteApplication deletes an application in IQ with the given id
func (iq *IQ) DeleteApplication(applicationID string) error {
	iq.del(fmt.Sprintf("%s/%s", iqRestApplication, applicationID))
	return nil // Always returns an error, so...
}

// DeleteOrganization deletes an organization in IQ with the given id
func (iq *IQ) DeleteOrganization(organizationID string) error {
	// return iq.del(fmt.Sprintf(iqRestOrganizationPrivate, "%s", organizationID))
	url := fmt.Sprintf(iqRestOrganizationPrivate, iq.host, organizationID)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	sessionCookies, err := iq.getSessionCookies()
	if err != nil {
		return err
	}
	for _, cookie := range sessionCookies {
		req.AddCookie(cookie)
		if cookie.Name == iqSessionCookieName {
			req.Header.Add("X-CSRF-TOKEN", cookie.Value)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	// _, resp, err := iq.httpReq(req)
	if resp.StatusCode != http.StatusNoContent {
		return err
	}

	return nil
}

// EvaluateComponents evaluates the list of components
func (iq *IQ) EvaluateComponents(components []IQComponent, applicationID string) (eval *IQEvaluation, err error) {
	request, err := json.Marshal(iqEvaluationRequest{Components: components})
	if err != nil {
		return
	}

	requestEndpoint := fmt.Sprintf(iqRestEvaluation, "%s", applicationID)
	body, _, err := iq.post(requestEndpoint, request)
	if err != nil {
		return
	}

	var results iqEvaluationRequestResponse
	if err = json.Unmarshal(body, &results); err != nil {
		return
	}

	resultsEndpoint := fmt.Sprintf(iqRestEvaluationResults, "%s", results.ApplicationID, results.ResultID)
	getEvaluationResults := func() (*IQEvaluation, error) {
		body, resp, err := iq.get(resultsEndpoint)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode == http.StatusNotFound {
			return nil, nil
		}

		var eval IQEvaluation
		if err = json.Unmarshal(body, &eval); err != nil {
			return nil, err
		}

		return &eval, nil
	}

	ticker := time.NewTicker(10 * time.Second)
	done := make(chan bool, 1)
	go func() {
		for _ = range ticker.C {
			if eval, err = getEvaluationResults(); eval != nil {
				ticker.Stop()
				done <- true
			}
		}
	}()
	<-done

	return
}

// EvaluateComponentsAsFirewall evaluates the list of components using Root Organization only
func (iq *IQ) EvaluateComponentsAsFirewall(components []IQComponent) (eval *IQEvaluation, err error) {
	// Create temp application
	_, appName, appID, err := iq.createTempApplication()
	if err != nil {
		return
	}
	defer iq.deleteTempApplication(appName)

	// Evaluate components
	eval, err = iq.EvaluateComponents(components, appID)
	if err != nil {
		return
	}

	return
}

// NewIQ creates a new IQ object
func NewIQ(host, username, password string) (*IQ, error) {
	iq := new(IQ)
	iq.host = host
	iq.username = username
	iq.password = password
	return iq, nil
}
