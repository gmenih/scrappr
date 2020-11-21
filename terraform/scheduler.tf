resource "google_cloud_scheduler_job" "scrape_realestate_mb_apt" {
  name        = "scrape_realestate_mb_apt"
  description = "Scrape realestate"
  schedule    = "*/30 6-22 * * *"

  pubsub_target {
    # topic.id is the topic's full resource name.
    topic_name = google_pubsub_topic.trigger_scraping.id
    data = base64encode(jsonencode({
      "url"  = "https://www.nepremicnine.net/oglasi-prodaja/podravska/stanovanje/"
      "type" = "apartment"
    }))
  }
}

resource "google_cloud_scheduler_job" "scrape_realestate_mb_house" {
  name        = "scrape_realestate_mb_house"
  description = "Scrape realestate"
  schedule    = "*/35 6-22 * * *"

  pubsub_target {
    # topic.id is the topic's full resource name.
    topic_name = google_pubsub_topic.trigger_scraping.id
    data = base64encode(jsonencode({
      "url"  = "https://www.nepremicnine.net/oglasi-prodaja/podravska/hisa/"
      "type" = "house"
    }))
  }
}
