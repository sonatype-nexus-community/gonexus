package iqwebhooks

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var tests = []struct {
	t WebhookEventType
	e interface{}
}{
	{WebhookEventApplicationEvaluation, WebhookApplicationEvaluation{Initiator: "dummy", ID: "foobar"}},
	{WebhookEventViolationAlert, WebhookViolationAlert{Initiator: "dummy"}},
	{WebhookEventPolicyManagement, WebhookPolicyManagement{policyOwner{ID: "dummy1", Name: "foobar1"}}},
	{WebhookEventLicenseOverride, WebhookLicenseOverride{licenseOverride{ID: "dummy1", OwnerID: "foobar1"}}},
	{WebhookEventSecurityOverride, WebhookSecurityOverride{securityVulnerabilityOverride{ID: "dummy1", OwnerID: "foobar1"}}},
}

func TestParseRequest(t *testing.T) {
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
	appEvalEvents, _ := ApplicationEvaluationEvents()
	violationAlertEvents, _ := ViolationAlertEvents()
	policyManagementEvents, _ := PolicyManagementEvents()
	licenseOverrideEvents, _ := LicenseOverrideEvents()
	securityOverrideEvents, _ := SecurityOverrideEvents()

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
		case got := <-policyManagementEvents:
			if test.t != WebhookEventPolicyManagement {
				t.Error("did not identify expected type")
				t.Error("want", test.t)
			}
			t.Log(got)
		case got := <-licenseOverrideEvents:
			if test.t != WebhookEventLicenseOverride {
				t.Error("did not identify expected type")
				t.Error("want", test.t)
			}
			t.Log(got)
		case got := <-securityOverrideEvents:
			if test.t != WebhookEventSecurityOverride {
				t.Error("did not identify expected type")
				t.Error("want", test.t)
			}
			t.Log(got)
		default:
			t.Error("did not get expected anything")
		}
	}
}

func ExampleListen() {
	appEvalEvents, _ := ApplicationEvaluationEvents()
	violationAlertEvents, _ := ViolationAlertEvents()
	policyMgmtEvents, _ := PolicyManagementEvents()
	licenseOverride, _ := LicenseOverrideEvents()
	securityOverride, _ := SecurityOverrideEvents()

	go func() {
		for {
			select {
			case <-appEvalEvents:
				log.Println("Received Application Evaluation event")
			case <-violationAlertEvents:
				log.Println("Received Violation Alert event")
			case <-policyMgmtEvents:
				log.Println("Received Policy Management event")
			case <-licenseOverride:
				log.Println("Received License Overridden event")
			case <-securityOverride:
				log.Println("Received Security Vulnerability Overridden event")
			default:
			}
		}
	}()

	http.HandleFunc("/ingest", Listen)

	log.Fatal(http.ListenAndServe(":9876", nil))
}
