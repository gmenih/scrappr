resource "google_datastore_index" "id_createdAt" {
  kind = "price"
  properties {
    name = "id"
    direction = "ASCENDING"
  }
  properties {
    name = "createdAt"
    direction = "ASCENDING"
  }
}

resource "google_datastore_index" "id_price" {
  kind = "price"
  properties {
    name = "id"
    direction = "ASCENDING"
  }
  properties {
    name = "price"
    direction = "ASCENDING"
  }
}
