package nexusiq

import (
	"encoding/json"
	"fmt"
	// "io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	// "strings"
	"testing"
)

var dummyReportInfos = []ReportInfo{
	{
		ApplicationID:           "app1ID",
		EmbeddableReportHTMLURL: "WhoEmbedsThis?",
		EvaluationDate:          "evalDate",
		ReportDataURL:           "dataURL",
		ReportHTMLURL:           "htmlURL",
		ReportPdfURL:            "pdfURL",
		Stage:                   StageBuild,
	},
}

func reportsTestIQ(t *testing.T) (iq IQ, mock *httptest.Server, err error) {
	return newTestIQ(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dump, _ := httputil.DumpRequest(r, true)
		t.Logf("%q\n", dump)

		switch {
		case r.Method == http.MethodGet && r.URL.Path[1:] == restReports:
			infos, err := json.Marshal(dummyReportInfos)
			if err != nil {
				t.Fatal(err)
			}

			fmt.Fprintln(w, string(infos))
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))
}

func TestGetAllReportInfos(t *testing.T) {
	iq, mock, err := reportsTestIQ(t)
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	infos, err := GetAllReportInfos(iq)
	if err != nil {
		t.Error(err)
	}

	if len(infos) != len(dummyReportInfos) {
		t.Errorf("Got %d results instead of the expected %d", len(infos), len(dummyReportInfos))
	}

	for i, f := range infos {
		if !f.Equals(&dummyReportInfos[i]) {
			t.Fatal("Did not get expected report info")
		}
	}
}

func TestGetReportInfosByAppID(t *testing.T) {
	t.Skip("TODO")
}

func TestGetRawReportByAppID(t *testing.T) {
	t.Skip("TODO")
}

func TestGetPolicyReportByAppID(t *testing.T) {
	t.Skip("TODO")
}

func TestGetReportByAppID(t *testing.T) {
	t.Skip("TODO")
}
