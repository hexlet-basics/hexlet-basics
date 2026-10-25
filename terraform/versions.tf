terraform {
  required_version = "~> 1.17.0"

  required_providers {
    yandex = {
      source  = "yandex-cloud/yandex"
      version = "~> 0.217.0"
    }

    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "~> 5.28.0"
    }

    external = {
      source  = "hashicorp/external"
      version = "~> 2.4.0"
    }
  }
}
