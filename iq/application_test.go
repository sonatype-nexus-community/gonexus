package nexusiq

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"strings"
	"testing"
)

var dummyApps = []Application{
	Application{ID: "app1InternalId", PublicID: "app1PubId", Name: "app1Name", OrganizationID: "app1OrgId"},
	Application{ID: "app2InternalId", PublicID: "app2PubId", Name: "app2Name", OrganizationID: "app2OrgId"},
	Application{ID: "app3InternalId", PublicID: "app3PubId", Name: "app3Name", OrganizationID: "app3OrgId"},
	Application{ID: "app4InternalId", PublicID: "app4PubId", Name: "app4Name", OrganizationID: "app4OrgId"},
}

func getAppByPublicID(pubID string) (app Application, i int, ok bool) {
	for i, app = range dummyApps {
		if app.PublicID == pubID {
			return app, i, true
		}
	}
	return
}

func applicationTestIQ(t *testing.T) (iq IQ, mock *httptest.Server, err error) {
	return newTestIQ(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		/*
			getAppByPublicID := func(pubId string) (app Application, i int, ok bool) {
				for i, app = range dummyApps {
					if app.PublicID == pubId {
						return app, i, true
					}
				}
				return
			}
		*/

		dump, _ := httputil.DumpRequest(r, true)
		t.Logf("%q\n", dump)

		switch {
		case r.Method == http.MethodGet && r.URL.String()[1:] == restApplication:
			apps, err := json.Marshal(allAppsResponse{dummyApps})
			if err != nil {
				t.Fatal(err)
			}

			fmt.Fprintln(w, string(apps))
		case r.Method == http.MethodGet && strings.HasPrefix(r.URL.String()[1:], restApplicationByPublic[:len(restApplicationByPublic)-2]):
			pubID := strings.Replace(r.URL.RawQuery, "publicId=", "", -1)
			if app, _, ok := getAppByPublicID(pubID); ok {
				resp, err := json.Marshal(iqAppDetailsResponse{[]Application{app}})
				if err != nil {
					t.Fatal(err)
				}

				fmt.Fprintln(w, string(resp))
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		case r.Method == http.MethodPost:
			defer r.Body.Close()
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
			}
			var appReq iqNewAppRequest
			if err = json.Unmarshal(body, &appReq); err != nil {
				w.WriteHeader(http.StatusBadRequest)
			}
			app := Application{
				ID:              appReq.Name + "InternalId",
				PublicID:        appReq.PublicID,
				Name:            appReq.Name,
				OrganizationID:  appReq.OrganizationID,
				ContactUserName: appReq.ContactUserName,
			}

			// for _, t := range appReq.ApplicationTags {
			// 	app.ApplicationTags = append(app.ApplicationTags, t)
			// }

			dummyApps = append(dummyApps, app)

			resp, err := json.Marshal(app)
			if err != nil {
				w.WriteHeader(http.StatusTeapot)
			}

			fmt.Fprintln(w, string(resp))
		case r.Method == http.MethodDelete:
			pubID := strings.Replace(r.URL.Path[1:], restApplication+"/", "", 1)
			if _, i, ok := getAppByPublicID(pubID); ok {
				copy(dummyApps[i:], dummyApps[i+1:])
				dummyApps[len(dummyApps)-1] = Application{}
				dummyApps = dummyApps[:len(dummyApps)-1]

				w.WriteHeader(http.StatusNoContent)
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))
}

func TestGetAllApplications(t *testing.T) {
	iq, mock, err := applicationTestIQ(t)
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	applications, err := GetAllApplications(iq)
	if err != nil {
		t.Error(err)
	}

	if len(applications) != len(dummyApps) {
		t.Errorf("Expected %d applications but found %d", len(dummyApps), len(applications))
	}

	for i, app := range applications {
		if !app.Equals(&dummyApps[i]) {
			t.Error("Did not get back expected applications")
		}
	}

	t.Logf("%v\n", applications)
}

func TestGetApplicationByPublicID(t *testing.T) {
	iq, mock, err := applicationTestIQ(t)
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	dummyAppsIdx := 2

	app, err := GetApplicationByPublicID(iq, dummyApps[dummyAppsIdx].PublicID)
	if err != nil {
		t.Error(err)
	}

	if !dummyApps[dummyAppsIdx].Equals(app) {
		t.Errorf("Did not retrieve the expected app: %v\n", app)
	}

	t.Log(app)
}

func TestCreateApplication(t *testing.T) {
	iq, mock, err := applicationTestIQ(t)
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	createdApp := Application{PublicID: "createdApp", Name: "createdApp", OrganizationID: "createdAppOrgId"}

	createdApp.ID, err = CreateApplication(iq, createdApp.Name, createdApp.OrganizationID)
	if err != nil {
		t.Fatal(err)
	}

	app, err := GetApplicationByPublicID(iq, createdApp.PublicID)
	if err != nil {
		t.Fatal(err)
	}

	if !createdApp.Equals(app) {
		t.Errorf("Did not retrieve the expected app: %v\n", app)
	}
}

func TestDeleteApplication(t *testing.T) {
	iq, mock, err := applicationTestIQ(t)
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	deleteMeApp := Application{PublicID: "deleteMeApp", Name: "deleteMeApp", OrganizationID: "deleteMeAppOrgId"}

	deleteMeApp.ID, err = CreateApplication(iq, deleteMeApp.Name, deleteMeApp.OrganizationID)
	if err != nil {
		t.Fatal(err)
	}

	if err := DeleteApplication(iq, deleteMeApp.PublicID); err != nil {
		t.Fatal(err)
	}

	if _, err := GetApplicationByPublicID(iq, deleteMeApp.PublicID); err == nil {
		t.Fatal("App was not deleted")
	}
}

func ExampleGetAllApplications() {
	iq, err := New("http://localhost:8070", "username", "password")
	if err != nil {
		panic(err)
	}

	applications, err := GetAllApplications(iq)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", applications)
}
