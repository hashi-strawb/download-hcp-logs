# Example Go Script to pull HCP Vault Audit Logs

WARNING: This makes use of unstable preview APIs which could change at any time!

This is also meant as an example of how pulling logs from the HCP Vault APIs could work.
You probably do not want to use this directly. Instead, fork the repo, and customise it to your requirements.

USE AT YOUR OWN PERIL


## Why?

HCP Vault has audit logs! Woo!

HCP Vault only supports streaming those to [a small number of places](https://www.hashicorp.com/blog/hcp-vault-adds-3-new-observability-integrations)! Boo!

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

(you can also put these environment variables in a `.env` file, if you omit the `export`)

## Running

```
go run .
```

Example output:

```
$ go run main.go | head -n 2 | jq .
INFO[0000] Response: &{0xc0000be1e0 2021-12-01T14:49:43.342Z v1.8.5 0xc00007ae20 vault-cluster 0xc000250d80 RUNNING}
INFO[0001] Log ID: 639d8cc4-52d6-4bd8-ba86-9aef948e8348
INFO[0001] State: PENDING
INFO[0011] State: CREATING
INFO[0021] State: READY
INFO[0021] Log Payload: &{vault-cluster https://hcp-data-plane-blob-prod.s3.amazonaws.com/REDACTED 2021-12-02T10:12:03.669Z 2021-12-02T10:12:03.669Z 639d8cc4-52d6-4bd8-ba86-9aef948e8348 2021-12-02T10:11:50.905Z 2021-12-02T09:41:50.905Z 0xc0001022d0 READY}
INFO[0021] Download URL: https://hcp-data-plane-blob-prod.s3.amazonaws.com/REDACTED

{
  "time": "2021-12-02T09:01:09.745135603Z",
  "type": "request",
  "auth": {
    "client_token": "hmac-sha256:d443d8ce8764c852c2b4b202f7c2d4a70d575c3353b5ad1759c5629bdcbdea55",
    "accessor": "hmac-sha256:78313a654c3363a9301496d5440e44468c039f84b1bc285db050f9ab6722d1f0",
    "display_name": "jwt-msp",
    "policies": [
      "default",
      "msp"
    ],
    "token_policies": [
      "default",
      "msp"
    ],
    "metadata": {
      "role": "msp"
    },
    "entity_id": "c4e52b80-f836-7453-734d-1c0ef6c3ab23",
    "token_type": "service",
    "token_ttl": 600,
    "token_issue_time": "2021-12-02T09:01:08Z"
  },
  "request": {
    "id": "cae52ccd-5564-4a7a-cd4d-f8cc025bca17",
    "operation": "read",
    "mount_type": "kv",
    "client_token": "hmac-sha256:d443d8ce8764c852c2b4b202f7c2d4a70d575c3353b5ad1759c5629bdcbdea55",
    "client_token_accessor": "hmac-sha256:78313a654c3363a9301496d5440e44468c039f84b1bc285db050f9ab6722d1f0",
    "namespace": {
      "id": "root"
    },
    "path": "hcp-metadata/health",
    "remote_address": "172.25.20.18"
  }
}
{
  "time": "2021-12-02T09:01:09.745668159Z",
  "type": "response",
  "auth": {
    "client_token": "hmac-sha256:d443d8ce8764c852c2b4b202f7c2d4a70d575c3353b5ad1759c5629bdcbdea55",
    "accessor": "hmac-sha256:78313a654c3363a9301496d5440e44468c039f84b1bc285db050f9ab6722d1f0",
    "display_name": "jwt-msp",
    "policies": [
      "default",
      "msp"
    ],
    "token_policies": [
      "default",
      "msp"
    ],
    "metadata": {
      "role": "msp"
    },
    "entity_id": "c4e52b80-f836-7453-734d-1c0ef6c3ab23",
    "token_type": "service",
    "token_ttl": 600,
    "token_issue_time": "2021-12-02T09:01:08Z"
  },
  "request": {
    "id": "cae52ccd-5564-4a7a-cd4d-f8cc025bca17",
    "operation": "read",
    "mount_type": "kv",
    "client_token": "hmac-sha256:d443d8ce8764c852c2b4b202f7c2d4a70d575c3353b5ad1759c5629bdcbdea55",
    "client_token_accessor": "hmac-sha256:78313a654c3363a9301496d5440e44468c039f84b1bc285db050f9ab6722d1f0",
    "namespace": {
      "id": "root"
    },
    "path": "hcp-metadata/health",
    "remote_address": "172.25.20.18"
  },
  "response": {
    "mount_type": "kv",
    "secret": {},
    "data": {
      "data": {
        "write_timestamp": "hmac-sha256:98c697f37af07149b690257e12c6bcf14b5ffc1cb51556f5d583792a2c3ef5cf"
      }
    }
  }
}
```

## Known Limitations

The API will not respect log intervals specfified with a duration of less than 1 hour.

As the API is not yet stable, the example code provided in this repository will not attempt to compensate.

This can be demonstrated by running:

```
go run . | jq .time
```
Examples of this behaviour is presented below.

### Logs downloaded with greater time interval than requested

For example:
```
INFO[1233] Downloading logs from 2021-12-03T16:38:39.656Z to 2021-12-03T16:39:50.989Z
```

* Earliest logs: "2021-12-03T16:02:44.108136091Z"
* Latest logs: "2021-12-03T16:39:03.208386076Z"


### Log download lag

Additionally, the most recent logs downloaded may be a few minutes behind what was requested.

For example:
```
INFO[1457] Downloading logs from 2021-12-03T16:42:13.567Z to 2021-12-03T16:43:35.028Z
```

Latest logs: "2021-12-03T16:40:47.134782281Z"

### Logs around the top of the hour

Also, logs around the top of the hour do not behave as you may expect.

For example:

```
INFO[0071] Downloading logs from 2021-12-03T17:03:49.392Z to 2021-12-03T17:05:00.436Z
```

Earliest logs: "2021-12-03T17:02:50.099642805Z"

and
```
INFO[0000] Downloading logs from 2021-12-03T16:55:52.026Z to 2021-12-03T17:05:52.026Z
```

Latest logs: "2021-12-03T17:04:44.576945466Z"






## TODO

* [X] Download logs in a loop
* [X] Support for .env for env vars
* [ ] Refactor the whole thing
  * [ ] Use Location from `vault.NewListParams()` by cluster-id rather than hard-coding


Out-of-scope (until a stable API is available)
* [ ] Use an external progress file to keep track of the last timestamp which was pulled
* [ ] Optionally output logs to a file
* [ ] Dedupe logs from download
* [ ] Build a binary and release
* [ ] Sample Vault Agent config in README to populate .env
