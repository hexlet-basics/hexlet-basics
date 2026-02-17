terraform {
  required_version = "~> 1.14.0"

  required_providers {
    yandex = {
      source  = "yandex-cloud/yandex"
      version = "~> 0.187.0"
    }

    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "~> 5.17.0"
    }

    external = {
      source  = "hashicorp/external"
      version = "~> 2.3.0"
    }
  }
}
