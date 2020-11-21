resource "google_storage_bucket" "function_source" {
  name = "function_source"
}

resource "google_storage_bucket" "graph_images" {
  name = "graph_images"
}
