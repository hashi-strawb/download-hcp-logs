# Example Go Script to pull HCP Vault Audit Logs

WARNING: This makes use of unstable preview APIs which could change at any time!

USE AT YOUR OWN PERIL


## Why?

HCP Vault has audit logs! Woo!

HCP Vault only supports streaming those to a small number of places! Boo!

HCP Vault also only supports downloading audit logs through the HCP Portal manually! Double boo!

Fortunately, there is an official API on the way for this.

In the short term, we can use the preview API

This script uses https://github.com/hashicorp/hcp-sdk-go to interact with that.



## Requirements

Environment variables to define the cluster:

```
export HCP_ORGANIZATION_ID=<YOUR HCP ORG HERE>
export HCP_PROJECT_ID=<YOUR HCP PROJECT HERE>
export HCP_VAULT_CLUSTER_ID=vault-cluster
```

Environment variables for HCP Auth:

```
export HCP_CLIENT_ID=<YOUR CLIENT ID HERE>
export HCP_CLIENT_SECRET=<YOUR CLIENT SECRET HERE>
```

## Running

```
go run main.go
```

Example output:

```
$ go run main.go
Response: &models.HashicorpCloudVault20201125Cluster{Config:(*models.HashicorpCloudVault20201125ClusterConfig)(0xc0004acfc0), CreatedAt:strfmt.DateTime{wall:0x1464cf58, ext:63773966983, loc:(*time.Location)(nil)}, CurrentVersion:"v1.8.5", DNSNames:(*models.HashicorpCloudVault20201125ClusterDNSNames)(0xc0000e2040), ID:"vault-cluster", Location:(*models.HashicorpCloudLocationLocation)(0xc0000bb9e0), State:"RUNNING"}

Response: "a6c7a6d5-0fa5-4193-9647-ad70f13f34da"

State: PENDING
State: CREATING
State: READYResponse: &models.HashicorpCloudVault20201125AuditLog{ClusterID:"vault-cluster", DownloadURL:"https://hcp-data-plane-blob-prod.s3.amazonaws.com/225af347-0fd9-41b1-8571-bed0d3ef665e/auditlogs/98b8bfa8-3115-4831-a598-434d02f83786?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=ASIAZDKPXWD4FHMAK4HM%2F20211201%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20211201T160154Z&X-Amz-Expires=900&X-Amz-Security-Token=FwoGZXIvYXdzEFEaDIKka9XFbRoODVjZSiLjAVJ4U5tMj8mqcuc9c2hzVpvCy6jj2TtLieIkyjTRbFOuqWmLPgdjMkg56tLc7dQBrxBwG3qjXxVQQ%2Fll%2BCsANoiLH6WoHhBZUYz3S0SthOnYSU5E66oTUhpjOMnglKud4drTDYKR2Ljgvpjz0sAmS0Ynko9CVHRwbJIAefpbj4p0MKnVCI6IzmGazfNJxsvIO8EFtG7UbVjDf2tcKcm90oroU0W3tdoHdA6NE6JXM2AuN3M4vsMQfHCqLtjUJ%2F6xxIL3n4yTH6SulQex4IXXbMSzZC5z04c8ZLVrbcs8UhFTZhDIKNGono0GMi2kWNkUFLqdPY3%2FhdLOKeKv%2FTNjrqVW5FSWHQtdJ3LMrTAKfV9jmBRorwPKnFY%3D&X-Amz-SignedHeaders=host&X-Amz-Signature=246779b825d4f1bd9b410bb8413a48444a22c04f690b1a6dc2f99cf2de945c41", ExpiresAt:strfmt.DateTime{wall:0x2c99e860, ext:63773971304, loc:(*time.Location)(nil)}, FinishedAt:strfmt.DateTime{wall:0x2c99e860, ext:63773971304, loc:(*time.Location)(nil)}, ID:"a6c7a6d5-0fa5-4193-9647-ad70f13f34da", IntervalEnd:strfmt.DateTime{wall:0xa7d8c0, ext:63773971294, loc:(*time.Location)(nil)}, IntervalStart:strfmt.DateTime{wall:0xa7d8c0, ext:63773969494, loc:(*time.Location)(nil)}, Location:(*models.HashicorpCloudLocationLocation)(0xc00016a450), State:"READY"}

Download URL: "https://hcp-data-plane-blob-prod.s3.amazonaws.com/225af347-0fd9-41b1-8571-bed0d3ef665e/auditlogs/98b8bfa8-3115-4831-a598-434d02f83786?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=ASIAZDKPXWD4FHMAK4HM%2F20211201%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20211201T160154Z&X-Amz-Expires=900&X-Amz-Security-Token=FwoGZXIvYXdzEFEaDIKka9XFbRoODVjZSiLjAVJ4U5tMj8mqcuc9c2hzVpvCy6jj2TtLieIkyjTRbFOuqWmLPgdjMkg56tLc7dQBrxBwG3qjXxVQQ%2Fll%2BCsANoiLH6WoHhBZUYz3S0SthOnYSU5E66oTUhpjOMnglKud4drTDYKR2Ljgvpjz0sAmS0Ynko9CVHRwbJIAefpbj4p0MKnVCI6IzmGazfNJxsvIO8EFtG7UbVjDf2tcKcm90oroU0W3tdoHdA6NE6JXM2AuN3M4vsMQfHCqLtjUJ%2F6xxIL3n4yTH6SulQex4IXXbMSzZC5z04c8ZLVrbcs8UhFTZhDIKNGono0GMi2kWNkUFLqdPY3%2FhdLOKeKv%2FTNjrqVW5FSWHQtdJ3LMrTAKfV9jmBRorwPKnFY%3D&X-Amz-SignedHeaders=host&X-Amz-Signature=246779b825d4f1bd9b410bb8413a48444a22c04f690b1a6dc2f99cf2de945c41"
```

You can then download from that URL, gunzip it, then do whatever you like with your audit logs
