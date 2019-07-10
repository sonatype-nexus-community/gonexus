package nexusrm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"testing"
)

var dummyRepos = []Repository{
	{Name: "repo-maven", Format: "maven", Type: "hosted", URL: "http://localhost:8081/blah/repo-maven"},
	{Name: "repo-nuget", Format: "nuget", Type: "hosted", URL: "http://localhost:8081/blah/repo-nuget"},
	{Name: "repo-pypi", Format: "pypi", Type: "group", URL: "http://localhost:8081/blah/repo-pypi"},
	{Name: "repo-npm", Format: "npm", Type: "proxy", URL: "http://localhost:8081/blah/repo-npm"}, //, Attributes: {Proxy: {RemoteURL: "http://bestest.repo"}}},
}

func repositoriesTestRM(t *testing.T) (rm RM, mock *httptest.Server, err error) {
	return newTestRM(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dump, _ := httputil.DumpRequest(r, true)
		t.Logf("%q\n", dump)

		switch {
		case r.Method == http.MethodGet:
			repos, err := json.Marshal(dummyRepos)
			if err != nil {
				t.Fatal(err)
			}

			fmt.Fprintln(w, string(repos))
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))
}

func TestGetRepositories(t *testing.T) {
	rm, mock, err := repositoriesTestRM(t)
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repos, err := GetRepositories(rm)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%q\n", repos)

	for i, repo := range repos {
		if !repo.Equals(&dummyRepos[i]) {
			t.Error("Did not receive the expected repositories")
		}
	}
}
