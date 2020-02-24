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

resource "google_project_iam_member" "bigquery" {
  member = "serviceAccount:${google_service_account.beam.email}"
  role = "roles/bigquery.dataEditor"
}