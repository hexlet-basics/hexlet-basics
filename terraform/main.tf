terraform {
  backend "s3" {
    endpoints = {
      s3 = "https://storage.yandexcloud.net"
    }

    bucket = "hexlet-basics-terraform-state"
    region = "ru-central1"
    key    = "production_hexlet_basics.tfstate"

    skip_region_validation      = true
    skip_credentials_validation = true
    skip_requesting_account_id  = true
    skip_s3_checksum            = true
  }
}

# NOTE: needed TWC_TOKEN env variable
provider "twc" {}

resource "twc_project" "hexlet_basics" {
  name        = "Hexlet Basics"
  description = "Hexlet Basics infrastructure"
}

provider "kubernetes" {
  config_path = "../.kube/config"
}

provider "cloudflare" {
  email   = var.cloudflare_email
  api_key = var.cloudflare_api_key
}

provider "yandex" {
  cloud_id                 = var.yc.cloud_id
  folder_id                = var.yc.folder_id
  zone                     = var.yc.zone
  service_account_key_file = "yc_config.json"
}
