package nexusrm

import (
	"net/http"
)

const (
	restStatusReadable = "service/rest/v1/status"
	restStatusWritable = "service/rest/v1/status/writable"
)

// StatusReadable returns true if the RM instance can serve read requests
func StatusReadable(rm RM) (_ bool) {
	_, resp, err := rm.Get(restStatusReadable)
	return err == nil && resp.StatusCode == http.StatusOK
}

// StatusWritable returns true if the RM instance can serve read requests
func StatusWritable(rm RM) (_ bool) {
	_, resp, err := rm.Get(restStatusWritable)
	return err == nil && resp.StatusCode == http.StatusOK
}
