package nexusrm

import (
	"bytes"
	"fmt"
	"text/template"
)

/*
type blobStoreDelete struct {
	Name string
}

const groovyDeleteBlobStore = `blobStore.delete('{{.Name}}')`
*/

type blobStoreFile struct {
	Name, Path string
}

const groovyCreateFileBlobStore = `blobStore.createFileBlobStore('{{.Name}}', '{{.Path}}')`

// BlobStoreS3 encapsulates the needed options for creating an S3 blob store
type BlobStoreS3 struct {
	Name, BucketName, AwsAccessKey, AwsSecret, AwsIamRole, AwsRegion string
}

const groovyCreateS3BlobStore = `def config = [:]
config['bucket'] = '{{.BucketName}}'
config['accessKeyId'] = '{{.AwsAccessKey}}'
config['secretAccessKey'] = '{{.AwsSecret}}'
config['assumeRole'] = '{{.AwsIamRole}}'
config['region'] = '{{.AwsRegion}}'
blobStore.createS3BlobStore('{{.Name}}', config)`

type blobStoreGroup struct {
	Name       string
	BlobStores []string
}

const groovyCreateBlobStoreGroup = `blobStore.createBlobStoreGroup('{{.Name}}', [{{range .BlobStores}}'{{.}}',{{end}}], 'writeToFirst')`

/*
// DeleteBlobStore creates a blobstore
func DeleteBlobStore(rm RM, name string) error {
	tmpl, err := template.New("dbs").Parse(groovyDeleteBlobStore)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, blobStoreDelete{name})
	if err != nil {
		return err
	}

	_, err = ScriptRunOnce(rm, newAnonGroovyScript(buf.String()), nil)
	return err
}
*/

// CreateFileBlobStore creates a blobstore
func CreateFileBlobStore(rm RM, name, path string) error {
	tmpl, err := template.New("fbs").Parse(groovyCreateFileBlobStore)
	if err != nil {
		return fmt.Errorf("could not parse template: %v", err)
	}

	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, blobStoreFile{name, path})
	if err != nil {
		return fmt.Errorf("could not create file blobstore from template: %v", err)
	}

	_, err = ScriptRunOnce(rm, newAnonGroovyScript(buf.String()), nil)
	return fmt.Errorf("could not create file blobstore: %v", err)
}

// CreateBlobStoreGroup creates a blobstore
func CreateBlobStoreGroup(rm RM, name string, blobStores []string) error {
	tmpl, err := template.New("group").Parse(groovyCreateBlobStoreGroup)
	if err != nil {
		return fmt.Errorf("could not parse template: %v", err)
	}

	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, blobStoreGroup{name, blobStores})
	if err != nil {
		return fmt.Errorf("could not create group blobstore from template: %v", err)
	}

	_, err = ScriptRunOnce(rm, newAnonGroovyScript(buf.String()), nil)
	return fmt.Errorf("could not create group blobstore: %v", err)
}
