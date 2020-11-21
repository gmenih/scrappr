env = \
	GOOGLE_CLOUD_PROJECT="your-project" \
	GOOGLE_APPLICATION_CREDENTIALS="./config/google-key.json" \
	SPREADSHEET_ID="your-spreadsheet-id"

tf_env = \
	TF_VAR_gcp_project="your-project" \
	TF_VAR_gcp_region="europe-west3-c"

SHELL := /bin/zsh

run:
	$(env) go run src/local/$(file).go

run-scraper:
	make run file=scrape

run-grapher:
	make run file=graph

build:
	zip upload_to_google.zip -r go.mod go.sum ./{functions,realestate,scrape}

tf-plan:
	@echo "Generating TF plan"
	cd terraform && \
		$(tf_env) terraform plan -out plan.tfplan
	@cd -

tf-apply:
	cd terraform && \
	@terraform apply plan.tfplan \
	&& cd -

