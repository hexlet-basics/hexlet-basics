provider "google" {
  version     = "~> 2.7"
  region      = var.region
  credentials = file("google.key.json")
}

provider "cloudflare" {
  version = "~> 1.15.0"
  email    = var.cloudflare_email
  token  = var.cloudflare_api_key
}

provider "kubernetes" {
  version = "~> 1.7"

  # region = var.region
}

resource "google_project" "hexlet_basics" {
  name            = var.project_name
  project_id      = var.project_name
  billing_account = var.billing_account
  org_id          = var.org_id
}

provider "digitalocean" {
  version = "~> 1.16.0"
  token   =  "f8a961c5872bed9c15adb273efe7096109ad13708bd4b5d32ab4678756cc89ef"
}

output "project_id" {
  value = google_project.hexlet_basics.project_id
}
