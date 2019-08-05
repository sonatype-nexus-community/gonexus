package nexusiq

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

var dummyPolicyInfos = []PolicyInfo{
	{
		ID:          "policyInfo1ID",
		Name:        "policyInfo1",
		OwnerID:     "ROOT_ORGANIZATION",
		OwnerType:   "ORGANIZATION",
		ThreatLevel: 42,
		PolicyType:  "license",
	},
	{
		ID:          "policyInfo2ID",
		Name:        "policyInfo2",
		OwnerID:     "ROOT_ORGANIZATION",
		OwnerType:   "ORGANIZATION",
		ThreatLevel: 42,
		PolicyType:  "license",
	},
}

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
		if !reflect.DeepEqual(f, dummyPolicyInfos[i]) {
			t.Fatal("Did not get expected policy info")
		}
	}
}

func TestGetPolicyInfoByName(t *testing.T) {
	iq, mock := policiesTestIQ(t)
	defer mock.Close()

	expected := dummyPolicyInfos[0]

	info, err := GetPolicyInfoByName(iq, expected.Name)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(info, expected) {
		t.Fatal("Did not get expected policy info")
	}
}
