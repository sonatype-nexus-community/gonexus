package nexusiq

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var dummyLabels = []string{"dummyLabel1", "dummyLabel2"}
var testAppliedLabels = map[string]string{}

func componentLabelsTestFunc(t *testing.T, w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost:
		splitPath := strings.Split(r.URL.Path[1:], "/")
		componentID, label, appID := splitPath[3], splitPath[5], splitPath[7]
		key := fmt.Sprintf("%s:%s", componentID, appID)
		testAppliedLabels[key] = label

		w.WriteHeader(http.StatusNoContent)
	case r.Method == http.MethodDelete:
		splitPath := strings.Split(r.URL.Path[1:], "/")
		componentID, label, appID := splitPath[3], splitPath[5], splitPath[7]
		key := fmt.Sprintf("%s:%s", componentID, appID)

		appliedLabel, ok := testAppliedLabels[key]
		if !ok || appliedLabel != label {
			w.WriteHeader(http.StatusNotFound)
		} else {
			delete(testAppliedLabels, key)

			w.WriteHeader(http.StatusNoContent)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func componentLabelsTestIQ(t *testing.T) (iq IQ, mock *httptest.Server) {
	return newTestIQ(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasPrefix(r.URL.Path[1:], restApplication):
			applicationTestFunc(t, w, r)
		default:
			componentLabelsTestFunc(t, w, r)
		}
	})
}

func TestComponentLabelApply(t *testing.T) {
	iq, mock := componentLabelsTestIQ(t)
	defer mock.Close()

	label, component, appID := dummyLabels[0], dummyComponent, dummyApps[0].PublicID

	if err := ComponentLabelApply(iq, component, appID, label); err != nil {
		t.Error(err)
	}
}

func TestComponentLabelUnapply(t *testing.T) {
	iq, mock := componentLabelsTestIQ(t)
	defer mock.Close()

	label, component, appID := dummyLabels[0], dummyComponent, dummyApps[0].PublicID

	if err := ComponentLabelApply(iq, component, appID, label); err != nil {
		t.Fatal(err)
	}

	if err := ComponentLabelUnapply(iq, component, appID, label); err != nil {
		t.Error(err)
	}
}
