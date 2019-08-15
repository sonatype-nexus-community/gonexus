package webhooks

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestParseRequest(t *testing.T) {
	tests := []struct {
		t WebhookEventType
		e interface{}
	}{
		{WebhookEventApplicationEvaluation, WebhookApplicationEvaluation{
			Initiator: "dummy",
			ID:        "foobar",
		}},
		{WebhookEventViolationAlert, WebhookViolationAlert{
			Initiator: "dummy",
		}},
	}

	for _, test := range tests {
		buf, err := json.Marshal(test.e)
		if err != nil {
			t.Errorf("could not marshal test event: %v", err)
		}

		req := httptest.NewRequest(http.MethodPost, "http://foo.bar", bytes.NewBuffer(buf))
		req.Header.Set("Content-Type", "application/json; charset=UTF-8")
		req.Header.Set("User-Agent", "Sonatype_CLM_Server/1.70.0 (Java 1.8.0)")
		req.Header.Set("X-Nexus-Webhook-Id", string(test.t))

		got, err := parseRequest(req)
		if err != nil {
			t.Errorf("unexpected error parsing request: %v", err)
		}
		if got != test.t {
			t.Error("did not identify expected type")
			t.Error("got", got)
			t.Error("want", test.t)
		}
	}
}

func TestReceivingEventsFromRequest(t *testing.T) {
	tests := []struct {
		t WebhookEventType
		e interface{}
	}{
		{WebhookEventApplicationEvaluation, WebhookApplicationEvaluation{
			Initiator: "dummy",
			ID:        "foobar",
		}},
		{WebhookEventViolationAlert, WebhookViolationAlert{
			Initiator: "dummy",
		}},
	}

	appEvalEvents, _ := ApplicationEvaluationEvents()
	violationAlertEvents, _ := ViolationAlertEvents()

	for _, test := range tests {
		buf, err := json.Marshal(test.e)
		if err != nil {
			t.Errorf("could not marshal test event: %v", err)
		}

		req := httptest.NewRequest(http.MethodPost, "http://foo.bar", bytes.NewBuffer(buf))
		req.Header.Set("Content-Type", "application/json; charset=UTF-8")
		req.Header.Set("User-Agent", "Sonatype_CLM_Server/1.70.0 (Java 1.8.0)")
		req.Header.Set("X-Nexus-Webhook-Id", string(test.t))

		gotType, err := parseRequest(req)
		if err != nil {
			t.Errorf("unexpected error parsing request: %v", err)
		}
		if gotType != test.t {
			t.Error("did not identify expected type")
			t.Error("got", gotType)
			t.Error("want", test.t)
		}

		select {
		case got := <-appEvalEvents:
			if test.t != WebhookEventApplicationEvaluation {
				t.Error("did not identify expected type")
				t.Error("want", test.t)
			}
			t.Log(got)
		case got := <-violationAlertEvents:
			if test.t != WebhookEventViolationAlert {
				t.Error("did not identify expected type")
				t.Error("want", test.t)
			}
			t.Log(got)
		default:
			t.Error("did not get expected anything")
		}
	}
}
