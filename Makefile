env = GCP_PROJECT="your-project-id" \
	GOOGLE_APPLICATION_CREDENTIALS="./config/google-key.json"

run:
	$(env) go run local/main.go

build:
	zip upload_to_google.zip -r go.mod go.sum function.go scrape/