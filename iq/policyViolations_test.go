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

var dummyPolicyViolations = []ApplicationViolation{
	{
		Application: dummyApps[0],
		PolicyViolations: []PolicyViolation{
			{
				PolicyID:    dummyPolicyInfos[0].ID,
				PolicyName:  dummyPolicyInfos[0].Name,
				StageID:     StageBuild,
				ReportURL:   "foobar",
				ThreatLevel: dummyPolicyInfos[0].ThreatLevel,
			},
		},
	},
	{
		Application: dummyApps[1],
		PolicyViolations: []PolicyViolation{
			{
				PolicyID:    dummyPolicyInfos[1].ID,
				PolicyName:  dummyPolicyInfos[1].Name,
				StageID:     StageBuild,
				ReportURL:   "raboof",
				ThreatLevel: dummyPolicyInfos[1].ThreatLevel,
			},
		},
	},
}

func policyViolationsTestFunc(t *testing.T, w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet:
		policies := r.URL.Query()["p"]
		var resp violationResponse
		resp.ApplicationViolations = make([]ApplicationViolation, 0)
		for _, dummy := range dummyPolicyViolations {
			for _, p := range policies {
				for _, v := range dummy.PolicyViolations {
					if v.PolicyID == p {
						resp.ApplicationViolations = append(resp.ApplicationViolations, dummy)
					}
				}
			}
		}
		// TODO: error when a policyID doesn't match any of the violations

		violations, err := json.Marshal(resp)
		if err != nil {
			t.Fatal(err)
		}

		fmt.Fprintln(w, string(violations))
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func policyViolationsTestIQ(t *testing.T) (iq IQ, mock *httptest.Server) {
	return newTestIQ(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasPrefix(r.URL.Path[1:], restPolicies):
			policiesTestFunc(t, w, r)
		default:
			policyViolationsTestFunc(t, w, r)
		}
	})
}

func TestGetAllPolicyViolations(t *testing.T) {
	iq, mock := policyViolationsTestIQ(t)
	defer mock.Close()

	violations, err := GetAllPolicyViolations(iq)
	if err != nil {
		t.Error(err)
	}

	if len(violations) != len(dummyPolicyViolations) {
		t.Errorf("Got %d results instead of the expected %d", len(violations), len(dummyPolicyViolations))
	}

	for i, f := range violations {
		if !reflect.DeepEqual(f, dummyPolicyViolations[i]) {
			t.Fatal("Did not get expected policy violation")
		}
	}
}

func TestGetPolicyViolationsByName(t *testing.T) {
	iq, mock := policyViolationsTestIQ(t)
	defer mock.Close()

	expected := dummyPolicyViolations[0]

	violations, err := GetPolicyViolationsByName(iq, expected.PolicyViolations[0].PolicyName)
	if err != nil {
		t.Error(err)
	}

	if len(violations) != 1 {
		t.Fatalf("Received %d violations instead of the expected %d", len(violations), 1)
	}

	if !reflect.DeepEqual(violations[0], expected) {
		t.Fatal("Did not get expected policy violation")
	}
}
