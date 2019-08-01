package nexusiq

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var dummyComponentVersions = map[string][]string{
	dummyComponent.Hash: []string{
		"3.2.1",
		"3.3.2",
		"4.0.6",
		"4.1.31",
		"4.1.34",
		"4.1.36",
		"5.0.16",
		"5.0.18",
		"5.0.28",
	},
}

func componentVersionsTestFunc(t *testing.T, w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost:
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}

		var c Component
		if err = json.Unmarshal(body, &c); err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}

		versions, ok := dummyComponentVersions[c.Hash]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		resp, err := json.Marshal(versions)
		if err != nil {
			w.WriteHeader(http.StatusTeapot)
			return
		}

		fmt.Fprintln(w, string(resp))
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func componentVersionsTestIQ(t *testing.T) (iq IQ, mock *httptest.Server) {
	return newTestIQ(t, componentVersionsTestFunc)
}

func TestComponentVersions(t *testing.T) {
	iq, mock := componentVersionsTestIQ(t)
	defer mock.Close()

	versions, err := ComponentVersions(iq, dummyComponent)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%v\n", versions)

	expected := dummyComponentVersions[dummyComponent.Hash]
	if len(versions) != len(expected) {
		t.Errorf("Got %d results instead of the expected %d", len(versions), len(expected))
	}

	for i, v := range versions {
		if v != expected[i] {
			t.Error("Did not find the expected component version")
		}
	}
}
