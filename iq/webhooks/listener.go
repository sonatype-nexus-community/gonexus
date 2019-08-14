package webhooks

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func parseRequest(r *http.Request) (whtype WebhookEventType, err error) {
	ok, whtype := IsWebhookEvent(r)
	if !ok {
		return whtype, errors.New("not a valid IQ webhook")
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return whtype, errors.New("could not read payload")
	}

	switch whtype {
	case WebhookEventApplicationEvaluation:
		var event WebhookApplicationEvaluation
		if err = json.Unmarshal(body, &event); err != nil {
			sendApplicationEvaluationEvent(event)
		}
	case WebhookEventViolationAlert:
		var event WebhookViolationAlert
		if err = json.Unmarshal(body, &event); err != nil {
			sendViolationAlertEvents(event)
		}
	default:
		return whtype, fmt.Errorf("IQ webhook type '%s' not supported", whtype)
	}

	return whtype, nil
}

// Listen will handle any HTTP requests which are genuine Nexus IQ Webhooks
func Listen(w http.ResponseWriter, r *http.Request) {
	whtype, err := parseRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch whtype {
	case WebhookEventApplicationEvaluation:
		log.Println("Accepted Application Evaluation")
	case WebhookEventViolationAlert:
		log.Println("Accepted Violation Alert")
	}

	w.WriteHeader(http.StatusOK)
}
