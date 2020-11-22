resource "random_id" "function_id" {
  byte_length = 4  
}

resource "google_storage_bucket_object" "function_source" {
  name   = "index_${random_id.function_id.hex}.zip"
  bucket = google_storage_bucket.function_source.name
  source = "../dist/index.zip"
}

resource "google_cloudfunctions_function" "scrape_function" {
  name        = "scrape_function"
  description = "Scrape apartments"

  source_archive_bucket = google_storage_bucket.function_source.name
  source_archive_object = google_storage_bucket_object.function_source.name

  runtime             = "go113"
  entry_point         = "Scrape"
  available_memory_mb = 128
  timeout             = 540

  event_trigger {
    event_type = "google.pubsub.topic.publish"
    resource   = google_pubsub_topic.trigger_scraping.id
  }

  environment_variables = {
    "GCP_PROJECT" = var.gcp_project
  }
}

resource "google_cloudfunctions_function" "draw_graphs" {
  name        = "draw_graphs"
  description = "Draw graphs"

  source_archive_bucket = google_storage_bucket.function_source.name
  source_archive_object = google_storage_bucket_object.function_source.name

  runtime             = "go113"
  entry_point         = "DrawGraphs"
  available_memory_mb = 256
  timeout             = 540
  trigger_http        = true

  environment_variables = {
    "GCP_PROJECT" = var.gcp_project
  }
}
