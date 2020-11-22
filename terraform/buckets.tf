resource "google_storage_bucket" "function_source" {
  name = "${var.gcp_project}_function_source"
}

resource "google_storage_bucket" "graph_images" {
  name = "${var.gcp_project}_graph_images"
}
