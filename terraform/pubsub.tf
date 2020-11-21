resource "google_pubsub_topic" "trigger_scraping" {
  name = "trigger_scraping"
}

resource "google_pubsub_topic" "trigger_graphing" {
  name = "trigger_graphing"
}