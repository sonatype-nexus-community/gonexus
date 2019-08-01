package nexusiq

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path"
	"strings"
	"testing"
)

var dummyRetentionPolicies = map[string]DataRetentionPolicies{
	dummyOrgs[0].ID: DataRetentionPolicies{
		ApplicationReports: ApplicationReports{
			Stages: map[Stage]DataRetentionPolicy{
				StageDevelop: DataRetentionPolicy{
					InheritPolicy: false,
					EnablePurging: true,
					MaxAge:        "3 months",
				},
			},
		},
		SuccessMetrics: DataRetentionPolicy{
			InheritPolicy: false,
			EnablePurging: true,
			MaxAge:        "1 year",
		},
	},
}

func dataRetentionPoliciesTestFunc(t *testing.T, w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet:
		orgName := path.Base(r.URL.Path)

		policies, ok := dummyRetentionPolicies[orgName]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		resp, err := json.Marshal(policies)
		if err != nil {
			w.WriteHeader(http.StatusTeapot)
			return
		}

		fmt.Fprintln(w, string(resp))
	case r.Method == http.MethodPut:
		orgName := path.Base(r.URL.Path)

		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var policies DataRetentionPolicies
		if err = json.Unmarshal(body, &policies); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// TODO: when enablePurging is true, at least one purge criteria, maxAge or maxCount, needs to be specified.

		if existing, ok := dummyRetentionPolicies[orgName]; ok {
			for s, p := range policies.ApplicationReports.Stages {
				existing.ApplicationReports.Stages[s] = p
			}
			existing.SuccessMetrics = policies.SuccessMetrics
			dummyRetentionPolicies[orgName] = existing
		} else {
			dummyRetentionPolicies[orgName] = policies
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func dataRetentionPoliciesTestIQ(t *testing.T) (iq IQ, mock *httptest.Server) {
	return newTestIQ(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasPrefix(r.URL.Path[1:], restOrganization):
			organizationTestFunc(t, w, r)
		default:
			dataRetentionPoliciesTestFunc(t, w, r)
		}
	})
}

func TestGetRetentionPolicies(t *testing.T) {
	iq, mock := dataRetentionPoliciesTestIQ(t)
	defer mock.Close()

	policies, err := GetRetentionPolicies(iq, dummyOrgs[0].Name)
	if err != nil {
		t.Error(err)
	}

	expected := dummyRetentionPolicies[dummyOrgs[0].ID]
	if !policies.Equals(&expected) {
		t.Error("Did not find the expected retention policies")
	}
}

func TestSetRetentionPolicies(t *testing.T) {
	iq, mock := dataRetentionPoliciesTestIQ(t)
	defer mock.Close()

	// expected :=
	var expected = DataRetentionPolicies{
		ApplicationReports: ApplicationReports{
			Stages: map[Stage]DataRetentionPolicy{
				StageDevelop: DataRetentionPolicy{
					InheritPolicy: true,
					EnablePurging: false,
					MaxAge:        "42 months",
				},
			},
		},
		SuccessMetrics: DataRetentionPolicy{
			InheritPolicy: false,
			EnablePurging: true,
			MaxAge:        "1 year",
		},
	}

	var retentionRequest = DataRetentionPolicies{
		ApplicationReports: ApplicationReports{
			Stages: map[Stage]DataRetentionPolicy{
				StageDevelop: expected.ApplicationReports.Stages[StageDevelop],
			},
		},
	}

	err := SetRetentionPolicies(iq, dummyOrgs[0].Name, retentionRequest)
	if err != nil {
		t.Error(err)
	}

	got, err := GetRetentionPolicies(iq, dummyOrgs[0].Name)
	if err != nil {
		t.Error(err)
	}

	if !got.Equals(&expected) {
		t.Error("Did not find the expected retention policies")
	}
}
