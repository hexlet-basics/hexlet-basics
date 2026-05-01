terraform {
  required_version = "~> 1.15.0"

  required_providers {
    yandex = {
      source  = "yandex-cloud/yandex"
      version = "~> 0.201.0"
    }

    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "~> 5.19.0"
    }

    external = {
      source  = "hashicorp/external"
      version = "~> 2.3.0"
    }
  }
}
