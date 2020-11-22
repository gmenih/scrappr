variable "scheduled_scrapes" {
  default = {
    "podravska_apt" = {
      url  = "https://www.nepremicnine.net/oglasi-prodaja/podravska/stanovanje/"
      type = "apartment"
      schedule = "0,30 6-22 * * *"
    }
    "podravska_house" = {
      url  = "https://www.nepremicnine.net/oglasi-prodaja/podravska/hisa/"
      type = "house"
       schedule = "5,35 6-22 * * *"
    }
  }
}

resource "google_cloud_scheduler_job" "scrape_realestate" {
  for_each = var.scheduled_scrapes

  name        = "scrape_realestate_${each.key}"
  description = "Scrape realestate ${each.key}"
  schedule    = each.value.schedule != null ? each.value.schedule : "0 * * * *"

  pubsub_target {
    # topic.id is the topic's full resource name.
    topic_name = google_pubsub_topic.trigger_scraping.id
    data = base64encode(jsonencode({
      "url"  = "${each.value.url}"
      "type" = "${each.value.type}"
    }))
  }
}
