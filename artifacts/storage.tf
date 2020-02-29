resource "google_storage_bucket" "input" {
  name = "${var.project}-input"
}

resource "google_storage_bucket" "output" {
  name = "${var.project}-output"
}

