resource "google_project_service" "containerregistry" {
  service = "containerregistry.googleapis.com"
  disable_on_destroy = false
}

resource "google_project_service" "kubernetes" {
  service = "container.googleapis.com"
  disable_on_destroy = false
}

