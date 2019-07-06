package nexusiq

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/hokiegeek/gonexus"
)

const restOrganization = "api/v2/organizations"
const restApplication = "api/v2/applications"
const restEvaluation = "api/v2/evaluation/applications/%s"
const restEvaluationResults = "api/v2/evaluation/applications/%s/results/%s"

// IQ holds basic and state info on the IQ Server we will connect to
type IQ struct {
	nexus.DefaultServer
}

// New creates a new IQ instance
func New(host, username, password string) (*IQ, error) {
	iq := new(IQ)
	iq.Host = host
	iq.Username = username
	iq.Password = password
	return iq, nil
}

// GetApplicationDetailsByName returns details on the named IQ application
func GetApplicationDetailsByName(iq *IQ, applicationName string) (appInfo *ApplicationDetails, err error) {
	endpoint := fmt.Sprintf("%s?publicId=%s", restApplication, applicationName)

	body, _, err := iq.Get(endpoint)
	if err != nil {
		return nil, err
	}

	var resp iqAppDetailsResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	if len(resp.Applications) == 0 {
		return nil, errors.New("Application not found")
	}

	return &resp.Applications[0], nil
}

// CreateOrganization creates an organization in IQ with the given name
func CreateOrganization(iq *IQ, name string) (string, error) {
	request, err := json.Marshal(iqNewOrgRequest{Name: name})
	if err != nil {
		return "", err
	}

	body, _, err := iq.Post(restOrganization, request)
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
func CreateApplication(iq *IQ, name, organizationID string) (string, error) {
	request, err := json.Marshal(iqNewAppRequest{Name: name, PublicID: name, OrganizationID: organizationID})
	if err != nil {
		return "", err
	}

	body, _, err := iq.Post(restApplication, request)
	if err != nil {
		return "", err
	}

	var resp ApplicationDetails
	if err = json.Unmarshal(body, &resp); err != nil {
		return "", err
	}
	return resp.ID, nil
}

// DeleteApplication deletes an application in IQ with the given id
func DeleteApplication(iq *IQ, applicationID string) error {
	iq.Del(fmt.Sprintf("%s/%s", restApplication, applicationID))
	return nil // Always returns an error, so...
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

	resultsEndpoint := fmt.Sprintf(restEvaluationResults, results.ApplicationID, results.ResultID)
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
