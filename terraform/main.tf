terraform {
  backend "s3" {
    endpoints = {
      s3 = "https://storage.yandexcloud.net"
    }

    bucket = "code-basics-terraform-state"
    region = "ru-central1"
    key    = "production_code_basics.tfstate"

    skip_region_validation      = true
    skip_credentials_validation = true
    skip_requesting_account_id  = true
    skip_s3_checksum            = true
  }
}

data "external" "helm-secrets" {
  program = ["helm", "secrets", "decrypt", "--terraform", "../k8s/secrets.yaml"]
}

locals {
  data = yamldecode(base64decode(data.external.helm-secrets.result.content_base64))
}

provider "cloudflare" {
  email   = local.data.terraform.cloudflare.email
  api_key = local.data.terraform.cloudflare.api_key
}

provider "yandex" {
  cloud_id                 = local.data.terraform.yc.cloud_id
  folder_id                = local.data.terraform.yc.folder_id
  zone                     = local.data.terraform.yc.zone
  service_account_key_file = "yc_config.json"
}
