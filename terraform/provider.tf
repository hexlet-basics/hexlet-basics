provider "digitalocean" {
  version = "~> 1.22.2"
  token = var.digitalocean_token
}

provider "google" {
  version = "~> 3.35.0"
  region = var.region
  credentials = file("google.key.json")
}
