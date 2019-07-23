package nexusrm

type SearchQueryBuilder interface {
	build() string
}

/*
continuationToken string (query)

sort string (query)	[Available values : group, name, version, repository]
direction string (query) [Available values : asc, desc]
q string (query)

repository string (query)
format string (query)
tag string (query)
group string (query)
name string (query)
version string (query)
md5 string (query)
sha1 string (query)
sha256 string (query)
sha512 string (query)
prerelease string (query)

docker.imageName string (query)
docker.imageTag string (query)
docker.layerId string (query)
docker.contentDigest string (query)
maven.groupId string (query)
maven.artifactId string (query)
maven.baseVersion string (query)
maven.extension string (query)
maven.classifier string (query)
npm.scope string (query)
nuget.id string (query)
nuget.tags string (query)
pypi.classifiers string (query)
pypi.description string (query)
pypi.keywords string (query)
pypi.summary string (query)
rubygems.description string (query)
rubygems.platform string (query)
rubygems.summary string (query)
yum.architecture string (query)
*/
