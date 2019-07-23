package nexusrm

const restAssets = "service/rest/v1/assets?repository=%s"
const restListAssetsByRepo = "service/rest/v1/assets?repository=%s"

// GetAssets returns a list of assets in the indicated repository
func GetAssets(rm RM, repo string) (items []RepositoryItemAsset, err error) {
	return
}

// GetAsset returns an asset by ID
func GetAsset(rm RM, id string) (items RepositoryItemAsset, err error) {
	return
}

// DeleteAssetByID deletes the asset indicated by ID
func DeleteAssetByID(rm RM, id string) error {
	return nil
}
