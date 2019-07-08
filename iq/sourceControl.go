package nexusiq

import (
	"encoding/json"
	"fmt"

	"github.com/hokiegeek/gonexus"
)

const restSourceControl = "api/v2/sourceControl/%s"
const restSourceControlDelete = "api/v2/sourceControl/%s/%s"

// SourceControlEntry describes a Source Control entry in IQ
type SourceControlEntry struct {
	ID            string `json:"id,omitempty"`
	ApplicationID string `json:"applicationId"`
	RepositoryURL string `json:"repositoryUrl"`
	Token         string `json:"token"`
}

// GetSourceControlEntry lists of all of the Source Control entries for the given application
func GetSourceControlEntry(iq nexus.Server, applicationID string) (entry SourceControlEntry, err error) {
	appInfo, err := GetApplicationDetailsByPublicID(iq, applicationID)
	if err != nil {
		return
	}

	endpoint := fmt.Sprintf(restSourceControl, appInfo.ID)

	body, _, err := iq.Get(endpoint)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &entry)

	return
}

// CreateSourceControlEntry creates a source control entry in IQ
func CreateSourceControlEntry(iq nexus.Server, applicationID, repositoryURL, token string) error {
	appInfo, err := GetApplicationDetailsByPublicID(iq, applicationID)
	if err != nil {
		return err
	}

	request, err := json.Marshal(SourceControlEntry{"", appInfo.ID, repositoryURL, token})
	if err != nil {
		return err
	}

	endpoint := fmt.Sprintf(restSourceControl, appInfo.ID)

	_, _, err = iq.Post(endpoint, request)
	if err != nil {
		return err
	}

	return nil
}

// UpdateSourceControlEntry updates a source control entry in IQ
func UpdateSourceControlEntry(iq nexus.Server, applicationID, repositoryURL, token string) error {
	appInfo, err := GetApplicationDetailsByPublicID(iq, applicationID)
	if err != nil {
		return err
	}

	request, err := json.Marshal(SourceControlEntry{"", appInfo.ID, repositoryURL, token})
	if err != nil {
		return err
	}

	endpoint := fmt.Sprintf(restSourceControl, appInfo.ID)

	_, _, err = iq.Put(endpoint, request)
	if err != nil {
		return err
	}

	return nil
}

// DeleteSourceControlEntry deletes a source control entry in IQ
func DeleteSourceControlEntry(iq nexus.Server, applicationID, sourceControlID string) error {
	appInfo, err := GetApplicationDetailsByPublicID(iq, applicationID)
	if err != nil {
		return err
	}

	endpoint := fmt.Sprintf(restSourceControlDelete, appInfo.ID, sourceControlID)

	_, err = iq.Del(endpoint)
	if err != nil {
		return err
	}

	return nil
}
