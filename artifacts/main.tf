resource "google_project_service" "kubernetes" {
  service = "container.googleapis.com"
  disable_on_destroy = false
}

resource "google_storage_bucket" "input" {
  name = "${var.project}-input"
}

resource "google_storage_bucket" "output" {
  name = "${var.project}-output"
}

resource "google_service_account" "beam" {
  account_id = "beamworker"
  display_name = "beamworker"
}

resource "google_storage_bucket_iam_member" "output" {
  bucket = google_storage_bucket.output.name
  member = "serviceAccount:${google_service_account.beam.email}"
  role = "roles/storage.objectAdmin"
}

resource "google_storage_bucket_iam_member" "input" {
  bucket = google_storage_bucket.input.name
  member = "serviceAccount:${google_service_account.beam.email}"
  role = "roles/storage.objectAdmin"
}

resource "google_storage_bucket_iam_member" "container" {
  bucket = "artifacts.${var.project}.appspot.com"
  member = "serviceAccount:${google_service_account.beam.email}"
  role = "roles/storage.objectViewer"
}

resource "google_container_cluster" "default" {
  name = "beam"
  location = "us-west1-a"
  depends_on = [google_project_service.kubernetes]
  initial_node_count = 1
  node_config {
    service_account = google_service_account.beam.email
    oauth_scopes = [
      "https://www.googleapis.com/auth/cloud-platform"
    ]
  }
}

