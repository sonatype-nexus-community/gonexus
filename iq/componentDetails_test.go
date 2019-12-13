package nexusiq

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
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
				if deets.Component.Hash == c.Hash || reflect.DeepEqual(deets.Component.ComponentID, c.ComponentID) {
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
	return newTestIQ(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		switch {
		// case r.URL.Path[1:] == restOrganization:
		// 	organizationTestFunc(t, w, r)
		case r.URL.Path[1:] == restApplication:
			applicationTestFunc(t, w, r)
		case r.URL.Path[1:] == restReports:
			reportsTestFunc(t, w, r)
		default:
			componentDetailsTestFunc(t, w, r)
		}
	})
}

func TestGetComponent(t *testing.T) {
	iq, mock := componentDetailsTestIQ(t)
	defer mock.Close()

	type args struct {
		iq        IQ
		component Component
	}
	tests := []struct {
		name    string
		args    args
		want    ComponentDetail
		wantErr bool
	}{
		{
			"good component",
			args{iq, dummyComponentDetails[0].Component},
			dummyComponentDetails[0],
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetComponent(tt.args.iq, tt.args.component)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetComponent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetComponent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetComponents(t *testing.T) {
	iq, mock := componentDetailsTestIQ(t)
	defer mock.Close()

	expected := dummyComponentDetails[0]

	details, err := GetComponents(iq, []Component{expected.Component})
	if err != nil {
		t.Error(err)
	}
	t.Logf("%v\n", details)

	if len(details) != 1 {
		t.Fatalf("Received %d results but expected 1\n", len(details))
	}

	if !reflect.DeepEqual(details[0], expected) {
		t.Errorf("Did not receive expected component details")
	}
}

func TestGetComponentsByApplication(t *testing.T) {
	t.Skip("TODO")
	iq, mock := componentDetailsTestIQ(t)
	defer mock.Close()

	type args struct {
		iq          IQ
		appPublicID string
	}
	tests := []struct {
		name    string
		args    args
		want    []ComponentDetail
		wantErr bool
	}{
		{
			"best case",
			args{iq, dummyApps[0].PublicID},
			dummyComponentDetails,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetComponentsByApplication(tt.args.iq, tt.args.appPublicID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetComponentsByApplication() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetComponentsByApplication() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAllComponents(t *testing.T) {
	t.Skip("TODO")
	iq, mock := componentDetailsTestIQ(t)
	defer mock.Close()

	type args struct {
		iq IQ
	}
	tests := []struct {
		name    string
		args    args
		want    []ComponentDetail
		wantErr bool
	}{
		{
			"best case",
			args{iq},
			dummyComponentDetails,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAllComponents(tt.args.iq)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllComponents() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllComponents() = %v, want %v", got, tt.want)
			}
		})
	}
}
