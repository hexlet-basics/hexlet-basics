terraform {
  required_providers {
    cloudflare = {
      source = "terraform-providers/cloudflare"
    }
    digitalocean = {
      source = "terraform-providers/digitalocean"
    }
    google = {
      source = "hashicorp/google"
    }
    kubernetes = {
      source = "hashicorp/kubernetes"
    }
  }
  required_version = ">= 0.13"
}
