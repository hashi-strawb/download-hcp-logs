package main

import (
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"time"

	"github.com/go-openapi/strfmt"
	vault "github.com/hashicorp/hcp-sdk-go/clients/cloud-vault-service/preview/2020-11-25/client/vault_service"
	"github.com/hashicorp/hcp-sdk-go/httpclient"
	_ "github.com/joho/godotenv/autoload"
	log "github.com/sirupsen/logrus"
)

func main() {

	// Initialize SDK http client
	cl, err := httpclient.New(httpclient.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Import versioned client for each service.
	vaultClient := vault.New(cl, nil)

	// These IDs can be obtained from the portal URL
	orgID := os.Getenv("HCP_ORGANIZATION_ID")
	projID := os.Getenv("HCP_PROJECT_ID")
	clusterID := os.Getenv("HCP_VAULT_CLUSTER_ID")

	//
	// List Clusters
	// (we're not actually doing anything here, but it would be useful in future)
	//

	listParams3 := vault.NewListParams()
	listParams3.LocationOrganizationID = orgID
	listParams3.LocationProjectID = projID

	resp3, err := vaultClient.List(listParams3, nil)
	if err != nil {
		log.Fatal(err)
	}
	if len(resp3.Payload.Clusters) > 0 {
		log.Infof("Response: %v", resp3.Payload.Clusters[0])
	} else {
		log.Infof("Response: %v", resp3.Payload)
	}

	//
	// Download and print audit logs
	//

	now := time.Now()
	startTime := strfmt.DateTime(now.Add(-10 * time.Minute))
	endTime := strfmt.DateTime(now)

	for {
		downloadAndPrintLogs(vaultClient,
			clusterID, orgID, projID,
			startTime, endTime)

		// Wait at least a minute between requests, otherwise we get duplicates
		time.Sleep(1 * time.Minute)

		startTime = endTime
		endTime = strfmt.DateTime(time.Now())

		// TODO: ensure startTime is at least 10m before endTime to account for
		// peculiarities with the logs around the top of the hour
	}

}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func gunzip(data []byte) ([]byte, error) {
	b := bytes.NewBuffer(data)

	var r io.Reader
	r, err := gzip.NewReader(b)
	if err != nil {
		return []byte{}, err
	}

	var resB bytes.Buffer
	_, err = resB.ReadFrom(r)
	if err != nil {
		return []byte{}, err
	}

	return resB.Bytes(), nil
}
