package nexusiq

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

var dummyPolicyInfos = []PolicyInfo{}

func policiesTestFunc(t *testing.T, w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet:
		infos, err := json.Marshal(policiesList{dummyPolicyInfos})
		if err != nil {
			t.Fatal(err)
		}

		fmt.Fprintln(w, string(infos))
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func policiesTestIQ(t *testing.T) (iq IQ, mock *httptest.Server) {
	return newTestIQ(t, policiesTestFunc)
}

func TestGetPolicies(t *testing.T) {
	iq, mock := policiesTestIQ(t)
	defer mock.Close()

	infos, err := GetPolicies(iq)
	if err != nil {
		t.Error(err)
	}

	if len(infos) != len(dummyPolicyInfos) {
		t.Errorf("Got %d results instead of the expected %d", len(infos), len(dummyPolicyInfos))
	}

	for i, f := range infos {
		if !f.Equals(&dummyPolicyInfos[i]) {
			t.Fatal("Did not get expected policy info")
		}
	}
}
