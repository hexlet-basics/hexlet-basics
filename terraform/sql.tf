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
}

resource "digitalocean_database_db" "hexlet_basics_prod" {
  name     = var.db_name
  cluster_id = digitalocean_database_cluster.hexlet_basics.id
}

resource "digitalocean_database_db" "postgres" {
  name     = "postgres"
  cluster_id = digitalocean_database_cluster.hexlet_basics.id
}
