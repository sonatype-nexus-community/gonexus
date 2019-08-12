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
	dummyApps[0].ID + ":" + StageBuild: {
		Component: dummyComponent,
		VersionChanges: []remediationVersionChange{
			{
				Type: RemediationTypeNoViolations,
			},
			{
				Type: RemediationTypeNonFailing,
			},
		},
	},
	dummyOrgs[0].ID + ":" + StageBuild: {
		Component: dummyComponent,
		VersionChanges: []remediationVersionChange{
			{
				Type: RemediationTypeNoViolations,
			},
			{
				Type: RemediationTypeNonFailing,
			},
		},
	},
	dummyApps[0].ID: {
		Component: dummyComponent,
		VersionChanges: []remediationVersionChange{
			{
				Type: RemediationTypeNoViolations,
			},
			{
				Type: RemediationTypeNonFailing,
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

	var key string
	switch {
	case r.Method == http.MethodPost && strings.HasPrefix(r.URL.Path[1:], restRemediationByAppPrefix):
		key = strings.ReplaceAll(r.URL.Path[1:], restRemediationByAppPrefix, "")
	case r.Method == http.MethodPost && strings.HasPrefix(r.URL.Path[1:], restRemediationByOrgPrefix):
		key = strings.ReplaceAll(r.URL.Path[1:], restRemediationByOrgPrefix, "")
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if stage, ok := r.URL.Query()["stageId"]; ok && len(stage) > 0 {
		key = fmt.Sprintf("%s:%s", key, stage[0])
	}
	respond(key, w)
}

func compRemediationTestIQ(t *testing.T) (iq IQ, mock *httptest.Server) {
	return newTestIQ(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path[1:] == restApplication:
			applicationTestFunc(t, w, r)
		case strings.HasPrefix(r.URL.Path[1:], restOrganization):
			organizationTestFunc(t, w, r)
		case strings.HasPrefix(r.URL.Path[1:], restApplication) && strings.HasSuffix(r.URL.Path, "/raw"):
			reportsTestFunc(t, w, r)
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

func TestRemediationByAppReport(t *testing.T) {
	// t.Skip("TODO")
	iq, mock := compRemediationTestIQ(t)
	defer mock.Close()

	appIdx, reportID := 0, "0"

	got, err := GetRemediationsByAppReport(iq, dummyApps[appIdx].PublicID, reportID)
	if err != nil {
		t.Error(err)
	}

	want := []Remediation{dummyRemediations[dummyApps[appIdx].ID]}
	want[0].Component = Component{
		Hash:       want[0].Component.Hash,
		PackageURL: want[0].Component.PackageURL,
	}
	if !reflect.DeepEqual(want, got) {
		t.Error("Did not receive the expected remediation")
		t.Error("got", got)
		t.Error("want", want)
	}
}
