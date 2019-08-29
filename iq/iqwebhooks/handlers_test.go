package iqwebhooks

import (
	"reflect"
	"testing"
)

func TestApplicationEvaluationEvents(t *testing.T) {
	tests := []struct {
		want   WebhookApplicationEvaluation
		closer bool
	}{
		{WebhookApplicationEvaluation{Initiator: "dummy1", ID: "foobar1"}, false},
		{WebhookApplicationEvaluation{Initiator: "dummy2", ID: "foobar2"}, true},
	}

	for _, test := range tests {
		events, close := ApplicationEvaluationEvents()
		if test.closer {
			defer close()
		}
		sendApplicationEvaluationEvent(test.want)

		got := <-events
		if !reflect.DeepEqual(got, test.want) {
			t.Error("Did not get expected event")
			t.Error("got", got)
			t.Error("want", test.want)
		}
	}
}

func TestViolationAlertEvents(t *testing.T) {
	tests := []struct {
		want   WebhookViolationAlert
		closer bool
	}{
		{WebhookViolationAlert{Initiator: "dummy1"}, false},
		{WebhookViolationAlert{Initiator: "dummy2"}, true},
	}

	for _, test := range tests {
		events, close := ViolationAlertEvents()
		if test.closer {
			defer close()
		}
		sendViolationAlertEvent(test.want)

		got := <-events
		if !reflect.DeepEqual(got, test.want) {
			t.Error("Did not get expected event")
			t.Error("got", got)
			t.Error("want", test.want)
		}
	}
}

func TestPolicyManagementEvents(t *testing.T) {
	tests := []struct {
		want   WebhookPolicyManagement
		closer bool
	}{
		{WebhookPolicyManagement{policyOwner{ID: "dummy1", Name: "foobar1"}}, false},
		{WebhookPolicyManagement{policyOwner{ID: "dummy2", Name: "foobar2"}}, true},
	}

	for _, test := range tests {
		events, close := PolicyManagementEvents()
		if test.closer {
			defer close()
		}
		sendPolicyManagementEvent(test.want)

		got := <-events
		if !reflect.DeepEqual(got, test.want) {
			t.Error("Did not get expected event")
			t.Error("got", got)
			t.Error("want", test.want)
		}
	}
}

func TestLicenseOverrideEvents(t *testing.T) {
	tests := []struct {
		want   WebhookLicenseOverride
		closer bool
	}{
		{WebhookLicenseOverride{licenseOverride{ID: "dummy1", OwnerID: "foobar1"}}, false},
		{WebhookLicenseOverride{licenseOverride{ID: "dummy2", OwnerID: "foobar2"}}, true},
	}

	for _, test := range tests {
		events, close := LicenseOverrideEvents()
		if test.closer {
			defer close()
		}
		sendLicenseOverrideEvent(test.want)

		got := <-events
		if !reflect.DeepEqual(got, test.want) {
			t.Error("Did not get expected event")
			t.Error("got", got)
			t.Error("want", test.want)
		}
	}
}

func TestSecurityOverrideEvents(t *testing.T) {
	tests := []struct {
		want   WebhookSecurityOverride
		closer bool
	}{
		{WebhookSecurityOverride{securityVulnerabilityOverride{ID: "dummy1", OwnerID: "foobar1"}}, false},
		{WebhookSecurityOverride{securityVulnerabilityOverride{ID: "dummy2", OwnerID: "foobar2"}}, true},
	}

	for _, test := range tests {
		events, close := SecurityOverrideEvents()
		if test.closer {
			defer close()
		}
		sendSecurityOverrideEvent(test.want)

		got := <-events
		if !reflect.DeepEqual(got, test.want) {
			t.Error("Did not get expected event")
			t.Error("got", got)
			t.Error("want", test.want)
		}
	}
}
