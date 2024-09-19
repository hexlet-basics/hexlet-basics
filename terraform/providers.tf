terraform {
  required_version = ">=1.0.0"

  required_providers {
    digitalocean = {
      source = "digitalocean/digitalocean"
      version = "~> 2.0"
    }
    cloudflare = {
      source = "cloudflare/cloudflare"
      version = "~> 3.0"
    }
  }

  backend "gcs" {
    bucket      = "hexlet-basics-terraform-state-2"
    prefix      = "hb/tfstate"
    credentials = "./google.key.json"
  }
}

provider "digitalocean" {
  spaces_access_id  = var.do_spaces_access_id
  spaces_secret_key = var.do_spaces_secret_key
}

provider "cloudflare" {
  email   = var.cloudflare_email
  api_key = var.cloudflare_api_key
}
