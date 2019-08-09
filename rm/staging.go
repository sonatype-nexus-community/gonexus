package nexusrm

import "fmt"

// service/rest/v1/staging/move/{repository}
const (
	restStaging       = "service/rest/v1/staging/move/%s"
	restStagingDelete = "service/rest/v1/staging/delete"
)

type stagingResponse struct {
	Status  int64  `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Destination     string            `json:"destination"`
		ComponentsMoved []componentsMoved `json:"components moved"`
	} `json:"data"`
}

type componentsMoved struct {
	Name    string `json:"name"`
	Group   string `json:"group"`
	Version string `json:"version"`
}

type stagingDeletionResponse struct {
	Status  int64  `json:"status"`
	Message string `json:"message"`
	Data    struct {
		ComponentsDeleted []componentsDeleted `json:"components deleted"`
	} `json:"data"`
}

type componentsDeleted struct {
	Repository string `json:"repository"`
	Group      string `json:"group"`
	Name       string `json:"name"`
	Version    string `json:"version"`
}

// StagingMove promotes components which match a set of criteria
func StagingMove(rm RM, query QueryBuilder) error {
	endpoint := fmt.Sprintf("%s?%s", restStaging, query.Build())

	// TODO: handle response
	_, _, err := rm.Post(endpoint, nil)
	return err
}

// StagingDelete removes components which have been staged
func StagingDelete(rm RM, query QueryBuilder) error {
	endpoint := fmt.Sprintf("%s?%s", restStaging, query.Build())

	_, err := rm.Del(endpoint)
	return err
}
