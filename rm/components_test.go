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

const (
	dummyContinuationToken = "go_on..."
	dummyNewComponentID    = "newComponentID"
)

func componentsTestRM(t *testing.T) (rm RM, mock *httptest.Server, err error) {
	return newTestRM(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		getComponentByID := func(id string) (c RepositoryItem, i int, ok bool) {
			for repo := range dummyComponents {
				for i, c = range dummyComponents[repo] {
					if c.ID == id {
						return c, i, true
					}
				}
			}
			return
		}

		dump, _ := httputil.DumpRequest(r, true)
		t.Logf("%q\n", dump)

		switch {
		case r.Method == http.MethodGet && r.URL.Path[1:] == restComponents:
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
		case r.Method == http.MethodGet && strings.HasPrefix(r.URL.String()[1:], restComponents[:len(restComponents)-2]):
			cID := strings.Replace(r.URL.Path[1:], restComponents+"/", "", 1)
			t.Log(cID)
			if c, _, ok := getComponentByID(cID); ok {
				resp, err := json.Marshal(c)
				if err != nil {
					t.Fatal(err)
				}

				fmt.Fprintln(w, string(resp))
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		case r.Method == http.MethodPost:
			repo := r.URL.Query().Get("repository")
			// TODO check that is valid repository. http 422 if no repo
			// 403 no perms
			// ... might get 100 too

			if err := r.ParseMultipartForm(32 << 20); err != nil {
				w.WriteHeader(http.StatusBadRequest)
			}

			component := RepositoryItem{ID: dummyNewComponentID, Repository: repo}

			for k, v := range r.Form {
				switch k {
				case "maven2.groupId":
					component.Format = "maven2"
					component.Group = v[0]
				case "maven2.artifactId":
					component.Format = "maven2"
					component.Name = v[0]
				case "maven2.version":
					component.Format = "maven2"
					component.Version = v[0]
					// case "maven2.packaging":
					// case "maven2.tag":
					// case "maven2.generate-pom":
				}
				// t.Logf("%s = %s\n", k, v)
			}

			dummyComponents[repo] = append(dummyComponents[repo], component)

			w.WriteHeader(http.StatusNoContent)
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

func TestGetComponentByID(t *testing.T) {
	rm, mock, err := componentsTestRM(t)
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	expectedComponent := dummyComponents["test-repo1"][0]

	component, err := GetComponentByID(rm, expectedComponent.ID)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%q\n", component)

	if !component.Equals(&expectedComponent) {
		t.Error("Did not receive expected component")
	}
}

func componentUploader(t *testing.T, expected RepositoryItem, upload uploadComponent) {
	rm, mock, err := componentsTestRM(t)
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	if err = UploadComponent(rm, expected.Repository, upload); err != nil {
		t.Error(err)
	}

	expected.ID = dummyNewComponentID

	component, err := GetComponentByID(rm, expected.ID)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%q\n", component)

	if !component.Equals(&expected) {
		t.Error("Did not receive expected component")
	}
}

func TestUploadComponentMaven(t *testing.T) {
	expected := RepositoryItem{
		ID:         "newComponent",
		Repository: "test-repo1",
		Format:     "maven2",
		Group:      "org.test",
		Name:       "testComponent3",
		Version:    "3.0.0"}

	upload := UploadComponentMaven{
		GroupID:    expected.Group,
		ArtifactID: expected.Name,
		Version:    expected.Version}
	// Assets                                       [3]UploadAssetMaven

	componentUploader(t, expected, upload)
}

func TestDeleteComponentByID(t *testing.T) {
	t.Skip("TODO")
	/*
		rm, mock, err := componentsTestRM(t)
			if err != nil {
				t.Fatal(err)
			}
			defer mock.Close()

			deleteMe := RepositoryItem{Repository: "no-repo", Format: "maven2", Group: "org.test", Name: "testComponent1", Version: "1.0.0"}

			deleteMe.ID, err = CreateApplication(iq, deleteMe.Name, deleteMe.OrganizationID)
			if err != nil {
				t.Fatal(err)
			}

			if err := DeleteApplication(iq, deleteMe.PublicID); err != nil {
				t.Fatal(err)
			}

			if _, err := GetApplicationByPublicID(iq, deleteMe.PublicID); err == nil {
				t.Fatal("App was not deleted")
			}
	*/
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
