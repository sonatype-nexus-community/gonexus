package nexusiq

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

const reportDataURLFormat = "api/v2/applications/%s/reports/%s/raw"

var dummyReportInfos = []ReportInfo{
	{
		ApplicationID:           dummyApps[0].ID,
		EmbeddableReportHTMLURL: "WhoEmbedsThis?",
		EvaluationDateStr:       "evalDate",
		ReportDataURL:           fmt.Sprintf(reportDataURLFormat, dummyApps[0].PublicID, "0"),
		ReportHTMLURL:           "htmlURL",
		ReportPdfURL:            "pdfURL",
		Stage:                   StageBuild,
	},
	{
		ApplicationID:           dummyApps[1].ID,
		EmbeddableReportHTMLURL: "WhoEmbedsThis?",
		EvaluationDateStr:       "evalDate",
		ReportDataURL:           fmt.Sprintf(reportDataURLFormat, dummyApps[0].PublicID, "1"),
		ReportHTMLURL:           "htmlURL",
		ReportPdfURL:            "pdfURL",
		Stage:                   StageBuild,
	},
}

var dummyRawReports = map[string]ReportRaw{
	dummyReportInfos[0].ReportDataURL: ReportRaw{
		Components: []rawReportComponent{
			{
				Component: Component{
					Hash: "045c37a03be19f3e0db8",
					ComponentID: &ComponentIdentifier{
						Format: "maven",
						Coordinates: Coordinates{
							ArtifactID: "jackson-databind",
							GroupID:    "com.fasterxml.jackson.core",
							Version:    "2.6.1",
							Extension:  "jar",
						},
					},
				},
			},
		},
		MatchSummary: rawReportMatchSummary{KnownComponentCount: 11, TotalComponentCount: 111},
		ReportInfo:   dummyReportInfos[0],
	},
	dummyReportInfos[1].ReportDataURL: ReportRaw{},
}

var dummyPolicyReports = map[string]ReportPolicy{
	strings.Replace(dummyReportInfos[0].ReportDataURL, "/raw", "/policy", 1): ReportPolicy{
		Application: dummyApps[0],
		Components: []PolicyReportComponent{
			{
				Component: Component{
					Hash: "045c37a03be19f3e0db8",
					ComponentID: &ComponentIdentifier{
						Format: "maven",
						Coordinates: Coordinates{
							ArtifactID: "jackson-databind",
							GroupID:    "com.fasterxml.jackson.core",
							Version:    "2.6.1",
							Extension:  "jar",
						},
					},
				},
			},
		},
		Counts: policyReportCounts{
			ExactlyMatchedComponentCount:      1,
			GrandfatheredPolicyViolationCount: 1,
			PartiallyMatchedComponentCount:    1,
			TotalComponentCount:               1,
		},
		ReportTime:  54,
		ReportTitle: "foobar",
		ReportInfo:  dummyReportInfos[0],
	},
	strings.Replace(dummyReportInfos[1].ReportDataURL, "/raw", "/policy", 1): ReportPolicy{},
}

func reportsTestFunc(t *testing.T, w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet && r.URL.Path[1:] == restReports:
		infos, err := json.Marshal(dummyReportInfos)
		if err != nil {
			t.Fatal(err)
		}

		fmt.Fprintln(w, string(infos))
	case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path[1:], restReports):
		appID := strings.Replace(r.URL.Path[1:], restReports+"/", "", 1)

		var found bool
		for _, r := range dummyReportInfos {
			if r.ApplicationID == appID {
				found = true
				resp, err := json.Marshal([]ReportInfo{r})
				if err != nil {
					t.Fatal(err)
				}

				fmt.Fprintln(w, string(resp))
			}
		}

		if !found {
			w.WriteHeader(http.StatusNotFound)
		}
	case r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/raw"):
		if raw, ok := dummyRawReports[r.URL.Path[1:]]; ok {
			resp, err := json.Marshal(raw)
			if err != nil {
				t.Fatal(err)
			}

			fmt.Fprintln(w, string(resp))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	case r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/policy"):
		if policy, ok := dummyPolicyReports[r.URL.Path[1:]]; ok {
			resp, err := json.Marshal(policy)
			if err != nil {
				t.Fatal(err)
			}

			fmt.Fprintln(w, string(resp))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	default:
		t.Log("wtf", r.URL.Path)
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func reportsTestIQ(t *testing.T) (iq IQ, mock *httptest.Server) {
	return newTestIQ(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path[1:] == restOrganization:
			organizationTestFunc(t, w, r)
		case r.URL.Path[1:] == restApplication:
			applicationTestFunc(t, w, r)
		default:
			reportsTestFunc(t, w, r)
		}
	})
}

func TestGetAllReportInfos(t *testing.T) {
	iq, mock := reportsTestIQ(t)
	defer mock.Close()

	infos, err := GetAllReportInfos(iq)
	if err != nil {
		t.Error(err)
	}

	if len(infos) != len(dummyReportInfos) {
		t.Errorf("Got %d results instead of the expected %d", len(infos), len(dummyReportInfos))
	}

	for i, f := range infos {
		if !reflect.DeepEqual(f, dummyReportInfos[i]) {
			t.Fatal("Did not get expected report info")
		}
	}
}

func TestGetReportInfosByAppID(t *testing.T) {
	iq, mock := reportsTestIQ(t)
	defer mock.Close()

	testIdx := 0

	infos, err := GetReportInfosByAppID(iq, dummyApps[testIdx].PublicID)
	if err != nil {
		t.Error(err)
	}

	if len(infos) != 1 {
		t.Errorf("Got %d results instead of the expected 1", len(infos))
	}

	if !reflect.DeepEqual(infos[0], dummyReportInfos[testIdx]) {
		t.Fatal("Did not get expected report info")
	}
}

func Test_getRawReportByAppReportID(t *testing.T) {
	iq, mock := reportsTestIQ(t)
	defer mock.Close()

	testIdx := 0

	report, err := getRawReportByAppReportID(iq, dummyApps[testIdx].PublicID, fmt.Sprintf("%d", testIdx))
	if err != nil {
		t.Fatal(err)
	}

	dummy := dummyRawReports[dummyReportInfos[testIdx].ReportDataURL]
	if !reflect.DeepEqual(report, dummy) {
		t.Error("Did not get expected raw report")
	}
}

func TestGetRawReportByAppID(t *testing.T) {
	iq, mock := reportsTestIQ(t)
	defer mock.Close()

	testIdx := 0

	report, err := GetRawReportByAppID(iq, dummyApps[testIdx].PublicID, dummyReportInfos[testIdx].Stage)
	if err != nil {
		t.Fatal(err)
	}

	dummy := dummyRawReports[dummyReportInfos[testIdx].ReportDataURL]
	if !reflect.DeepEqual(report, dummy) {
		t.Error("Did not get expected raw report")
	}
}

func TestGetPolicyReportByAppID(t *testing.T) {
	iq, mock := reportsTestIQ(t)
	defer mock.Close()

	testIdx := 0

	report, err := GetPolicyReportByAppID(iq, dummyApps[testIdx].PublicID, dummyReportInfos[testIdx].Stage)
	if err != nil {
		t.Fatal(err)
	}

	dummy := dummyPolicyReports[strings.Replace(dummyReportInfos[testIdx].ReportDataURL, "/raw", "/policy", 1)]
	if !reflect.DeepEqual(report, dummy) {
		t.Error("Did not get expected policy report")
	}
}

func TestGetReportByAppID(t *testing.T) {
	iq, mock := reportsTestIQ(t)
	defer mock.Close()

	testIdx := 0

	report, err := GetReportByAppID(iq, dummyApps[testIdx].PublicID, dummyReportInfos[testIdx].Stage)
	if err != nil {
		t.Fatal(err)
	}

	dummyRaw := dummyRawReports[dummyReportInfos[testIdx].ReportDataURL]
	if !reflect.DeepEqual(report.Raw, dummyRaw) {
		t.Error("Did not get expected raw report")
	}

	dummyPolicy := dummyPolicyReports[strings.Replace(dummyReportInfos[testIdx].ReportDataURL, "/raw", "/policy", 1)]
	if !reflect.DeepEqual(report.Policy, dummyPolicy) {
		t.Error("Did not get expected policy report")
	}
}

func TestGetReportInfosByOrganization(t *testing.T) {
	iq, mock := reportsTestIQ(t)
	defer mock.Close()

	type args struct {
		iq               IQ
		organizationName string
	}
	tests := []struct {
		name      string
		args      args
		wantInfos []ReportInfo
		wantErr   bool
	}{
		{
			"test1",
			args{iq, dummyOrgs[0].Name},
			[]ReportInfo{dummyReportInfos[0]},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotInfos, err := GetReportInfosByOrganization(tt.args.iq, tt.args.organizationName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetReportInfosByOrganization() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotInfos, tt.wantInfos) {
				t.Errorf("GetReportInfosByOrganization() = %v, want %v", gotInfos, tt.wantInfos)
			}
		})
	}
}
