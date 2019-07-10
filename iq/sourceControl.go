package nexusiq

import (
	"encoding/json"
	"fmt"
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

func getSourceControlEntryByInternalID(iq IQ, applicationID string) (entry SourceControlEntry, err error) {
	endpoint := fmt.Sprintf(restSourceControl, applicationID)

	body, _, err := iq.Get(endpoint)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &entry)

	return
}

// GetSourceControlEntry lists of all of the Source Control entries for the given application
func GetSourceControlEntry(iq IQ, applicationID string) (entry SourceControlEntry, err error) {
	appInfo, err := GetApplicationDetailsByPublicID(iq, applicationID)
	if err != nil {
		return
	}

	return getSourceControlEntryByInternalID(iq, appInfo.ID)
}

// GetAllSourceControlEntries lists of all of the Source Control entries in the IQ instance
func GetAllSourceControlEntries(iq IQ) (entries []SourceControlEntry, err error) {
	apps, err := GetAllApplications(iq)
	if err != nil {
		return
	}

	for _, app := range apps {
		if entry, err := getSourceControlEntryByInternalID(iq, app.ID); err == nil {
			entries = append(entries, entry)
		}
	}

	return
}

// CreateSourceControlEntry creates a source control entry in IQ
func CreateSourceControlEntry(iq IQ, applicationID, repositoryURL, token string) error {
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
func UpdateSourceControlEntry(iq IQ, applicationID, repositoryURL, token string) error {
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

func deleteSourceControlEntry(iq IQ, appInternalID, sourceControlID string) error {
	endpoint := fmt.Sprintf(restSourceControlDelete, appInternalID, sourceControlID)

	_, err := iq.Del(endpoint)
	if err != nil {
		return err
	}

	return nil
}

// DeleteSourceControlEntry deletes a source control entry in IQ
func DeleteSourceControlEntry(iq IQ, applicationID, sourceControlID string) error {
	appInfo, err := GetApplicationDetailsByPublicID(iq, applicationID)
	if err != nil {
		return err
	}

	return deleteSourceControlEntry(iq, appInfo.ID, sourceControlID)
}

// DeleteSourceControlEntryByApp deletes a source control entry in IQ for the given application
func DeleteSourceControlEntryByApp(iq IQ, applicationID string) error {
	appInfo, err := GetApplicationDetailsByPublicID(iq, applicationID)
	if err != nil {
		return err
	}

	entry, err := getSourceControlEntryByInternalID(iq, appInfo.ID)
	if err != nil {
		return err
	}

	return deleteSourceControlEntry(iq, appInfo.ID, entry.ID)
}

// DeleteSourceControlEntryByEntry deletes a source control entry in IQ for the given entry ID
/*
func DeleteSourceControlEntryByEntry(iq IQ, sourceControlID string) error {
	entry, err := getSourceControlEntryByInternalID(iq, appInfo.ID)
	if err != nil {
		return err
	}

	return deleteSourceControlEntry(iq, entry.ApplicationID, entry.ID)
}
*/
