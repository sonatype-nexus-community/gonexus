package nexusiq

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"testing"
)

func createTempApplication(t *testing.T) (orgID string, appName string, appID string, err error) {
	rand.Seed(time.Now().UnixNano())
	name := strconv.Itoa(rand.Int())

	iq := getTestIQ(t)

	orgID, err = CreateOrganization(iq, name)
	if err != nil {
		return
	}

	appName = fmt.Sprintf("%s_app", name)

	appID, err = CreateApplication(iq, appName, orgID)
	if err != nil {
		return
	}

	return
}

func deleteTempApplication(t *testing.T, applicationName string) error {
	iq := getTestIQ(t)

	appInfo, err := GetApplicationDetailsByPublicID(iq, applicationName)
	if err != nil {
		return err
	}

	if err := DeleteApplication(iq, appInfo.ID); err != nil {
		return err
	}

	// if err := DeleteOrganization(iq, appInfo.OrganizationID); err != nil {
	// 	return err
	// }

	return nil
}

func ExampleGetAllApplications() {
	iq, err := New("http://localhost:8070", "username", "password")
	if err != nil {
		panic(err)
	}

	applications, err := GetAllApplications(iq)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", applications)
}
