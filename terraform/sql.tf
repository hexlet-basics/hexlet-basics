resource "digitalocean_database_cluster" "hexlet_basics" {
  name = "master3"
  engine = "pg"
  version = "11"

  size = "db-s-1vcpu-1gb"
  region = "fra1"
  node_count = 1
}

resource "digitalocean_database_user" "hexlet_basics" {
  name     = "hexlet_basics"
  cluster_id = digitalocean_database_cluster.hexlet_basics.id
  # password = var.db_password
}

resource "digitalocean_database_db" "hexlet_basics_prod" {
  name     = var.db_name
  cluster_id = digitalocean_database_cluster.hexlet_basics.id
}

resource "digitalocean_database_db" "postgres" {
  name     = "postgres"
  cluster_id = digitalocean_database_cluster.hexlet_basics.id
}

resource "digitalocean_database_firewall" "k8s" {
  depends_on = [digitalocean_kubernetes_cluster.hexlet_basics_2]

  cluster_id = digitalocean_database_cluster.hexlet_basics.id
  rule {
    type  = "k8s"
    value = digitalocean_kubernetes_cluster.hexlet_basics_2.id
  }
}

resource "digitalocean_database_cluster" "hexlet-basics-redis" {
  name       = "hexlet-basics-redis"
  engine     = "redis"
  version    = "6"
  size       = "db-s-1vcpu-1gb"
  region     = "fra1"
  node_count = 1
}

# FIXME return with new k8s cluster
# resource "digitalocean_database_firewall" "redis-k8s" {
#   depends_on = [digitalocean_kubernetes_cluster.hexlet_basics_2]

#   cluster_id = digitalocean_database_cluster.hexlet-basics-redis.id
#   rule {
#     type  = "k8s"
#     value = digitalocean_kubernetes_cluster.hexlet_basics_2.id
#   }
# }
