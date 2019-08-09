package nexusrm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func searchTestFunc(t *testing.T, w http.ResponseWriter, r *http.Request) {
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
}

func searchTestRM(t *testing.T) (rm RM, mock *httptest.Server) {
	return newTestRM(t, searchTestFunc)
}

func TestSearchComponents(t *testing.T) {
	rm, mock := searchTestRM(t)
	defer mock.Close()

	repo := "repo-maven"

	query := NewSearchQueryBuilder().Repository(repo)
	components, err := SearchComponents(rm, query)
	if err != nil {
		t.Fatalf("Did not complete search: %v", err)
	}

	t.Logf("%q\n", components)

	if len(components) != len(dummyComponents[repo]) {
		t.Errorf("Received %d components instead of %d\n", len(components), len(dummyComponents[repo]))
	}

	for i, c := range components {
		if !reflect.DeepEqual(c, dummyComponents[repo][i]) {
			t.Fatal("Did not receive expected components")
		}
	}
}

func TestSearchAssets(t *testing.T) {
	rm, mock := searchTestRM(t)
	defer mock.Close()

	repo := "repo-maven"

	query := NewSearchQueryBuilder().Repository(repo)
	assets, err := SearchAssets(rm, query)
	if err != nil {
		t.Error(err)
	}

	t.Logf("%q\n", assets)

	if len(assets) != len(dummyAssets[repo]) {
		t.Errorf("Received %d assets instead of %d\n", len(assets), len(dummyAssets[repo]))
	}

	for i, c := range assets {
		if !reflect.DeepEqual(c, dummyAssets[repo][i]) {
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
