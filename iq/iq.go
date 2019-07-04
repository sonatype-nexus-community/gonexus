package nexusiq

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"time"
	// "net/http/httputil"

	"github.com/hokiegeek/gonexus"
)

const iqRestOrganizationPrivate = "rest/organization/%s"
const iqRestSessionPrivate = "rest/user/session"
const iqRestOrganization = "api/v2/organizations"
const iqRestApplication = "api/v2/applications"
const iqRestEvaluation = "api/v2/evaluation/applications/%s"
const iqRestEvaluationResults = "api/v2/evaluation/applications/%s/results/%s"

// IQ holds basic and state info on the IQ Server we will connect to
type IQ struct {
	nexus.Server
}

func (iq *IQ) getApplicationInfoByName(applicationName string) (appInfo *iqAppInfo, err error) {
	endpoint := fmt.Sprintf("%s?publicId=%s", iqRestApplication, applicationName)

	body, _, err := iq.Get(endpoint)
	if err != nil {
		return nil, err
	}

	var resp iqAppInfoResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	if len(resp.Applications) == 0 {
		return nil, errors.New("Application not found")
	}

	return &resp.Applications[0], nil
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

func (iq *IQ) newPrivateRequest(method, endpoint string, payload io.Reader) (*http.Request, error) {
	req, err := iq.NewRequest(method, endpoint, payload)
	if err != nil {
		return nil, err
	}

	_, resp, err := iq.Get(iqRestSessionPrivate)
	if err != nil {
		return nil, err
	}

	for _, cookie := range resp.Cookies() {
		req.AddCookie(cookie)
		if cookie.Name == "CLM-CSRF-TOKEN" {
			req.Header.Add("X-CSRF-TOKEN", cookie.Value)
		}
	}

	return req, nil
}

// CreateOrganization creates an organization in IQ with the given name
func (iq *IQ) CreateOrganization(name string) (string, error) {
	request, err := json.Marshal(iqNewOrgRequest{Name: name})
	if err != nil {
		return "", err
	}

	body, _, err := iq.Post(iqRestOrganization, request)
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

	body, _, err := iq.Post(iqRestApplication, request)
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
	iq.Del(fmt.Sprintf("%s/%s", iqRestApplication, applicationID))
	return nil // Always returns an error, so...
}

// DeleteOrganization deletes an organization in IQ with the given id
func (iq *IQ) DeleteOrganization(organizationID string) error {
	url := fmt.Sprintf(iqRestOrganizationPrivate, organizationID)

	req, err := iq.newPrivateRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	_, resp, err := iq.Do(req)
	if err != nil || resp.StatusCode != http.StatusNoContent {
		return err
	}

	return nil
}

// EvaluateComponents evaluates the list of components
func (iq *IQ) EvaluateComponents(components []Component, applicationID string) (eval *Evaluation, err error) {
	request, err := json.Marshal(iqEvaluationRequest{Components: components})
	if err != nil {
		return
	}

	requestEndpoint := fmt.Sprintf(iqRestEvaluation, applicationID)
	body, _, err := iq.Post(requestEndpoint, request)
	if err != nil {
		return
	}

	var results iqEvaluationRequestResponse
	if err = json.Unmarshal(body, &results); err != nil {
		return
	}

	resultsEndpoint := fmt.Sprintf(iqRestEvaluationResults, results.ApplicationID, results.ResultID)
	ticker := time.NewTicker(5 * time.Second)
	done := make(chan bool, 1)
	go func() {
		getEvaluationResults := func() (*Evaluation, error) {
			body, resp, err := iq.Get(resultsEndpoint)
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
				if eval, err = getEvaluationResults(); eval != nil || err != nil {
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

// EvaluateComponentsAsFirewall evaluates the list of components using Root Organization only
func (iq *IQ) EvaluateComponentsAsFirewall(components []Component) (eval *Evaluation, err error) {
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

// New creates a new IQ instance
func New(host, username, password string) (*IQ, error) {
	iq := new(IQ)
	iq.Host = host
	iq.Username = username
	iq.Password = password
	return iq, nil
}
