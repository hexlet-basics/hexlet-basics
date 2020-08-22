resource "google_sql_database_instance" "master" {
  project          = var.project_name
  name             = "master3"
  database_version = "POSTGRES_9_6"
  region           = var.region

  settings {
    tier              = "db-f1-micro"
    availability_type = "REGIONAL"

    backup_configuration {
      enabled = true
    }
  }
}

resource "google_sql_user" "hexlet_basics" {
  project  = var.project_name
  name     = "hexlet_basics"
  instance = google_sql_database_instance.master.name
  password = var.db_password
}

resource "google_sql_database" "hexlet_basics_prod" {
  project  = var.project_name
  name     = var.db_name
  instance = google_sql_database_instance.master.name
}

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
