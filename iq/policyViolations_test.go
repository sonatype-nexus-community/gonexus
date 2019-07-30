package nexusiq

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var dummyPolicyViolations = []ApplicationViolation{}

func policyViolationsTestFunc(t *testing.T, w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet:
		violations, err := json.Marshal(violationResponse{dummyPolicyViolations})
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
		if !f.Equals(&dummyPolicyViolations[i]) {
			t.Fatal("Did not get expected policy violation")
		}
	}
}

func TestGetPolicyViolationsByName(t *testing.T) {
	t.Skip("TODO")
}
