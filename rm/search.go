package nexusrm

/*
import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
)
*/

const (
	restSearchComponents = "service/rest/v1/search"
	restSearchAssets     = "service/rest/v1/search/assets"
	// restSearchAssetsDownload = "service/rest/v1/search/assets/download"
)

type searchComponentsResponse struct {
	Items             []RepositoryItem `json:"items"`
	ContinuationToken string           `json:"continuationToken"`
}

type searchAssetsResponse struct {
	Items             []RepositoryItemAssets `json:"items"`
	ContinuationToken string                 `json:"continuationToken"`
}

func search(rm RM, endpoint string, queryBuilder SearchQueryBuilder) ([]byte, error) {
	return []byte{}, nil
}

// func SearchComponents(rm RM, TODO) ([]RepositoryItem, error) {
//	b, err := search(rm, restSearchComponents, TODO)
// if err := json.Unmarshal(body, &item); err != nil {
// 	return item, doError(err)
// }
// }

// func SearchAssets(rm RM, TODO) ([]RepositoryItemAssets, error) {
// }

// GetComponents returns a list of components in the indicated repository
/*
func GetComponents(rm RM, repo string) ([]RepositoryItem, error) {
	continuation := ""
}
*/
