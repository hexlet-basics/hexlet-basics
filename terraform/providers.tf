terraform {
  required_version = ">=1.0.0"

  required_providers {
    digitalocean = {
      source = "digitalocean/digitalocean"
      version = "~> 2.0"
    }
    kubernetes = {
      source = "hashicorp/kubernetes"
      version = ">= 2.0.0"
    }
    cloudflare = {
      source = "cloudflare/cloudflare"
      version = "~> 3.0"
    }
  }

  backend "gcs" {
    bucket      = "hexlet-basics-terraform-state-2"
    prefix      = "hb/tfstate"
    credentials = "./google.key.json"
  }
}

provider "digitalocean" {
  spaces_access_id  = var.do_spaces_access_id
  spaces_secret_key = var.do_spaces_secret_key
}

# NOTE При создании нового кластера данных в data.digitalocean_kubernetes_cluster.hexlet_basics_cluster_data еще не будет
provider "kubernetes" {
  # host  = data.digitalocean_kubernetes_cluster.hexlet_basics_cluster_data_3.endpoint
  # token = data.digitalocean_kubernetes_cluster.hexlet_basics_cluster_data_3.kube_config[0].token
  # cluster_ca_certificate = base64decode(
  #   data.digitalocean_kubernetes_cluster.hexlet_basics_cluster_data_3.kube_config[0].cluster_ca_certificate
  # )

  config_path = "../.kube/config"
}

provider "cloudflare" {
  email   = var.cloudflare_email
  api_key = var.cloudflare_api_key
}
