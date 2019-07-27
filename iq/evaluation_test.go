package nexusiq

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	activeEval            iqEvaluationRequestResponse
	activeEvalAttempts    = 0
	activeEvalMaxAttempts = 1
	activeEvalResult      Evaluation
)

const restEvaluationResults = "api/v2/evaluation/applications/%s/results/%s"

func evaluationTestFunc(t *testing.T, w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet:
		expectedEndpoint := fmt.Sprintf(restEvaluationResults, activeEval.ApplicationID, activeEval.ResultID)
		if r.URL.String()[1:] != expectedEndpoint {
			t.Fatalf("Did not find expected URL: %s", expectedEndpoint)
		}

		if activeEvalAttempts < activeEvalMaxAttempts {
			activeEvalAttempts++
			w.WriteHeader(http.StatusNotFound)
		}

		resp, err := json.Marshal(activeEvalResult)
		if err != nil {
			w.WriteHeader(http.StatusTeapot)
		}

		fmt.Fprintln(w, string(resp))
	case r.Method == http.MethodPost:
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}

		var req iqEvaluationRequest
		if err = json.Unmarshal(body, &req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}

		activeEval.ResultID = "dummyResultID"
		activeEval.SubmittedDate = "Tomorrow and tomorrow and tomorrow"
		activeEval.ApplicationID = strings.Replace(r.URL.Path[1:], restEvaluation[:len(restEvaluation)-2], "", 1)
		activeEval.ResultsURL = fmt.Sprintf(restEvaluationResults, activeEval.ApplicationID, activeEval.ResultID)

		// Populate the Evaluation object
		activeEvalResult.SubmittedDate = activeEval.SubmittedDate
		activeEvalResult.ApplicationID = activeEval.ApplicationID
		activeEvalResult.EvaluationDate = "Today"
		for _, c := range req.Components {
			activeEvalResult.Results = append(activeEvalResult.Results, ComponentEvaluationResult{Component: c})
		}

		resp, err := json.Marshal(activeEval)
		if err != nil {
			w.WriteHeader(http.StatusTeapot)
		}

		fmt.Fprintln(w, string(resp))
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func evaluationTestIQ(t *testing.T) (iq IQ, mock *httptest.Server) {
	return newTestIQ(t, evaluationTestFunc)
}

func TestEvaluateComponents(t *testing.T) {
	iq, mock := evaluationTestIQ(t)
	defer mock.Close()

	var dummy Component
	dummy.Hash = "045c37a03be19f3e0db8"
	dummy.ComponentID = new(ComponentIdentifier)
	dummy.ComponentID.Format = "maven"
	dummy.ComponentID.Coordinates.ArtifactID = "jackson-databind"
	dummy.ComponentID.Coordinates.GroupID = "com.fasterxml.jackson.core"
	dummy.ComponentID.Coordinates.Version = "2.6.1"
	dummy.ComponentID.Coordinates.Extension = "jar"

	appID := "dummyAppId"

	report, err := EvaluateComponents(iq, []Component{dummy}, appID)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%v\n", report)

	if report.ApplicationID != appID {
		t.Errorf("AppID %s not found in report", appID)
	}

	if len(report.Results) != 1 {
		t.Errorf("Got %d results instead of the expected 1", len(report.Results))
	}

	reportComponent := report.Results[0].Component
	if !dummy.Equals(&reportComponent) {
		t.Error("Did not find the expected Component in evaluation results")
	}
}
