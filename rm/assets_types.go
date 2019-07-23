package nexusrm

type repositoryItemAssetsChecksum struct {
	Sha1 string `json:"sha1"`
	Md5  string `json:"md5"`
}

// RepositoryItemAsset describes the assets associated with a component
type RepositoryItemAsset struct {
	DownloadURL string                       `json:"downloadUrl"`
	Path        string                       `json:"path"`
	ID          string                       `json:"id"`
	Repository  string                       `json:"repository"`
	Format      string                       `json:"format"`
	Checksum    repositoryItemAssetsChecksum `json:"checksum"`
}

// Equals compares two RepositoryItemAssets objects
func (a *RepositoryItemAsset) Equals(b *RepositoryItemAsset) (_ bool) {
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
	Items             []RepositoryItemAsset `json:"items"`
	ContinuationToken string                `json:"continuationToken"`
}
