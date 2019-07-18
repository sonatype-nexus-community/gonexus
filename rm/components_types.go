package nexusrm

import (
	"fmt"
	"os"
)

// UploadComponent defines the interface which describes a component to upload
type uploadComponent interface {
	formData() (fields map[string]string, files map[string]*os.File)
}

// UploadAssetMaven encapsulates data needed to upload an maven2 asset
type UploadAssetMaven struct {
	File                  *os.File
	Classifier, Extension string
}

// UploadComponentMaven encapsulates data needed to upload an maven2 component
type UploadComponentMaven struct {
	GroupID, ArtifactID, Version, Packaging, Tag string
	GeneratePom                                  bool
	Assets                                       []UploadAssetMaven
}

func (a UploadComponentMaven) formData() (map[string]string, map[string]*os.File) {
	fields := make(map[string]string)
	files := make(map[string]*os.File)

	fields["maven2.groupId"] = a.GroupID
	fields["maven2.artifactId"] = a.ArtifactID
	fields["maven2.version"] = a.Version
	fields["maven2.packaging"] = a.Packaging
	fields["maven2.tag"] = a.Tag
	fields["maven2.generate-pom"] = fmt.Sprintf("%T", a.GeneratePom)

	for i, a := range a.Assets {
		fieldName := fmt.Sprintf("maven2.asset%d", i)
		files[fieldName] = a.File
		fields[fieldName+".classifier"] = a.Classifier
		fields[fieldName+".extension"] = a.Extension
	}

	return fields, files
}

// UploadComponentNpm encapsulates data needed to upload an NPM component
type UploadComponentNpm struct {
	File *os.File
	Tag  string
}
