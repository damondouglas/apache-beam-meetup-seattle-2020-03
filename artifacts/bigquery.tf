resource "google_bigquery_dataset" "beam" {
  dataset_id = "beam"
}

resource "google_bigquery_table" "patients" {
  dataset_id = google_bigquery_dataset.beam.dataset_id
  table_id = "patients"
  schema = <<EOF
[
    {
        "name": "mrn",
        "type": "STRING",
        "mode": "REQUIRED"
    },
    {
        "name": "allergies",
        "type": "STRING",
        "mode": "REQUIRED"
    }
]
EOF
}