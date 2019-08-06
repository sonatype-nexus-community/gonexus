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

var dummyComponents = map[string][]RepositoryItem{
	"repo-maven": []RepositoryItem{
		{
			ID: "component1id", Repository: "repo-maven", Format: "maven2", Group: "org.test", Name: "testComponent1", Version: "1.0.0",
			Assets: []RepositoryItemAsset{dummyAssets["repo-maven"][0]},
		},
		{
			ID: "component2id", Repository: "repo-maven", Format: "maven2", Group: "org.test", Name: "testComponent2", Version: "2.0.0",
			Assets: []RepositoryItemAsset{dummyAssets["repo-maven"][1]},
		},
		{
			ID: "component3id", Repository: "repo-maven", Format: "maven2", Group: "org.test", Name: "testComponent3", Version: "3.0.0",
			Assets: []RepositoryItemAsset{dummyAssets["repo-maven"][2]},
		},
	},
	"repo-npm": []RepositoryItem{
		{
			ID: "component4id", Repository: "repo-npm", Format: "maven2", Group: "org.test", Name: "testComponent4", Version: "4.0.0",
			Assets: []RepositoryItemAsset{dummyAssets["repo-npm"][0]},
		},
	},
}

const dummyNewComponentID = "newComponentID"

func componentsTestFunc(t *testing.T, w http.ResponseWriter, r *http.Request) {
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

	switch {
	case r.Method == http.MethodGet && r.URL.Path[1:] == restComponents:
		query := r.URL.Query()
		repo := query["repository"][0]

		lastComponentIdx := len(dummyComponents[repo]) - 1
		var components listComponentsResponse
		token, ok := query["continuationToken"]
		switch {
		case !ok:
			components.Items = dummyComponents[repo][:lastComponentIdx]
			components.ContinuationToken = dummyContinuationToken
		case token[0] == dummyContinuationToken:
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
			t.Log("crap")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		component := RepositoryItem{ID: dummyNewComponentID, Repository: repo}

		// I have no idea how I want to do this, moving forward
		for k, v := range r.Form {
			switch {
			case k == "repository":
				repo = v[0]
			case k == "maven2.groupId":
				component.Format = "maven2"
				component.Group = v[0]
			case k == "maven2.artifactId":
				component.Format = "maven2"
				component.Name = v[0]
			case k == "maven2.version":
				component.Format = "maven2"
				component.Version = v[0]
			case k == "maven2.packaging":
				component.Format = "maven2"
			case k == "maven2.tag":
				component.Format = "maven2"
			case k == "maven2.generate-pom":
				component.Format = "maven2"
			case strings.HasPrefix(k, "maven2.asset"):
				component.Format = "maven2"
			default:
				t.Logf("Did not recognize form field: %s\n", k)
				w.WriteHeader(http.StatusBadRequest)
			}
			// nContent-Disposition: form-data; name=\"maven2.asset1\" ; filename=\"/tmp/test.jar\"\r\nContent-Type: application/octet-stream\
			// t.Logf("%s = %s\n", k, v)
		}

		dummyComponents[repo] = append(dummyComponents[repo], component)

		w.WriteHeader(http.StatusNoContent)
	case r.Method == http.MethodDelete:
		cID := strings.Replace(r.URL.Path[1:], restComponents+"/", "", 1)
		t.Log(cID)
		if c, i, ok := getComponentByID(cID); ok {
			copy(dummyComponents[c.Repository][i:], dummyComponents[c.Repository][i+1:])
			dummyComponents[c.Repository][len(dummyComponents)-1] = RepositoryItem{}
			dummyComponents[c.Repository] = dummyComponents[c.Repository][:len(dummyComponents)-1]

			w.WriteHeader(http.StatusNoContent)

			resp, err := json.Marshal(c)
			if err != nil {
				t.Fatal(err)
			}

			fmt.Fprintln(w, string(resp))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func componentsTestRM(t *testing.T) (rm RM, mock *httptest.Server) {
	return newTestRM(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasPrefix(r.URL.Path[1:], restRepositories):
			repositoriesTestFunc(t, w, r)
		default:
			componentsTestFunc(t, w, r)
		}
	})
}

func getComponentsTester(t *testing.T, repo string) {
	rm, mock := componentsTestRM(t)
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
		if !reflect.DeepEqual(c, dummyComponents[repo][i]) {
			t.Fatal("Did not receive expected components")
		}
	}
}

func TestGetComponentsNoPaging(t *testing.T) {
	getComponentsTester(t, "repo-npm")
}

func TestGetComponentsPaging(t *testing.T) {
	getComponentsTester(t, "repo-maven")
}

func TestGetComponentByID(t *testing.T) {
	rm, mock := componentsTestRM(t)
	defer mock.Close()

	expectedComponent := dummyComponents["repo-maven"][0]

	component, err := GetComponentByID(rm, expectedComponent.ID)
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("%q\n", component)

	if !reflect.DeepEqual(component, expectedComponent) {
		t.Error("Did not receive expected component")
	}
}

func componentUploader(t *testing.T, expected RepositoryItem, upload UploadComponentWriter) {
	// func componentUploader(t *testing.T, expected RepositoryItem, coordinate string, file io.Reader) {
	rm, mock := componentsTestRM(t)
	defer mock.Close()

	// if err := UploadComponent(rm, expected.Repository, coordinate, file); err != nil {
	if err := UploadComponent(rm, expected.Repository, upload); err != nil {
		t.Error(err)
	}

	expected.ID = dummyNewComponentID

	component, err := GetComponentByID(rm, expected.ID)
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("%q\n", component)

	if !reflect.DeepEqual(component, expected) {
		t.Error("Did not receive expected component")
	}
}

func TestUploadComponentMaven(t *testing.T) {
	coord := "org.test:testComponent3:3.0.0"
	coordSlice := strings.Split(coord, ":")
	expected := RepositoryItem{
		Repository: "repo-maven",
		Format:     "maven2",
		Group:      coordSlice[0],
		Name:       coordSlice[1],
		Version:    coordSlice[2],
		/*
			Assets: []RepositoryItemAssets{RepositoryItemAssets{
				DownloadURL: "",
				Path:        "",
				ID:          "",
				Repository:  "repo-maven",
				Format:      "maven2",
				Checksum:    repositoryItemAssetsChecksum{Sha1: ""},
			},
			},
		*/
	}

	dummyFile := strings.NewReader("foobar")

	upload, err := NewUploadComponentMaven(coord, dummyFile)
	if err != nil {
		t.Fatal(err)
	}

	componentUploader(t, expected, upload)
}

func TestUploadComponentNpm(t *testing.T) {
	t.Skip("TODO")
	expected := RepositoryItem{
		Repository: "repo-npm",
		Format:     "npm",
		/*
			Group:      coordSlice[0],
			Name:       coordSlice[1],
			Version:    coordSlice[2],
				Assets: []RepositoryItemAssets{RepositoryItemAssets{
					DownloadURL: "",
					Path:        "",
					ID:          "",
					Repository:  "repo-maven",
					Format:      "maven2",
					Checksum:    repositoryItemAssetsChecksum{Sha1: ""},
				},
				},
		*/
	}

	dummyFile := strings.NewReader("foobar")

	componentUploader(t, expected, UploadComponentNpm{File: dummyFile})
}

func TestDeleteComponentByID(t *testing.T) {
	rm, mock := componentsTestRM(t)
	defer mock.Close()

	coord := "org.delete:componentDelete:0.0.0"
	coordSlice := strings.Split(coord, ":")
	deleteMe := RepositoryItem{
		ID:         "deleteMe",
		Repository: "repo-maven",
		Format:     "maven2",
		Group:      coordSlice[0],
		Name:       coordSlice[1],
		Version:    coordSlice[2]}

	upload, err := NewUploadComponentMaven(coord, nil)
	if err != nil {
		t.Fatal(err)
	}

	// if err = UploadComponent(rm, deleteMe.Repository, coord, nil); err != nil {
	if err = UploadComponent(rm, deleteMe.Repository, upload); err != nil {
		t.Error(err)
	}

	deleteMe.ID = dummyNewComponentID

	if err = DeleteComponentByID(rm, deleteMe.ID); err != nil {
		t.Fatal(err)
	}

	if _, err := GetComponentByID(rm, deleteMe.ID); err == nil {
		t.Errorf("Component not deleted: %v\n", err)
	}
}

func ExampleGetComponents() {
	rm, err := New("http://localhost:8081", "username", "password")
	if err != nil {
		panic(err)
	}

	items, err := GetComponents(rm, "maven-central")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%q\n", items)
}
