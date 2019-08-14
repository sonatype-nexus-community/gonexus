package webhooks

var appEvalEvents = make([]chan<- WebhookApplicationEvaluation, 0)
var violationAlertEvents = make([]chan<- WebhookViolationAlert, 0)

// ApplicationEvaluationEvents puts a new Application Evaluation on the given channel
func ApplicationEvaluationEvents(events chan<- WebhookApplicationEvaluation) {
	appEvalEvents = append(appEvalEvents, events)
}

// ViolationAlertEvents puts a new Violation Alert on the given channel
func ViolationAlertEvents(events chan<- WebhookViolationAlert) {
	violationAlertEvents = append(violationAlertEvents, events)
}

func sendApplicationEvaluationEvent(event WebhookApplicationEvaluation) {
	for _, c := range appEvalEvents {
		c <- event
	}
}

func sendViolationAlertEvents(event WebhookViolationAlert) {
	for _, c := range violationAlertEvents {
		c <- event
	}
}
