package nexusrm

// http://localhost:8081/service/rest/v1/assets?continuationToken=foo&repository=bar
const restListAssetsByRepo = "service/rest/v1/assets?repository=%s"

// RepositoryItemAssets describes the assets associated with a component
type RepositoryItemAssets struct {
	DownloadURL string `json:"downloadUrl"`
	Path        string `json:"path"`
	ID          string `json:"id"`
	Repository  string `json:"repository"`
	Format      string `json:"format"`
	Checksum    struct {
		Sha1 string `json:"sha1"`
		Md5  string `json:"md5"`
	} `json:"checksum"`
}

// Equals compares two RepositoryItemAssets objects
func (a *RepositoryItemAssets) Equals(b *RepositoryItemAssets) (_ bool) {
	if a == b {
		return true
	}

	if a.DownloadURL != b.DownloadURL {
		return
	}

	if a.Path != b.Path {
		return
	}

	if a.ID != b.ID {
		return
	}

	if a.Repository != b.Repository {
		return
	}

	if a.Format != b.Format {
		return
	}

	if a.Checksum.Sha1 != b.Checksum.Sha1 {
		return
	}

	if a.Checksum.Md5 != b.Checksum.Md5 {
		return
	}

	return true
}

type listAssetsResponse struct {
	Items             []RepositoryItemAssets `json:"items"`
	ContinuationToken string                 `json:"continuationToken"`
}

// GetAssets returns a list of assets in the indicated repository
func GetAssets(rm RM, repo string) (items []RepositoryItemAssets, err error) {
	return
}
