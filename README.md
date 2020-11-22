# scrappr

Nepremicnine scarper orchestrated by Terraform, running on GCP.

## Setting it up

### GCP Console

Create a GCP project, and generate a service account with the following roles:
* Cloud Functions Admin
* Cloud Functions Developer
* Cloud Scheduler Admin
* Service Account User
* Pub/Sub Admin
* Storage Admin
* Storage Object Admin

Generate a new key for the service account, and store it in `./config/terraform-key.json`

This can be achieved by running (probably, untested):

```shell
export PROJECT_ID="your-project-id"

gcloud iam service-acount create terraform-user \
    --description="Terraform service account" \
    --display-name="Terraform User"

gcloud project add-iam-policy-binding $PROJECT_ID \
    --member="serviceAccount:terraform-user@$PROJECT_ID.iam.gserviceaccount.com" \
    --role="roles/cloudfunctions.admin" \
    --role="roles/cloudfunctions.developer" \
    --role="roles/cloudscheduler.admin" \
    --role="roles/iam.serviceAccountUser" \
    --role="roles/pubsub.admin" \
    --role="roles/storage.admin" \
    --role="roles/storage.objectAdmin"

```


### Terraform
Open the `Makefile` and edit the `gcp_project` variable (and any others that you think are incorrect)

Then set up Terraform via the `make` command:
```shell
# build the zip archive which will run in Google Cloud Functions
make build

# initialize terraform providers
make tf-init

# plan terraform resources
make tf-plan

# apply
make tf-apply
```

## Running it locally

`make run-scraper`