package nexusiq

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

var dummyMetrics = []struct {
	firstTimePeriod, lastTimePeriod time.Time
	apps, orgs                      map[string]struct{}
	metrics                         []Metrics
}{
	{
		time.Now().Add(-(14 * (24 * time.Hour))),
		time.Now(),
		map[string]struct{}{dummyApps[0].ID: struct{}{}},
		nil,
		[]Metrics{
			{
				ApplicationID:       dummyApps[0].ID,
				ApplicationPublicID: dummyApps[0].PublicID,
				ApplicationName:     dummyApps[0].Name,
				OrganizationID:      dummyOrgs[0].ID,
				OrganizationName:    dummyOrgs[0].Name,
				Aggregations: []aggregation{
					{
						TimePeriodStart:    "2018-08-01",
						EvaluationCount:    4,
						MttrModerateThreat: 1885594213,
						MttrCriticalThreat: 74576,
						DiscoveredCounts: map[violationCountType]violationCounts{
							licenseCount:  violationCounts{0, 0, 0, 4},
							otherCount:    violationCounts{0, 3, 0, 0},
							qualityCount:  violationCounts{0, 0, 0, 0},
							securityCount: violationCounts{0, 0, 1, 0},
						},
						FixedCounts: map[violationCountType]violationCounts{
							licenseCount:  violationCounts{0, 0, 0, 4},
							otherCount:    violationCounts{0, 3, 0, 0},
							qualityCount:  violationCounts{0, 0, 0, 0},
							securityCount: violationCounts{0, 0, 0, 0},
						},
						OpenCountsAtTimePeriodEnd: map[violationCountType]violationCounts{
							licenseCount:  violationCounts{0, 0, 0, 8},
							otherCount:    violationCounts{0, 4, 0, 0},
							qualityCount:  violationCounts{21, 5, 0, 0},
							securityCount: violationCounts{0, 1, 4, 3},
						},
					},
				},
			},
		},
	},
}

func reportMetricsTestFunc(t *testing.T, w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost:
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		var req metricRequest
		if err = json.Unmarshal(body, &req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}

		var reqFirst, reqLast time.Time

		switch req.TimePeriod {
		case metricsMonthly:
			reqFirst, _ = time.Parse("2006-01", req.FirstTimePeriod)
			if req.LastTimePeriod != "" {
				reqLast, _ = time.Parse("2006-01", req.LastTimePeriod)
			}
		case metricsWeekly:
			// TODO
		}

		metrics := make([]Metrics, 0)
		for _, m := range dummyMetrics {
			var found bool
			t.Log("shit1", reqFirst, m.firstTimePeriod)
			if m.firstTimePeriod.Before(reqFirst) {
				continue
			}

			if reqLast.After(m.lastTimePeriod) {
				t.Log("shit2", reqLast, m.lastTimePeriod)
				continue
			}

			if len(req.ApplicationIDS) > 0 {
				for _, reqApp := range req.ApplicationIDS {
					if _, ok := m.apps[reqApp]; ok {
						found = true
					}
				}

				if !found {
					continue
				}
			}

			if len(req.OrganizationIDS) > 0 {
				for _, reqOrg := range req.OrganizationIDS {
					if _, ok := m.orgs[reqOrg]; ok {
						found = true
					}
				}

				if !found {
					continue
				}
			}

			metrics = append(metrics, m.metrics...)
		}

		resp, err := json.Marshal(metrics)
		if err != nil {
			w.WriteHeader(http.StatusTeapot)
		}

		fmt.Fprintln(w, string(resp))
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func reportMetricsTestIQ(t *testing.T) (IQ, *httptest.Server) {
	return newTestIQ(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path[1:] == restOrganization:
			organizationTestFunc(t, w, r)
		case r.URL.Path[1:] == restApplication:
			applicationTestFunc(t, w, r)
		default:
			reportMetricsTestFunc(t, w, r)
		}
	})
}

func TestMetricsRequestBuilder(t *testing.T) {
	var tests = []struct {
		input MetricsRequestBuilder
		want  metricRequest
	}{
		{
			input: MetricsRequestBuilder{
				timePeriod:      metricsMonthly,
				firstTimePeriod: time.Now(),
			},
			want: metricRequest{
				TimePeriod:      metricsMonthly,
				FirstTimePeriod: time.Now().Format("2006-01"),
			},
		},
		{
			input: MetricsRequestBuilder{
				timePeriod:      metricsMonthly,
				firstTimePeriod: time.Now(),
				apps:            []string{dummyApps[0].PublicID},
			},
			want: metricRequest{
				TimePeriod:      metricsMonthly,
				FirstTimePeriod: time.Now().Format("2006-01"),
				ApplicationIDS:  []string{dummyApps[0].ID},
			},
		},
		{
			MetricsRequestBuilder{
				metricsWeekly,
				time.Date(2009, 11, 1, 1, 1, 1, 1, time.UTC),
				time.Date(2019, 1, 1, 1, 1, 1, 1, time.UTC),
				[]string{dummyApps[0].PublicID},
				[]string{dummyOrgs[0].Name},
			},
			metricRequest{
				TimePeriod:      metricsWeekly,
				FirstTimePeriod: "2009-W44",
				LastTimePeriod:  "2019-W01",
				ApplicationIDS:  []string{dummyApps[0].ID},
				OrganizationIDS: []string{dummyOrgs[0].ID},
			},
		},
	}

	iq, mock := reportMetricsTestIQ(t)
	defer mock.Close()

	for _, test := range tests {
		got, err := test.input.build(iq)
		if err != nil {
			t.Errorf("Unexpected error building metrics request: %v", err)
			t.Error("input", test.input)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Error("Did not get expected request")
			t.Error("got", got)
			t.Error("want", test.want)
		}
	}
}

func TestGenerateMetrics(t *testing.T) {
	iq, mock := reportMetricsTestIQ(t)
	defer mock.Close()

	tests := []struct {
		input *MetricsRequestBuilder
		want  []Metrics
	}{
		{
			input: NewMetricsRequestBuilder().Monthly().StartingOn(time.Now().Add((120 * (24 * time.Hour)))).WithApplication(dummyApps[0].PublicID),
			want:  []Metrics{},
		},
		{
			input: NewMetricsRequestBuilder().Monthly().StartingOn(time.Now().Add(-(30 * (24 * time.Hour)))).WithApplication(dummyApps[0].PublicID),
			want:  dummyMetrics[0].metrics,
		},
	}

	for _, test := range tests {
		got, err := GenerateMetrics(iq, test.input)
		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Error("Did not get expected metrics")
			t.Error("got", got)
			t.Error("want", test.want)
		}
	}
}

func ExampleGenerateMetrics() {
	iq, err := New("http://localhost:8070", "admin", "admin123")
	if err != nil {
		panic(err)
	}

	reqLastYear := NewMetricsRequestBuilder().Monthly().StartingOn(time.Now().Add(-(24 * time.Hour) * 365)).WithApplication("WebGoat")

	metrics, err := GenerateMetrics(iq, reqLastYear)
	if err != nil {
		panic(err)
	}

	buf, err := json.MarshalIndent(metrics, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(buf))
}
