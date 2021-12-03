package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-openapi/strfmt"
	sharedModels "github.com/hashicorp/hcp-sdk-go/clients/cloud-shared/v1/models"
	"github.com/hashicorp/hcp-sdk-go/clients/cloud-vault-service/preview/2020-11-25/client/vault_service"
	vault "github.com/hashicorp/hcp-sdk-go/clients/cloud-vault-service/preview/2020-11-25/client/vault_service"
	"github.com/hashicorp/hcp-sdk-go/clients/cloud-vault-service/preview/2020-11-25/models"

	log "github.com/sirupsen/logrus"
)

// TODO: Create an AuditLogRequest type for this
func downloadAndPrintLogs(vaultClient vault.ClientService,
	clusterID, orgID, projID string,
	start, end strfmt.DateTime) {
	log.Infof("Downloading logs from %v to %v", start, end)

	//
	// Export logs from a specified hour
	//
	auditParams := vault.NewFetchAuditLogParams()
	auditParams.LocationOrganizationID = orgID
	auditParams.LocationProjectID = projID
	auditParams.ClusterID = clusterID
	auditParams.Body = &models.HashicorpCloudVault20201125FetchAuditLogRequest{
		ClusterID:     clusterID,
		IntervalStart: start,
		IntervalEnd:   end,
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

	// TODO: these logs will always start at the top of the hour, regardless of the startTime

}
