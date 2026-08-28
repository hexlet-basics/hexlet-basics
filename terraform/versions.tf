terraform {
  required_version = "~> 1.15.0"

  required_providers {
    yandex = {
      source  = "yandex-cloud/yandex"
    }

    cloudflare = {
      source  = "cloudflare/cloudflare"
    }

    # external = {
    #   source  = "hashicorp/external"
    #   version = "~> 2.4.0"
    # }
  }
}
