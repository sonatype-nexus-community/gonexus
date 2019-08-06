package nexusrm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

var dummyRepos = []Repository{
	{Name: "repo-maven", Format: "maven2", Type: "hosted", URL: "http://localhost:8081/blah/repo-maven"},
	{Name: "repo-nuget", Format: "nuget", Type: "hosted", URL: "http://localhost:8081/blah/repo-nuget"},
	{Name: "repo-pypi", Format: "pypi", Type: "group", URL: "http://localhost:8081/blah/repo-pypi"},
	{Name: "repo-npm", Format: "npm", Type: "proxy", URL: "http://localhost:8081/blah/repo-npm"}, //, Attributes: {Proxy: {RemoteURL: "http://bestest.repo"}}},
}

func repositoriesTestFunc(t *testing.T, w http.ResponseWriter, r *http.Request) {
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
}

func repositoriesTestRM(t *testing.T) (rm RM, mock *httptest.Server) {
	return newTestRM(t, repositoriesTestFunc)
}

func TestGetRepositories(t *testing.T) {
	rm, mock := repositoriesTestRM(t)
	defer mock.Close()

	repos, err := GetRepositories(rm)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%q\n", repos)

	for i, repo := range repos {
		if !reflect.DeepEqual(repo, dummyRepos[i]) {
			t.Error("Did not receive the expected repositories")
		}
	}
}

func TestGetRepositoryByName(t *testing.T) {
	rm, mock := repositoriesTestRM(t)
	defer mock.Close()

	dummyRepoIdx := 0

	repo, err := GetRepositoryByName(rm, dummyRepos[dummyRepoIdx].Name)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%q\n", repo)

	if !reflect.DeepEqual(repo, dummyRepos[dummyRepoIdx]) {
		t.Error("Did not receive the expected repositories")
	}
}
