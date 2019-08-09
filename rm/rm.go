package nexusrm

import (
	"bytes"

	nexus "github.com/sonatype-nexus-community/gonexus"
)

// RM is the interface which any Repository Manager implementation would need to satisfy
type RM interface {
	nexus.Client
}

type rmClient struct {
	nexus.DefaultClient
}

// New creates a new Repository Manager instance
func New(host, username, password string) (RM, error) {
	rm := new(rmClient)
	rm.Host = host
	rm.Username = username
	rm.Password = password
	return rm, nil
}

// QueryBuilder allows you to build a search query
type QueryBuilder struct {
	criteria map[string]string
}

func (b *QueryBuilder) addCriteria(c, v string) *QueryBuilder {
	b.criteria[c] = v
	return b
}

func (b *QueryBuilder) buildCriteria(buf *bytes.Buffer) {
	for k, v := range b.criteria {
		buf.WriteString("&")
		buf.WriteString(k)
		buf.WriteString("=")
		buf.WriteString(v)
	}
}

// Build will build the assembled search query
func (b *QueryBuilder) Build() string {
	var buf bytes.Buffer

	b.buildCriteria(&buf)

	return buf.String()
}

// Q allows specifying a keyword search
func (b *QueryBuilder) Q(v string) *QueryBuilder {
	return b.addCriteria("q", v)
}

// Repository allows specifying the repository to search
func (b *QueryBuilder) Repository(v string) *QueryBuilder {
	return b.addCriteria("repository", v)
}

// Format allows specifiying the format to filter by
func (b *QueryBuilder) Format(v string) *QueryBuilder {
	return b.addCriteria("format", v)
}

// Tag allows specifiying a tag to filter by
func (b *QueryBuilder) Tag(v string) *QueryBuilder {
	return b.addCriteria("tag", v)
}

// Group allows specifiying a group to filter by
func (b *QueryBuilder) Group(v string) *QueryBuilder {
	return b.addCriteria("group", v)
}

// Name allows specifiying a name to filter by
func (b *QueryBuilder) Name(v string) *QueryBuilder {
	return b.addCriteria("name", v)
}

// Version allows specifiying a version to filter by
func (b *QueryBuilder) Version(v string) *QueryBuilder {
	return b.addCriteria("version", v)
}

// Md5 allows specifiying an md5 sum to filter by
func (b *QueryBuilder) Md5(v string) *QueryBuilder {
	return b.addCriteria("md5", v)
}

// Sha1 allows specifiying an sha1 sum to filter by
func (b *QueryBuilder) Sha1(v string) *QueryBuilder {
	return b.addCriteria("sha1", v)
}

// Sha256 allows specifiying an sha256 sum to filter by
func (b *QueryBuilder) Sha256(v string) *QueryBuilder {
	return b.addCriteria("sha256", v)
}

// Sha512 allows specifiying an sha512 sum to filter by
func (b *QueryBuilder) Sha512(v string) *QueryBuilder {
	return b.addCriteria("sha512", v)
}

// Prerelease allows specifiying a prerelease qualifier to filter by
func (b *QueryBuilder) Prerelease(v string) *QueryBuilder {
	return b.addCriteria("prerelease", v)
}

// DockerImageName allows specifiying the name of a docker image to filter by
func (b *QueryBuilder) DockerImageName(v string) *QueryBuilder {
	return b.addCriteria("docker.imageName", v)
}

// DockerImageTag allows specifiying the tag of a docker image to filter by
func (b *QueryBuilder) DockerImageTag(v string) *QueryBuilder {
	return b.addCriteria("docker.imageTag", v)
}

// DockerLayerID allows specifiying the ID of a docker image layer to filter by
func (b *QueryBuilder) DockerLayerID(v string) *QueryBuilder {
	return b.addCriteria("docker.layerId", v)
}

// DockerContentDigest allows specifiying the digest of docker layers to filter by
func (b *QueryBuilder) DockerContentDigest(v string) *QueryBuilder {
	return b.addCriteria("docker.contentDigest", v)
}

// MavenGroupID allows specifiying the group name/id of maven component to filter by
func (b *QueryBuilder) MavenGroupID(v string) *QueryBuilder {
	return b.addCriteria("maven.groupId", v)
}

// MavenArtifactID allows specifiying the artifact id of maven component to filter by
func (b *QueryBuilder) MavenArtifactID(v string) *QueryBuilder {
	return b.addCriteria("maven.artifactId", v)
}

// MavenBaseVersion allows specifiying the version of maven component to filter by
func (b *QueryBuilder) MavenBaseVersion(v string) *QueryBuilder {
	return b.addCriteria("maven.baseVersion", v)
}

// MavenExtension allows specifiying the extension of maven component to filter by
func (b *QueryBuilder) MavenExtension(v string) *QueryBuilder {
	return b.addCriteria("maven.extension", v)
}

// MavenClassifier allows specifiying the classifier of maven component to filter by
func (b *QueryBuilder) MavenClassifier(v string) *QueryBuilder {
	return b.addCriteria("maven.classifier", v)
}

// NpmScope allows specifiying the scope of an NPM component to filter by
func (b *QueryBuilder) NpmScope(v string) *QueryBuilder {
	return b.addCriteria("npm.scope", v)
}

// NugetID allows specifiying the ID/name of a Nuget component to filter by
func (b *QueryBuilder) NugetID(v string) *QueryBuilder {
	return b.addCriteria("nuget.id", v)
}

// NugetTags allows specifiying the tags of a Nuget component to filter by
func (b *QueryBuilder) NugetTags(v string) *QueryBuilder {
	return b.addCriteria("nuget.tags", v)
}

// PypiClassifiers allows specifiying the classifiers of a pypi component to filter by
func (b *QueryBuilder) PypiClassifiers(v string) *QueryBuilder {
	return b.addCriteria("pypi.classifiers", v)
}

// PypiDescription allows specifiying the description of a pypi component to filter by
func (b *QueryBuilder) PypiDescription(v string) *QueryBuilder {
	return b.addCriteria("pypi.description", v)
}

// PypiKeywords allows specifiying the keywords of a pypi component to filter by
func (b *QueryBuilder) PypiKeywords(v string) *QueryBuilder {
	return b.addCriteria("pypi.keywords", v)
}

// PypiSummary allows specifiying the summary of a pypi component to filter by
func (b *QueryBuilder) PypiSummary(v string) *QueryBuilder {
	return b.addCriteria("pypi.summary", v)
}

// RubygemsDescription allows specifiying the description of a ruby gem to filter by
func (b *QueryBuilder) RubygemsDescription(v string) *QueryBuilder {
	return b.addCriteria("rubygems.description", v)
}

// RubygemsPlatform allows specifiying the platform of a ruby gem to filter by
func (b *QueryBuilder) RubygemsPlatform(v string) *QueryBuilder {
	return b.addCriteria("rubygems.platform", v)
}

// RubygemsSummary allows specifiying the summary of a ruby gem to filter by
func (b *QueryBuilder) RubygemsSummary(v string) *QueryBuilder {
	return b.addCriteria("rubygems.summary", v)
}

// YumArchitecture allows specifiying the architecture of a Yum package to filter by
func (b *QueryBuilder) YumArchitecture(v string) *QueryBuilder {
	return b.addCriteria("yum.architecture", v)
}

// NewQueryBuilder creates a new instance of QueryBuilder
func NewQueryBuilder() *QueryBuilder {
	b := new(QueryBuilder)
	b.criteria = make(map[string]string)
	return b
}
