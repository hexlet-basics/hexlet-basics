terraform {
  backend "gcs" {
    bucket      = "hexlet-basics-terraform-state"
    prefix      = "production"
    credentials = "google.key.json"
  }
}
