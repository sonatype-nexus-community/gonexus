package nexusrm

import (
	"bytes"
)

type SearchSort int

const (
	None SearchSort = iota
	Group
	Name
	Version
	Repo
)

func (s SearchSort) String() string {
	return ""
}

type SearchDirection int

const (
	Asc SearchDirection = iota
	Desc
)

func (s SearchDirection) String() string {
	return ""
}

type SearchQueryBuilder struct {
	sort      SearchSort
	direction SearchDirection
	criteria  map[string]string
}

func (b *SearchQueryBuilder) addCriteria(c, v string) *SearchQueryBuilder {
	b.criteria[c] = v
	return b
}

func (b *SearchQueryBuilder) build() string {
	var buf bytes.Buffer

	if b.sort != None {
		buf.WriteString("sort=")
		buf.WriteString(b.sort.String())
		buf.WriteString("&")
	}

	buf.WriteString("direction=")
	buf.WriteString(b.direction.String())

	for k, v := range b.criteria {
		buf.WriteString("&")
		buf.WriteString(k)
		buf.WriteString("=")
		buf.WriteString(v)
	}

	return buf.String()
}

func (b *SearchQueryBuilder) Sort(v SearchSort) *SearchQueryBuilder {
	b.sort = v
	return b
}

func (b *SearchQueryBuilder) Direction(v SearchDirection) *SearchQueryBuilder {
	b.direction = v
	return b
}

func (b *SearchQueryBuilder) Q(v string) *SearchQueryBuilder {
	return b.addCriteria("q", v)
}

func (b *SearchQueryBuilder) Repository(v string) *SearchQueryBuilder {
	return b.addCriteria("repository", v)
}

func (b *SearchQueryBuilder) Format(v string) *SearchQueryBuilder {
	return b.addCriteria("format", v)
}

func (b *SearchQueryBuilder) Tag(v string) *SearchQueryBuilder {
	return b.addCriteria("tag", v)
}

func (b *SearchQueryBuilder) Group(v string) *SearchQueryBuilder {
	return b.addCriteria("group", v)
}

func (b *SearchQueryBuilder) Name(v string) *SearchQueryBuilder {
	return b.addCriteria("name", v)
}

func (b *SearchQueryBuilder) Version(v string) *SearchQueryBuilder {
	return b.addCriteria("version", v)
}

func (b *SearchQueryBuilder) Md5(v string) *SearchQueryBuilder {
	return b.addCriteria("md5", v)
}

func (b *SearchQueryBuilder) Sha1(v string) *SearchQueryBuilder {
	return b.addCriteria("sha1", v)
}

func (b *SearchQueryBuilder) Sha256(v string) *SearchQueryBuilder {
	return b.addCriteria("sha256", v)
}

func (b *SearchQueryBuilder) Sha512(v string) *SearchQueryBuilder {
	return b.addCriteria("sha512", v)
}

func (b *SearchQueryBuilder) Prerelease(v string) *SearchQueryBuilder {
	return b.addCriteria("prerelease", v)
}

func (b *SearchQueryBuilder) DockerImageName(v string) *SearchQueryBuilder {
	return b.addCriteria("docker.imageName", v)
}

func (b *SearchQueryBuilder) DockerImageTag(v string) *SearchQueryBuilder {
	return b.addCriteria("docker.imageTag", v)
}

func (b *SearchQueryBuilder) DockerLayerId(v string) *SearchQueryBuilder {
	return b.addCriteria("docker.layerId", v)
}

func (b *SearchQueryBuilder) DockerContentDigest(v string) *SearchQueryBuilder {
	return b.addCriteria("docker.contentDigest", v)
}

func (b *SearchQueryBuilder) MavenGroupId(v string) *SearchQueryBuilder {
	return b.addCriteria("maven.groupId", v)
}

func (b *SearchQueryBuilder) MavenArtifactId(v string) *SearchQueryBuilder {
	return b.addCriteria("maven.artifactId", v)
}

func (b *SearchQueryBuilder) MavenBaseVersion(v string) *SearchQueryBuilder {
	return b.addCriteria("maven.baseVersion", v)
}

func (b *SearchQueryBuilder) MavenExtension(v string) *SearchQueryBuilder {
	return b.addCriteria("maven.extension", v)
}

func (b *SearchQueryBuilder) MavenClassifier(v string) *SearchQueryBuilder {
	return b.addCriteria("maven.classifier", v)
}

func (b *SearchQueryBuilder) NpmScope(v string) *SearchQueryBuilder {
	return b.addCriteria("npm.scope", v)
}

func (b *SearchQueryBuilder) NugetId(v string) *SearchQueryBuilder {
	return b.addCriteria("nuget.id", v)
}

func (b *SearchQueryBuilder) NugetTags(v string) *SearchQueryBuilder {
	return b.addCriteria("nuget.tags", v)
}

func (b *SearchQueryBuilder) PypiClassifiers(v string) *SearchQueryBuilder {
	return b.addCriteria("pypi.classifiers", v)
}

func (b *SearchQueryBuilder) PypiDescription(v string) *SearchQueryBuilder {
	return b.addCriteria("pypi.description", v)
}

func (b *SearchQueryBuilder) PypiKeywords(v string) *SearchQueryBuilder {
	return b.addCriteria("pypi.keywords", v)
}

func (b *SearchQueryBuilder) PypiSummary(v string) *SearchQueryBuilder {
	return b.addCriteria("pypi.summary", v)
}

func (b *SearchQueryBuilder) RubygemsDescription(v string) *SearchQueryBuilder {
	return b.addCriteria("rubygems.description", v)
}

func (b *SearchQueryBuilder) RubygemsPlatform(v string) *SearchQueryBuilder {
	return b.addCriteria("rubygems.platform", v)
}

func (b *SearchQueryBuilder) RubygemsSummary(v string) *SearchQueryBuilder {
	return b.addCriteria("rubygems.summary", v)
}

func (b *SearchQueryBuilder) YumArchitecture(v string) *SearchQueryBuilder {
	return b.addCriteria("yum.architecture", v)
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
