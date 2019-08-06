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

var dummyAssets = map[string][]RepositoryItemAsset{
	"repo-maven": []RepositoryItemAsset{
		{
			ID:          "asset1id",
			DownloadURL: "http://localhost:8081/repository/repo-maven/org/test/testComponent1/1.0.0/testComponent1-1.0.0.jar",
			Path:        "org/test/testComponent1/1.0.0/testComponent1-1.0.0.jar",
			Repository:  "repo-maven",
			Format:      "maven2",
			Checksum:    repositoryItemAssetsChecksum{Sha1: "asset1sha1", Md5: "asset1md5"},
		},
		{
			ID:          "asset2id",
			DownloadURL: "http://localhost:8081/repository/repo-maven/org/test/testComponent2/2.0.0/testComponent2-2.0.0.jar",
			Path:        "org/test/testComponent2/2.0.0/testComponent2-2.0.0.jar",
			Repository:  "repo-maven",
			Format:      "maven2",
			Checksum:    repositoryItemAssetsChecksum{Sha1: "asset2sha1", Md5: "asset2md5"},
		},
		{
			ID:          "asset3id",
			DownloadURL: "http://localhost:8081/repository/repo-maven/org/test/testComponent3/3.0.0/testComponent3-3.0.0.jar",
			Path:        "org/test/testComponent3/3.0.0/testComponent3-3.0.0.jar",
			Repository:  "repo-maven",
			Format:      "maven2",
			Checksum:    repositoryItemAssetsChecksum{Sha1: "asset3sha1", Md5: "asset3md5"},
		},
	},
	"repo-npm": []RepositoryItemAsset{
		{
			ID:          "asset4id",
			DownloadURL: "http://localhost:8081/repository/repo-npm/testComponent4/-/4.0.0/testComponent4-4.0.0.tgz",
			Path:        "testComponent4/-/testComponent4-4.0.0.tgz",
			Repository:  "repo-npm",
			Format:      "npm",
			Checksum:    repositoryItemAssetsChecksum{Sha1: "asset4sha1", Md5: "asset4md5"},
		},
	},
}

const dummyNewAssetID = "newAssetID"

func assetsTestFunc(t *testing.T, w http.ResponseWriter, r *http.Request) {
	getAssetByID := func(id string) (a RepositoryItemAsset, i int, ok bool) {
		for repo := range dummyAssets {
			for i, a = range dummyAssets[repo] {
				if a.ID == id {
					return a, i, true
				}
			}
		}
		return
	}

	switch {
	case r.Method == http.MethodGet && strings.HasPrefix(r.URL.RequestURI()[1:], restListAssetsByRepo[:len(restListAssetsByRepo)-2]):
		query := r.URL.Query()
		repo := query["repository"][0]

		lastAssetIdx := len(dummyAssets[repo]) - 1
		var assets listAssetsResponse
		token, ok := query["continuationToken"]
		switch {
		case !ok:
			assets.Items = dummyAssets[repo][:lastAssetIdx]
			assets.ContinuationToken = dummyContinuationToken
		case token[0] == dummyContinuationToken:
			assets.Items = dummyAssets[repo][lastAssetIdx:]
		}

		resp, err := json.Marshal(assets)
		if err != nil {
			t.Fatal(err)
		}

		fmt.Fprintln(w, string(resp))
	case r.Method == http.MethodGet:
		aID := strings.Replace(r.URL.Path[1:], restAssets+"/", "", 1)
		t.Log(aID)
		if c, _, ok := getAssetByID(aID); ok {
			resp, err := json.Marshal(c)
			if err != nil {
				t.Fatal(err)
			}

			fmt.Fprintln(w, string(resp))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	case r.Method == http.MethodDelete:
		aID := strings.Replace(r.URL.Path[1:], restAssets+"/", "", 1)
		t.Log(aID)
		if a, i, ok := getAssetByID(aID); ok {
			copy(dummyAssets[a.Repository][i:], dummyAssets[a.Repository][i+1:])
			dummyAssets[a.Repository][len(dummyAssets)-1] = RepositoryItemAsset{}
			dummyAssets[a.Repository] = dummyAssets[a.Repository][:len(dummyAssets)-1]

			w.WriteHeader(http.StatusNoContent)

			resp, err := json.Marshal(a)
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

func assetsTestRM(t *testing.T) (rm RM, mock *httptest.Server) {
	return newTestRM(t, assetsTestFunc)
}

func getAssetsTester(t *testing.T, repo string) {
	rm, mock := assetsTestRM(t)
	defer mock.Close()

	assets, err := GetAssets(rm, repo)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%q\n", assets)

	if len(assets) != len(dummyAssets[repo]) {
		t.Errorf("Received %d assets instead of %d\n", len(assets), len(dummyAssets[repo]))
	}

	for i, c := range assets {
		if !reflect.DeepEqual(c, dummyAssets[repo][i]) {
			t.Fatal("Did not receive expected assets")
		}
	}
}

func TestGetAssetsNoPaging(t *testing.T) {
	getAssetsTester(t, "repo-npm")
}

func TestGetAssetsPaging(t *testing.T) {
	getAssetsTester(t, "repo-maven")
}

func TestGetAssetByID(t *testing.T) {
	rm, mock := assetsTestRM(t)
	defer mock.Close()

	expectedAsset := dummyAssets["repo-maven"][0]

	asset, err := GetAssetByID(rm, expectedAsset.ID)
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("%q\n", asset)

	if !reflect.DeepEqual(asset, expectedAsset) {
		t.Error("Did not receive expected asset")
	}
}

func TestDeleteAssetByID(t *testing.T) {
	rm, mock := assetsTestRM(t)
	defer mock.Close()

	deleteMe := RepositoryItemAsset{
		ID:          "testDeleteAsset",
		DownloadURL: "http://localhost:8081/repository/repo-maven/org/test/testDeleteAsset/1.2.3/testDeleteAsset-1.2.3.jar",
		Path:        "org/test/testDeleteAsset/1.2.3/testDeleteAsset-1.2.3.jar",
		Repository:  "repo-maven",
		Format:      "maven2",
		Checksum:    repositoryItemAssetsChecksum{Sha1: "assetDeletesha1", Md5: "assetDeletemd5"},
	}

	dummyAssets[deleteMe.Repository] = append(dummyAssets[deleteMe.Repository], deleteMe)

	if _, err := GetAssetByID(rm, deleteMe.ID); err != nil {
		t.Errorf("Error getting component: %v\n", err)
	}

	if err := DeleteAssetByID(rm, deleteMe.ID); err != nil {
		t.Fatal(err)
	}

	if _, err := GetAssetByID(rm, deleteMe.ID); err == nil {
		t.Errorf("Asset not deleted: %v\n", err)
	}
}
