package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"

	"os"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/hcp-sdk-go/clients/cloud-vault-service/preview/2020-11-25/client/vault_service"
	vault "github.com/hashicorp/hcp-sdk-go/clients/cloud-vault-service/preview/2020-11-25/client/vault_service"

	sharedModels "github.com/hashicorp/hcp-sdk-go/clients/cloud-shared/v1/models"
	"github.com/hashicorp/hcp-sdk-go/clients/cloud-vault-service/preview/2020-11-25/models"

	"github.com/hashicorp/hcp-sdk-go/httpclient"
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
	// Export logs from a specified hour
	//

	auditParams := vault.NewFetchAuditLogParams()
	auditParams.LocationOrganizationID = orgID
	auditParams.LocationProjectID = projID
	auditParams.ClusterID = clusterID

	// Based on https://github.com/hashicorp/cloud-vault-service/blob/057972c079c7f3a82cd6011229d516f31e1672cb/test/e2e/smoke_test.go#L795-L812
	n := time.Now()
	auditParams.Body = &models.HashicorpCloudVault20201125FetchAuditLogRequest{
		ClusterID:     clusterID,
		IntervalStart: strfmt.DateTime(n.Add(-30 * time.Minute)),
		IntervalEnd:   strfmt.DateTime(n),

		// TODO: pull this from cluster above
		Location: &sharedModels.HashicorpCloudLocationLocation{
			OrganizationID: orgID,
			ProjectID:      projID,
			Region: &sharedModels.HashicorpCloudLocationRegion{
				Provider: "aws",
				Region:   "eu-west-2",
			},
		},
	}

	respAudit, err := vaultClient.FetchAuditLog(auditParams, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Log ID: %v", respAudit.Payload.LogID)

	//
	// Wait for log to be available, then download
	//

	downloadURL := ""

	for {
		// Based on https://github.com/hashicorp/cloud-vault-service/blob/057972c079c7f3a82cd6011229d516f31e1672cb/test/e2e/smoke_test.go#L814-L821
		params := vault_service.NewGetAuditLogStatusParams()
		params.LogID = respAudit.Payload.LogID
		params.ClusterID = clusterID
		params.LocationOrganizationID = orgID
		params.LocationProjectID = projID
		auditStatus, err := vaultClient.GetAuditLogStatus(params, nil)
		if err != nil {
			log.Fatal(err)
		}
		log.Infof("State: %v", auditStatus.Payload.Log.State)
		if !contains([]string{"PENDING", "CREATING"},
			string(auditStatus.Payload.Log.State)) {

			log.Infof("Log Payload: %v", auditStatus.Payload.Log)

			downloadURL = auditStatus.Payload.Log.DownloadURL
			break
		}

		time.Sleep(10 * time.Second)
	}
	log.Infof("Download URL: %v", downloadURL)

	gzippedLogsResponse, err := http.Get(downloadURL)
	if err != nil {
		log.Fatal(err)
	}
	defer gzippedLogsResponse.Body.Close()

	if gzippedLogsResponse.StatusCode == http.StatusOK {
		gzippedLogs, err := io.ReadAll(gzippedLogsResponse.Body)
		if err != nil {
			log.Fatal(err)
		}
		logs, err := gunzip(gzippedLogs)

		fmt.Printf("%s\n", logs)
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
