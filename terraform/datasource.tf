data "digitalocean_kubernetes_cluster" "hexlet_basics_cluster_data" {
  name = var.cluster_name

  depends_on = [
    digitalocean_kubernetes_cluster.hexlet_basics_cluster
  ]
}

data "digitalocean_database_cluster" "postgres_db_data" {
  name = var.postgres_db_cluster_name

  depends_on = [
    digitalocean_database_cluster.postgres_db_cluster
  ]
}

data "digitalocean_database_cluster" "redis_db_data" {
  name = var.redis_db_cluster_name

  depends_on = [
    digitalocean_database_cluster.redis_db_cluster
  ]
}
