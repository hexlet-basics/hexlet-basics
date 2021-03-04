provider "cloudflare" {
  email    = var.cloudflare_email
  token  = var.cloudflare_api_key
}

provider "kubernetes" {
}

resource "google_project" "hexlet_basics" {
  name            = var.project_name
  project_id      = var.project_name
  billing_account = var.billing_account
  org_id          = var.org_id
}

provider "digitalocean" {
  token   =  var.digitalocean_token
}

output "project_id" {
  value = google_project.hexlet_basics.project_id
}
