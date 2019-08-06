package nexusiq

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

var (
	restRemediationByAppPrefix = strings.Split(restRemediationByApp, "%")[0]
	restRemediationByOrgPrefix = strings.Split(restRemediationByOrg, "%")[0]
)

var dummyRemediations = map[string]Remediation{
	"app1InternalId:build": {
		VersionChanges: []remediationVersionChange{
			{
				Type: remediationTypeNoViolations,
			},
			{
				Type: remediationTypeNonFailing,
			},
		},
	},
	"org1InternalId:build": {
		VersionChanges: []remediationVersionChange{
			{
				Type: remediationTypeNoViolations,
			},
			{
				Type: remediationTypeNonFailing,
			},
		},
	},
}

func compRemediationTestFunc(t *testing.T, w http.ResponseWriter, r *http.Request) {
	respond := func(key string, w http.ResponseWriter) {
		if r, ok := dummyRemediations[key]; ok {
			resp, err := json.Marshal(remediationResponse{r})
			if err != nil {
				w.WriteHeader(http.StatusTeapot)
			}

			fmt.Fprintln(w, string(resp))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}

	switch {
	case r.Method == http.MethodPost && strings.HasPrefix(r.URL.Path[1:], restRemediationByAppPrefix):
		key := fmt.Sprintf("%s:%s", strings.ReplaceAll(r.URL.Path[1:], restRemediationByAppPrefix, ""), r.URL.Query()["stageId"][0])
		respond(key, w)
	case r.Method == http.MethodPost && strings.HasPrefix(r.URL.Path[1:], restRemediationByOrgPrefix):
		key := fmt.Sprintf("%s:%s", strings.ReplaceAll(r.URL.Path[1:], restRemediationByOrgPrefix, ""), r.URL.Query()["stageId"][0])
		respond(key, w)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func compRemediationTestIQ(t *testing.T) (iq IQ, mock *httptest.Server) {
	return newTestIQ(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasPrefix(r.URL.Path[1:], restApplication):
			applicationTestFunc(t, w, r)
		case strings.HasPrefix(r.URL.Path[1:], restOrganization):
			organizationTestFunc(t, w, r)
		default:
			compRemediationTestFunc(t, w, r)
		}
	})
}

func TestRemediationByApp(t *testing.T) {
	iq, mock := compRemediationTestIQ(t)
	defer mock.Close()

	id, stage := dummyApps[0].PublicID, "build"

	remediation, err := GetRemediationByApp(iq, dummyComponent, stage, id)
	if err != nil {
		t.Error(err)
	}

	expected := dummyRemediations[dummyApps[0].ID+":"+stage]
	if !reflect.DeepEqual(remediation, expected) {
		t.Error("Did not receive the expected remediation")
	}
}

func TestRemediationByOrg(t *testing.T) {
	iq, mock := compRemediationTestIQ(t)
	defer mock.Close()

	id, stage := dummyOrgs[0].Name, "build"

	remediation, err := GetRemediationByOrg(iq, dummyComponent, stage, id)
	if err != nil {
		t.Error(err)
	}

	expected := dummyRemediations[dummyOrgs[0].ID+":"+stage]
	if !reflect.DeepEqual(remediation, expected) {
		t.Error("Did not receive the expected remediation")
	}
}
