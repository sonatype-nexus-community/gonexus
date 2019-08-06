package nexusiq

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

var dummyEntries = []SourceControlEntry{
	SourceControlEntry{ID: "entry1InternalId", ApplicationID: "app1InternalId", RepositoryURL: "entry1URL", Token: "entry1token"},
	SourceControlEntry{ID: "entry2InternalId", ApplicationID: "app2InternalId", RepositoryURL: "entry2URL", Token: "entry2token"},
	SourceControlEntry{ID: "entry3InternalId", ApplicationID: "app3InternalId", RepositoryURL: "entry3URL", Token: "entry3token"},
	// SourceControlEntry{ID: "entry4InternalId", ApplicationID: "app4InternalId", RepositoryURL: "entry4URL", Token: "entry4token"},
}

const newEntryID = "newEntryInternalId"

func sourceControlTestFunc(t *testing.T, w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet && strings.HasPrefix(r.URL.String()[1:], restSourceControl[:len(restSourceControl)-2]):
		appID := strings.Replace(r.URL.Path[1:], restSourceControl[:len(restSourceControl)-2], "", 1)

		var found bool
		for _, entry := range dummyEntries {
			if entry.ApplicationID == appID {
				resp, err := json.Marshal(entry)
				if err != nil {
					t.Error(err)
					http.Error(w, "WTF?", http.StatusTeapot)
				}
				found = true
				fmt.Fprintln(w, string(resp))
			}
		}
		if !found {
			w.WriteHeader(http.StatusNotFound)
		}
	case r.Method == http.MethodPost:
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}

		var entry SourceControlEntry
		if err = json.Unmarshal(body, &entry); err != nil {
			t.Error(err)
			w.WriteHeader(http.StatusTeapot)
		}
		entry.ID = newEntryID
		dummyEntries = append(dummyEntries, entry)
	case r.Method == http.MethodPut:
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}

		var entry SourceControlEntry
		if err = json.Unmarshal(body, &entry); err != nil {
			t.Error(err)
			w.WriteHeader(http.StatusTeapot)
		}
		for i, e := range dummyEntries {
			if e.ID == entry.ID {
				dummyEntries[i].ApplicationID = entry.ApplicationID
				dummyEntries[i].RepositoryURL = entry.RepositoryURL
				dummyEntries[i].Token = entry.Token
			}
		}
		dummyEntries = append(dummyEntries, entry)
	case r.Method == http.MethodDelete:
		splt := strings.Split(r.URL.Path, "/")
		id := splt[len(splt)-1]

		var found bool
		for i, e := range dummyEntries {
			if e.ID == id {
				found = true
				copy(dummyEntries[i:], dummyEntries[i+1:])
				dummyEntries[len(dummyEntries)-1] = SourceControlEntry{}
				dummyEntries = dummyEntries[:len(dummyEntries)-1]

				w.WriteHeader(http.StatusNoContent)
			}
		}

		if !found {
			w.WriteHeader(http.StatusNotFound)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func sourceControlTestIQ(t *testing.T) (iq IQ, mock *httptest.Server) {
	return newTestIQ(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasPrefix(r.URL.Path[1:], restApplication):
			applicationTestFunc(t, w, r)
		default:
			sourceControlTestFunc(t, w, r)
		}
	})
}

func TestGetSourceControlEntryByInternalID(t *testing.T) {
	iq, mock := sourceControlTestIQ(t)
	defer mock.Close()

	dummyEntryIdx := 2

	entry, err := getSourceControlEntryByInternalID(iq, dummyEntries[dummyEntryIdx].ApplicationID)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(entry, dummyEntries[dummyEntryIdx]) {
		t.Errorf("Did not receive expected entry")
	}

	t.Log(entry)
}

func TestGetAllSourceControlEntries(t *testing.T) {
	iq, mock := sourceControlTestIQ(t)
	defer mock.Close()

	entries, err := GetAllSourceControlEntries(iq)
	if err != nil {
		t.Error(err)
	}

	if len(entries) != len(dummyEntries) {
		t.Errorf("Received %d entries instead of the expected %d\n", len(entries), len(dummyEntries))
	}

	for _, entry := range entries {
		var found bool
		for _, dummy := range dummyEntries {
			if !reflect.DeepEqual(dummy, entry) {
				found = true
			}
		}
		if !found {
			t.Fatal("Entries received do not match expected")
		}
	}

	t.Logf("%v\n", entries)
}

func TestGetSourceControlEntry(t *testing.T) {
	iq, mock := sourceControlTestIQ(t)
	defer mock.Close()

	dummyEntryIdx := 0

	entry, err := GetSourceControlEntry(iq, dummyApps[dummyEntryIdx].PublicID)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(entry, dummyEntries[dummyEntryIdx]) {
		t.Errorf("Did not receive expected entry")
	}

	t.Log(entry)
}

func TestCreateSourceControlEntry(t *testing.T) {
	iq, mock := sourceControlTestIQ(t)
	defer mock.Close()

	createdEntry := SourceControlEntry{newEntryID, dummyApps[len(dummyApps)-1].ID, "createdEntryURL", "createEntryToken"}

	err := CreateSourceControlEntry(iq, dummyApps[len(dummyApps)-1].PublicID, createdEntry.RepositoryURL, createdEntry.Token)
	if err != nil {
		t.Error(err)
	}

	entry, err := GetSourceControlEntry(iq, dummyApps[len(dummyApps)-1].PublicID)
	if err != nil {
		t.Error(err)
	}
	t.Logf("Entry: %v\n", entry)

	if !reflect.DeepEqual(entry, createdEntry) {
		t.Errorf("Did not receive expected entry")
	}
}

func TestUpdateSourceControlEntry(t *testing.T) {
	iq, mock := sourceControlTestIQ(t)
	defer mock.Close()

	updatedEntryRepositoryURL := "updatedRepoURL"
	updatedEntryToken := "updatedToken"

	err := UpdateSourceControlEntry(iq, dummyApps[len(dummyApps)-2].PublicID, updatedEntryRepositoryURL, updatedEntryToken)
	if err != nil {
		t.Error(err)
	}

	entry, err := GetSourceControlEntry(iq, dummyApps[len(dummyApps)-2].PublicID)
	if err != nil {
		t.Error(err)
	}
	t.Logf("Entry: %v\n", entry)

	if entry.RepositoryURL != updatedEntryRepositoryURL {
		t.Errorf("Did not receive expected repository URL")
	}

	if entry.Token != updatedEntryToken {
		t.Errorf("Did not receive expected token")
	}
}

func TestDeleteSourceControlEntry(t *testing.T) {
	iq, mock := sourceControlTestIQ(t)
	defer mock.Close()

	app := dummyApps[len(dummyApps)-1]
	deleteMe := SourceControlEntry{newEntryID, app.ID, "deleteMeURL", "deleteMeToken"}

	if err := CreateSourceControlEntry(iq, app.PublicID, deleteMe.RepositoryURL, deleteMe.Token); err != nil {
		t.Error(err)
	}

	if err := DeleteSourceControlEntry(iq, app.PublicID, newEntryID); err != nil {
		t.Error(err)
	}

	if _, err := GetSourceControlEntry(iq, app.PublicID); err == nil {
		t.Error("Unexpectedly found entry which should have been deleted")
	}
}

func TestDeleteSourceControlEntryByApp(t *testing.T) {
	iq, mock := sourceControlTestIQ(t)
	defer mock.Close()

	app := dummyApps[len(dummyApps)-1]
	deleteMe := SourceControlEntry{newEntryID, app.ID, "deleteMeURL", "deleteMeToken"}

	if err := CreateSourceControlEntry(iq, app.PublicID, deleteMe.RepositoryURL, deleteMe.Token); err != nil {
		t.Error(err)
	}

	if err := DeleteSourceControlEntryByApp(iq, app.PublicID); err != nil {
		t.Error(err)
	}

	if _, err := GetSourceControlEntry(iq, app.PublicID); err == nil {
		t.Error("Unexpectedly found entry which should have been deleted")
	}
}
