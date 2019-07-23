package nexusrm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"strings"
	"testing"
)

func searchTestRM(t *testing.T) (rm RM, mock *httptest.Server, err error) {
	return newTestRM(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dump, _ := httputil.DumpRequest(r, true)
		t.Logf("%q\n", dump)

		switch {
		case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path[1:], restSearchAssets):
			// TODO: http.StatusNotAcceptable if only have cont token

			query := r.URL.Query()
			repo := query["repository"][0]

			var assets searchAssetsResponse
			if _, ok := dummyAssets[repo]; ok {
				lastComponentIdx := len(dummyAssets[repo]) - 1
				token, ok := query["continuationToken"]
				switch {
				case !ok:
					assets.Items = dummyAssets[repo][:lastComponentIdx]
					assets.ContinuationToken = dummyContinuationToken
				case token[0] == dummyContinuationToken:
					assets.Items = dummyAssets[repo][lastComponentIdx:]
				}
			}

			resp, err := json.Marshal(assets)
			if err != nil {
				t.Fatal(err)
			}

			fmt.Fprintln(w, string(resp))
		case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path[1:], restSearchComponents):
			// TODO: http.StatusNotAcceptable if only have cont token

			query := r.URL.Query()
			repo := query["repository"][0]

			var components searchComponentsResponse
			if _, ok := dummyComponents[repo]; ok {
				lastComponentIdx := len(dummyComponents[repo]) - 1
				token, ok := query["continuationToken"]
				switch {
				case !ok:
					components.Items = dummyComponents[repo][:lastComponentIdx]
					components.ContinuationToken = dummyContinuationToken
				case token[0] == dummyContinuationToken:
					components.Items = dummyComponents[repo][lastComponentIdx:]
				}
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

func TestSearchComponents(t *testing.T) {
	rm, mock, err := searchTestRM(t)
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := "repo-maven"

	query := NewSearchQueryBuilder().Repository(repo)
	components, err := SearchComponents(rm, query)
	if err != nil {
		panic(err)
	}

	t.Logf("%q\n", components)

	if len(components) != len(dummyComponents[repo]) {
		t.Errorf("Received %d components instead of %d\n", len(components), len(dummyComponents[repo]))
	}

	for i, c := range components {
		if !c.Equals(&dummyComponents[repo][i]) {
			t.Fatal("Did not receive expected components")
		}
	}
}

func TestSearchAssets(t *testing.T) {
	t.Skip("Waiting to implement assets endpoint tests")
	rm, mock, err := searchTestRM(t)
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := "repo-maven"

	query := NewSearchQueryBuilder().Repository(repo)
	assets, err := SearchAssets(rm, query)
	if err != nil {
		panic(err)
	}

	t.Logf("%q\n", assets)

	if len(assets) != len(dummyAssets[repo]) {
		t.Errorf("Received %d assets instead of %d\n", len(assets), len(dummyAssets[repo]))
	}

	for i, c := range assets {
		if !c.Equals(&dummyAssets[repo][i]) {
			t.Fatal("Did not receive expected components")
		}
	}
}

func ExampleSearchComponents() {
	rm, err := New("http://localhost:8081", "username", "password")
	if err != nil {
		panic(err)
	}

	query := NewSearchQueryBuilder().Repository("maven-releases")
	components, err := SearchComponents(rm, query)
	if err != nil {
		panic(err)
	}

	for _, c := range components {
		fmt.Println(c.Name)
	}
}
