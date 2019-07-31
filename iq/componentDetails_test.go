package nexusiq

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var dummyComponentDetails = []ComponentDetail{
	{
		Component: dummyComponent,
	},
}

func componentDetailsTestFunc(t *testing.T, w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost:
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}

		var req detailsRequest
		if err = json.Unmarshal(body, &req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}

		components := make([]ComponentDetail, 0)
		for _, c := range req.Components {
			for _, deets := range dummyComponentDetails {
				if deets.Component.Equals(&c) {
					components = append(components, deets)
				}
			}
		}

		resp, err := json.Marshal(detailsResponse{components})
		if err != nil {
			w.WriteHeader(http.StatusTeapot)
		}

		fmt.Fprintln(w, string(resp))
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func componentDetailsTestIQ(t *testing.T) (iq IQ, mock *httptest.Server) {
	return newTestIQ(t, componentDetailsTestFunc)
}

func TestGetComponent(t *testing.T) {
	iq, mock := componentDetailsTestIQ(t)
	defer mock.Close()

	expected := dummyComponentDetails[0]

	details, err := GetComponent(iq, []Component{expected.Component})
	if err != nil {
		t.Error(err)
	}
	t.Logf("%v\n", details)

	if len(details) != 1 {
		t.Fatalf("Received %d results but expected 1\n", len(details))
	}

	if !details[0].Equals(&expected) {
		t.Errorf("Did not receive expected component details")
	}
}
