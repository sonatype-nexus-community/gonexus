package nexusiq

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

var dummySearchResults = []SearchResult{}

func searchTestFunc(t *testing.T, w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet:
		query := r.URL.Query()
		coordsStr, err := url.QueryUnescape(query["coordinates"][0])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		var coords Coordinates
		if err := json.Unmarshal([]byte(coordsStr), &coords); err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}

		var result searchResponse
		if !coords.Equals(&dummyComponent.ComponentID.Coordinates) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			result.Criteria.ComponentIdentifier.Coordinates = coords
			result.Results = make([]SearchResult, 1)

			result.Results[0].ComponentIdentifier = *(dummyComponent.ComponentID)
			result.Results[0].Hash = dummyComponent.Hash
			result.Results[0].PackageURL = dummyComponent.PackageURL

			resp, err := json.Marshal(result)
			if err != nil {
				t.Fatal(err)
			}

			fmt.Fprintln(w, string(resp))
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func searchTestIQ(t *testing.T) (iq IQ, mock *httptest.Server) {
	return newTestIQ(t, searchTestFunc)
}

func TestSearchComponent(t *testing.T) {
	iq, mock := searchTestIQ(t)
	defer mock.Close()

	query := NewSearchQueryBuilder().Coordinates(dummyComponent.ComponentID.Coordinates)
	components, err := SearchComponents(iq, query)
	if err != nil {
		t.Fatalf("Did not complete search: %v", err)
	}

	t.Logf("%q\n", components)

	if len(components) != 1 {
		t.Errorf("Received %d components instead of %d\n", len(components), 1)
	}

	if !components[0].ComponentIdentifier.Equals(dummyComponent.ComponentID) {
		t.Fatal("Did not receive expected components")
	}
	// TODO: better comparison
}
