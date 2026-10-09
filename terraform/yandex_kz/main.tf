terraform {
  backend "s3" {
    endpoints = {
      s3 = "https://storage.yandexcloud.kz"
    }

    bucket = "hexlet-basics-terraform-state"
    region = "kz1"
    key    = "production_code_basics_kz.tfstate"

    skip_region_validation      = true
    skip_credentials_validation = true
    skip_requesting_account_id  = true
    skip_s3_checksum            = true
  }
}

data "external" "helm-secrets" {
  program = ["helm", "secrets", "decrypt", "--terraform", "../../k8s/secrets.yaml"]
}

locals {
  data = yamldecode(base64decode(data.external.helm-secrets.result.content_base64))
}

provider "yandex" {
  endpoint         = "api.yandexcloud.kz:443"
  storage_endpoint = "storage.yandexcloud.kz"

  cloud_id  = local.data.terraform.yc_kz.cloud_id
  folder_id = local.data.terraform.yc_kz.folder_id
  zone      = local.data.terraform.yc_kz.zone

  service_account_key_file = "yc_config.json"
  storage_access_key       = var.access_key
  storage_secret_key       = var.secret_key
}
