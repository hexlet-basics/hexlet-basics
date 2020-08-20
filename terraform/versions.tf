terraform {
  required_providers {
    digitalocean = {
      source = "terraform-providers/digitalocean"
    }
    google = {
      source = "hashicorp/google"
    }
  }
  required_version = ">= 0.13"
}
