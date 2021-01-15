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
  token   =  "60e6ec66e21557fce9f7033856abc30a54d148099d5f0f8a3ecdbf0f5927d6ed"
}

output "project_id" {
  value = google_project.hexlet_basics.project_id
}
