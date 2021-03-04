resource "digitalocean_kubernetes_cluster" "hexlet_basics_2" {
  version = "1.18.8-do.0"

  name = "hexlet-basics-2"
  region = "fra1"

  node_pool {
    name       = "hexlet-basics-node-pool-2"
    size       = "c-2"
    node_count = 2
  }
}

resource "digitalocean_kubernetes_cluster" "hexlet-basics-4" {
  version = "1.20.2-do.0"

  name = "hexlet-basics-4"
  region = "fra1"

  auto_upgrade = true

  node_pool {
    name       = "hexlet-basics-node-pool-4"
    size       = "c-2"
    node_count = 3
  }
}
