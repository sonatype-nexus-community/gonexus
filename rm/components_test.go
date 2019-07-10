package nexusrm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	// "strings"
	"testing"
)

var dummyComponents = map[string][]RepositoryItem{
	"test-repo1": []RepositoryItem{
		RepositoryItem{ID: "component1id", Repository: "test-repo1", Format: "maven2", Group: "org.test", Name: "testComponent1", Version: "1.0.0"},
		RepositoryItem{ID: "component2id", Repository: "test-repo1", Format: "maven2", Group: "org.test", Name: "testComponent2", Version: "2.0.0"},
		RepositoryItem{ID: "component3id", Repository: "test-repo1", Format: "maven2", Group: "org.test", Name: "testComponent3", Version: "3.0.0"},
	},
	"test-repo2": []RepositoryItem{
		RepositoryItem{ID: "component4id", Repository: "test-repo2", Format: "maven2", Group: "org.test", Name: "testComponent4", Version: "4.0.0"},
	},
}

const dummyContinuationToken = "go_on..."

func componentsTestRM(t *testing.T) (rm RM, mock *httptest.Server, err error) {
	return newTestRM(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dump, _ := httputil.DumpRequest(r, true)
		t.Logf("%q\n", dump)

		switch {
		case r.Method == http.MethodGet:
			repo := r.URL.Query().Get("repository")

			lastComponentIdx := len(dummyComponents[repo]) - 1
			var components listComponentsResponse
			if r.URL.Query().Get("continuationToken") == "" {
				components.Items = dummyComponents[repo][:lastComponentIdx]
				components.ContinuationToken = dummyContinuationToken
			} else {
				components.Items = dummyComponents[repo][lastComponentIdx:]
			}

			resp, err := json.Marshal(components)
			if err != nil {
				t.Fatal(err)
			}

			fmt.Fprintln(w, string(resp))
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))
}

func getComponentsTester(t *testing.T, repo string) {
	rm, mock, err := componentsTestRM(t)
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	components, err := GetComponents(rm, repo)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%q\n", components)

	if len(components) != len(dummyComponents[repo]) {
		t.Errorf("Received %d components instead of %d\n", len(components), len(dummyComponents[repo]))
	}

	for i, c := range components {
		if !c.Equals(&dummyComponents[repo][i]) {
			t.Fatal("Did not receive expected components")
		}
	}
}

func TestGetComponentsNoPaging(t *testing.T) {
	getComponentsTester(t, "test-repo2")
}

func TestGetComponentsPaging(t *testing.T) {
	getComponentsTester(t, "test-repo1")
}

func ExampleGetComponents() {
	rm, err := New("http://localhost:8081", "user", "password")
	if err != nil {
		panic(err)
	}

	items, err := GetComponents(rm, "maven-central")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%q\n", items)
}
