terraform {
  required_version = ">= 0.13"

  required_providers {
    yandex = {
      source = "yandex-cloud/yandex"
    }

    cloudflare = {
      source = "cloudflare/cloudflare"
    }
  }
}
