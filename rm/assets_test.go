package nexusrm

var dummyAssets = map[string][]RepositoryItemAsset{
	"repo-maven": []RepositoryItemAsset{
		/*
			RepositoryItem{ID: "component1id", Repository: "repo-maven", Format: "maven2", Group: "org.test", Name: "testComponent1", Version: "1.0.0"},
			RepositoryItem{ID: "component2id", Repository: "repo-maven", Format: "maven2", Group: "org.test", Name: "testComponent2", Version: "2.0.0"},
			RepositoryItem{ID: "component3id", Repository: "repo-maven", Format: "maven2", Group: "org.test", Name: "testComponent3", Version: "3.0.0"},
		*/
	},
	"repo-npm": []RepositoryItemAsset{
		// RepositoryItem{ID: "component4id", Repository: "repo-npm", Format: "maven2", Group: "org.test", Name: "testComponent4", Version: "4.0.0"},
	},
}

const dummyNewAssetID = "newAssetID"
