terraform {
  backend "s3" {
    endpoints = {
      s3 = "https://storage.yandexcloud.net"
    }

    bucket = "hexlet-basics-terraform-state"
    region = "ru-central1"
    key    = "production_timeweb.tfstate"

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
  config_path = "../../.kube/config_twc"
}
