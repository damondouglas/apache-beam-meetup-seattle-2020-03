
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
