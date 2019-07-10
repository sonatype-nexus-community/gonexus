package nexusiq

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"testing"
)

var dummyOrgs = []Organization{
	Organization{ID: "org1InternalId", Name: "org1Name"},
	Organization{ID: "org2InternalId", Name: "org2Name"},
	Organization{ID: "org3InternalId", Name: "org3Name"},
	Organization{ID: "org4InternalId", Name: "org4Name"},
}

func organizationTestIQ(t *testing.T) (iq IQ, mock *httptest.Server, err error) {
	return newTestIQ(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		/*
			getByPublicID := func(pubId string) (app ApplicationDetails, i int, ok bool) {
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
		case r.Method == http.MethodGet && r.URL.String()[1:] == restOrganization:
			orgs, err := json.Marshal(allOrgsResponse{dummyOrgs})
			if err != nil {
				t.Fatal(err)
			}

			fmt.Fprintln(w, string(orgs))
		case r.Method == http.MethodPost:
			defer r.Body.Close()
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
			}
			var orgReq iqNewOrgRequest
			if err = json.Unmarshal(body, &orgReq); err != nil {
				w.WriteHeader(http.StatusBadRequest)
			}
			org := Organization{
				ID:   orgReq.Name + "InternalId",
				Name: orgReq.Name,
			}

			// for _, t := range orgReq.Tags {
			// 	org.Tags = append(org.Tags, t)
			// }

			dummyOrgs = append(dummyOrgs, org)

			resp, err := json.Marshal(org)
			if err != nil {
				w.WriteHeader(http.StatusTeapot)
			}

			fmt.Fprintln(w, string(resp))
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))
}

func TestGetOranizationByName(t *testing.T) {
	iq, mock, err := organizationTestIQ(t)
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	dummyOrgsIdx := 2

	org, err := GetOrganizationByName(iq, dummyOrgs[dummyOrgsIdx].Name)
	if err != nil {
		t.Error(err)
	}

	if !dummyOrgs[dummyOrgsIdx].Equals(org) {
		t.Errorf("Did not retrieve the expected organization: %v\n", org)
	}

	t.Log(org)
}

func TestCreateOrganization(t *testing.T) {
	iq, mock, err := organizationTestIQ(t)
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	createdOrg := Organization{Name: "createdOrg"}

	createdOrg.ID, err = CreateOrganization(iq, createdOrg.Name)
	if err != nil {
		t.Fatal(err)
	}

	org, err := GetOrganizationByName(iq, createdOrg.Name)
	if err != nil {
		t.Fatal(err)
	}

	if !createdOrg.Equals(org) {
		t.Errorf("Did not retrieve the expected org: %v\n", org)
	}
}

func TestGetAllOrganizations(t *testing.T) {
	iq, mock, err := organizationTestIQ(t)
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	organizations, err := GetAllOrganizations(iq)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", organizations)
}

func ExampleCreateOrganization() {
	iq, err := New("http://localhost:8070", "user", "password")
	if err != nil {
		panic(err)
	}

	orgID, err := CreateOrganization(iq, "DatabaseTeam")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Organization ID: %s\n", orgID)
}
