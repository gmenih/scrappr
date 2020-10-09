env = 

run:
	GOOGLE_APPLICATION_CREDENTIALS=./config/google-key.json go run local/main.go

build:
	zip upload_to_google.zip -r go.mod go.sum function.go scrape/