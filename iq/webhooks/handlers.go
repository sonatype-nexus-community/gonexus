package webhooks

import (
	"sync"
)

var mu sync.Mutex

var (
	appEvalEvents        = make(map[chan<- WebhookApplicationEvaluation]struct{})
	violationAlertEvents = make(map[chan<- WebhookViolationAlert]struct{})
)

// ApplicationEvaluationEvents puts a new Application Evaluation on the given channel
func ApplicationEvaluationEvents() (<-chan WebhookApplicationEvaluation, func()) {
	events := make(chan WebhookApplicationEvaluation, 1)

	mu.Lock()
	appEvalEvents[events] = struct{}{}
	mu.Unlock()

	return events, func() {
		mu.Lock()
		defer mu.Unlock()
		delete(appEvalEvents, events)
	}
}

// ViolationAlertEvents puts a new Violation Alert on the given channel
func ViolationAlertEvents() (<-chan WebhookViolationAlert, func()) {
	events := make(chan WebhookViolationAlert, 1)

	mu.Lock()
	violationAlertEvents[events] = struct{}{}
	mu.Unlock()

	return events, func() {
		mu.Lock()
		defer mu.Unlock()
		delete(violationAlertEvents, events)
	}
}

func sendApplicationEvaluationEvent(event WebhookApplicationEvaluation) {
	mu.Lock()
	defer mu.Unlock()
	for c := range appEvalEvents {
		select {
		case c <- event:
		default:
		}
	}
}

func sendViolationAlertEvent(event WebhookViolationAlert) {
	mu.Lock()
	defer mu.Unlock()
	for c := range violationAlertEvents {
		select {
		case c <- event:
		default:
		}
	}
}
