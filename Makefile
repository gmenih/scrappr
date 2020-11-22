env = \
	GOOGLE_CLOUD_PROJECT="my-project" \
	GOOGLE_APPLICATION_CREDENTIALS="./config/google-key.json" \
	SPREADSHEET_ID="your-spreadsheet-id"

tf_env = \
	TF_VAR_gcp_project="my-project" \
	TF_VAR_gcp_region="europe-west1" \
	TF_VAR_gcp_config_file="$(realpath .)/config/terraform-key.json"

SHELL := /bin/zsh

run:
	$(env) go run src/local/$(file).go

run-scraper:
	make run file=scrape

run-grapher:
	make run file=graph

clean:
	rm -r dist
	mkdir -p dist

build: clean
	mkdir -p dist/tmp
	cp go.{sum,mod} ./dist/tmp
	cp -r ./src ./dist/tmp
	mv ./dist/tmp/src/*.go ./dist/tmp
	@cd ./dist/tmp && \
		zip ../index.zip -r ./*
	rm -r ./dist/tmp

tf-plan:
	@echo "Generating TF plan"
	cd terraform && \
		$(tf_env) terraform plan -out plan.tfplan
	@cd -

tf-apply:
	cd terraform && \
		terraform apply ./plan.tfplan

