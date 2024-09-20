terraform {
  required_version = ">= 0.13"

  required_providers {
    twc = {
      source = "tf.timeweb.cloud/timeweb-cloud/timeweb-cloud"
    }

    cloudflare = {
      source = "cloudflare/cloudflare"
    }

    yandex = {
      source = "yandex-cloud/yandex"
    }
  }
}
