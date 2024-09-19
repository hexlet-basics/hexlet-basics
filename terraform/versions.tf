terraform {
  required_version = ">= 0.13"

  required_providers {
    twc = {
      source = "tf.timeweb.cloud/timeweb-cloud/timeweb-cloud"
    }

    digitalocean = {
      source = "digitalocean/digitalocean"
    }

    cloudflare = {
      source = "cloudflare/cloudflare"
    }
  }
}
