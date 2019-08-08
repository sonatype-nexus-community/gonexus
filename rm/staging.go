package nexusrm

import nexus "github.com/sonatype-nexus-community/gonexus"

const restStaging = "service/rest/v1/staging/move/{repository}"

func StagingMove(rm RM, query nexus.SearchQueryBuilder) ([]RepositoryItem, error) {
	return nil, nil
}

func StagingDelete(rm RM, query nexus.SearchQueryBuilder) ([]RepositoryItem, error) {
	return nil, nil
}
