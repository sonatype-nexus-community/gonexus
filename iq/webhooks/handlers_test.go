package webhooks

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

func TestViolationAlertouEvents(t *testing.T) {
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
