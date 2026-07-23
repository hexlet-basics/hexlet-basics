terraform {
  required_version = "~> 1.15.0"

  required_providers {
    yandex = {
      source  = "yandex-cloud/yandex"
      version = "~> 0.219.0"
    }

    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "~> 5.22.0"
    }

    external = {
      source  = "hashicorp/external"
      version = "~> 2.4.0"
    }
  }
}
