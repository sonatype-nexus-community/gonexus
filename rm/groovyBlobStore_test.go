package nexusrm

import (
	"testing"
)

func TestCreateFileBlobStore(t *testing.T) {
	rm := getTestRM(t)

	err := CreateFileBlobStore(rm, "testname", "testpath")
	if err != nil {
		t.Error(err)
	}

	// TODO: list blobstores
}

func TestCreateBlobStoreGroup(t *testing.T) {
	rm := getTestRM(t)

	CreateFileBlobStore(rm, "f1", "pathf1")
	CreateFileBlobStore(rm, "f2", "pathf2")
	CreateFileBlobStore(rm, "f3", "pathf3")

	err := CreateBlobStoreGroup(rm, "grpname", []string{"f1", "f2", "f3"})
	if err != nil {
		t.Error(err)
	}
}

/*
func TestDeleteBlobStore(t *testing.T) {
	rm := getTestRM(t)

	err := DeleteBlobStore(rm, "testname")
	if err != nil {
		t.Error(err)
	}

	// TODO: list blobstores
}
*/
