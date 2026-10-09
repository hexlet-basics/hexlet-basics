terraform {
  required_version = ">= 1.11"

  required_providers {
    external = {
      source  = "hashicorp/external"
      version = "~> 2.4"
    }
    yandex = {
      source  = "yandex-cloud/yandex"
      version = "~> 0.219"
    }
  }
}
