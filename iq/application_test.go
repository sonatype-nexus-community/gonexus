package nexusiq

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

var dummyApps = []Application{
	{ID: "app1InternalId", PublicID: "app1PubId", Name: "app1Name", OrganizationID: "org1InternalId"},
	{ID: "app2InternalId", PublicID: "app2PubId", Name: "app2Name", OrganizationID: "org2InternalId"},
	{ID: "app3InternalId", PublicID: "app3PubId", Name: "app3Name", OrganizationID: "org3InternalId"},
	{ID: "app4InternalId", PublicID: "app4PubId", Name: "app4Name", OrganizationID: "org4InternalId"},
}

func getAppByPublicID(pubID string) (app Application, i int, ok bool) {
	for i, app = range dummyApps {
		if app.PublicID == pubID {
			return app, i, true
		}
	}
	return
}

func applicationTestFunc(t *testing.T, w http.ResponseWriter, r *http.Request) {
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
}

func applicationTestIQ(t *testing.T) (iq IQ, mock *httptest.Server) {
	return newTestIQ(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path[1:] == restOrganization:
			organizationTestFunc(t, w, r)
		default:
			applicationTestFunc(t, w, r)
		}
	})
}

func TestGetAllApplications(t *testing.T) {
	iq, mock := applicationTestIQ(t)
	defer mock.Close()

	applications, err := GetAllApplications(iq)
	if err != nil {
		t.Error(err)
	}

	if len(applications) != len(dummyApps) {
		t.Errorf("Expected %d applications but found %d", len(dummyApps), len(applications))
	}

	for i, app := range applications {
		if !reflect.DeepEqual(app, dummyApps[i]) {
			t.Error("Did not get back expected applications")
		}
	}

	t.Logf("%v\n", applications)
}

func TestGetApplicationByPublicID(t *testing.T) {
	iq, mock := applicationTestIQ(t)
	defer mock.Close()

	dummyAppsIdx := 2

	got, err := GetApplicationByPublicID(iq, dummyApps[dummyAppsIdx].PublicID)
	if err != nil {
		t.Error(err)
	}

	want := dummyApps[dummyAppsIdx]
	if !reflect.DeepEqual(*got, want) {
		t.Error("Did not retrieve the expected app")
		t.Error("got", got)
		t.Error("want", want)
	}

	t.Log(got)
}

func TestCreateApplication(t *testing.T) {
	iq, mock := applicationTestIQ(t)
	defer mock.Close()

	type args struct {
		iq             IQ
		name           string
		id             string
		organizationID string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			"best case",
			args{iq: iq, name: "createdApp", id: "createdApp", organizationID: "createdAppOrgId"},
			"createdAppInternalId",
			false,
		},
		{
			"missing name",
			args{iq: iq, id: "createdApp", organizationID: "createdAppOrgId"},
			"",
			true,
		},
		{
			"missing id",
			args{iq: iq, name: "createdApp", organizationID: "createdAppOrgId"},
			"",
			true,
		},
		{
			"missing org",
			args{iq: iq, name: "createdApp", id: "createdApp"},
			"",
			true,
		},
		{
			"missing options",
			args{iq: iq},
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateApplication(tt.args.iq, tt.args.name, tt.args.id, tt.args.organizationID)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateApplication() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CreateApplication() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeleteApplication(t *testing.T) {
	iq, mock := applicationTestIQ(t)
	defer mock.Close()

	deleteMeApp := Application{PublicID: "deleteMeApp", Name: "deleteMeApp", OrganizationID: "deleteMeAppOrgId"}

	var err error
	deleteMeApp.ID, err = CreateApplication(iq, deleteMeApp.Name, deleteMeApp.PublicID, deleteMeApp.OrganizationID)
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

func TestGetApplicationsByOrganization(t *testing.T) {
	iq, mock := applicationTestIQ(t)
	defer mock.Close()

	type args struct {
		iq               IQ
		organizationName string
	}
	tests := []struct {
		name    string
		args    args
		want    []Application
		wantErr bool
	}{
		{
			"good org",
			args{iq, dummyOrgs[0].Name},
			[]Application{dummyApps[0]},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetApplicationsByOrganization(tt.args.iq, tt.args.organizationName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetApplicationsByOrganization() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetApplicationsByOrganization() = %v, want %v", got, tt.want)
			}
		})
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

func ExampleCreateApplication() {
	iq, err := New("http://localhost:8070", "user", "password")
	if err != nil {
		panic(err)
	}

	appID, err := CreateApplication(iq, "name", "id", "organization")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Application ID: %s\n", appID)
}
