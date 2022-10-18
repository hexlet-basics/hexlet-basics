# --------------------------------------
# PRIMARY CLUSTER
# --------------------------------------

resource "digitalocean_kubernetes_cluster" "hexlet_basics_cluster_2" {
  name         = var.cluster_name_2
  region       = var.cluster_region

  auto_upgrade = true
  version      = data.digitalocean_kubernetes_versions.hexlet_basics_cluster_2.latest_version

  maintenance_policy {
    start_time  = "02:00"
    day         = "wednesday"
  }

  node_pool {
    name       = var.cluster_node_2_name
    size       = var.cluster_node_2_size
    node_count = 1
  }

  lifecycle {
    prevent_destroy = true
  }
}
resource "digitalocean_kubernetes_cluster" "hexlet_basics_cluster_3" {
  name         = var.cluster_name_3
  region       = var.cluster_region

  auto_upgrade = true
  version      = data.digitalocean_kubernetes_versions.hexlet_basics_cluster_2.latest_version

  maintenance_policy {
    start_time  = "02:00"
    day         = "wednesday"
  }

  node_pool {
    name       = var.cluster_node_3_name
    size       = var.cluster_node_3_size
    node_count = 3
  }

  lifecycle {
    prevent_destroy = true
  }
}

resource "digitalocean_kubernetes_node_pool" "hexlet_basics_nodes" {
  cluster_id = digitalocean_kubernetes_cluster.hexlet_basics_cluster_2.id
  name       = var.cluster_node_3_name
  size       = var.cluster_node_3_size
  node_count  = 3
}

locals {
  path_to_kubeconfig = "${path.root}/${var.rel_path_to_kubeconfig}"
}

# resource "local_file" "kubeconfig" {
#   depends_on = [
#     resource.digitalocean_kubernetes_cluster.hexlet_basics_cluster_2
#   ]

#   count      = var.write_kubeconfig ? 1 : 0
#   content    = resource.digitalocean_kubernetes_cluster.hexlet_basics_cluster_2.kube_config[0].raw_config
#   filename   = local.path_to_kubeconfig
#   file_permission = "0600"
# }

resource "local_file" "kubeconfig_2" {
  depends_on = [
    resource.digitalocean_kubernetes_cluster.hexlet_basics_cluster_3
  ]

  count      = var.write_kubeconfig ? 1 : 0
  content    = resource.digitalocean_kubernetes_cluster.hexlet_basics_cluster_3.kube_config[0].raw_config
  filename   = local.path_to_kubeconfig
  file_permission = "0600"
}

# --------------------------------------
# DATABASES
# --------------------------------------

# K8s Postges Databases
resource "digitalocean_database_cluster" "postgres_db_cluster" {
  name       = var.postgres_db_cluster_name
  engine     = "pg"
  version    = var.postgres_version
  size       = var.postgres_db_node_size
  region     = var.cluster_region
  node_count = 1

  lifecycle {
    prevent_destroy = true
  }
}

 resource "digitalocean_database_firewall" "postgres_db_firewall" {
   cluster_id = digitalocean_database_cluster.postgres_db_cluster.id

   rule {
     type  = "k8s"
     value = digitalocean_kubernetes_cluster.hexlet_basics_cluster_2.id
   }
   rule {
     type  = "k8s"
     value = digitalocean_kubernetes_cluster.hexlet_basics_cluster_3.id
   }
 }

# K8s Redis database
resource "digitalocean_database_cluster" "redis_db_cluster" {
  name       = var.redis_db_cluster_name
  engine     = "redis"
  version    = var.redis_version
  size       = var.redis_db_node_size
  region     = var.cluster_region
  node_count = 1

  lifecycle {
    prevent_destroy = true
  }
}

 resource "digitalocean_database_firewall" "redis_db_firewall" {
   cluster_id = digitalocean_database_cluster.redis_db_cluster.id

   rule {
     type  = "k8s"
     value = digitalocean_kubernetes_cluster.hexlet_basics_cluster_2.id
   }
   rule {
     type  = "k8s"
     value = digitalocean_kubernetes_cluster.hexlet_basics_cluster_3.id
   }
 }

resource "digitalocean_database_db" "postgres_db" {
  cluster_id = digitalocean_database_cluster.postgres_db_cluster.id
  name       = var.postgres_db_name
}

resource "digitalocean_database_user" "postgres_db_user" {
  cluster_id = digitalocean_database_cluster.postgres_db_cluster.id
  name       = var.postgres_db_user
}

# --------------------------------------
# MONITORING
# --------------------------------------
resource "digitalocean_monitor_alert" "disk_alert" {
  alerts {
    slack {
      url     = var.slack_notification_webhook
      channel = "#sideprojects-operation-auto"
    }
  }
  window      = "5m"
  type        = "v1/insights/droplet/disk_utilization_percent"
  compare     = "GreaterThan"
  value       = 80
  enabled     = true
  description = "Disk Utilization is running high"
}

resource "digitalocean_monitor_alert" "cpu_alert" {
  alerts {
    slack {
      url     = var.slack_notification_webhook
      channel = "#sideprojects-operation-auto"
    }
  }
  window      = "5m"
  type        = "v1/insights/droplet/cpu"
  compare     = "GreaterThan"
  value       = 70
  enabled     = true
  description = "CPU is running high"
}

resource "digitalocean_monitor_alert" "memory_alert" {
  alerts {
    slack {
      url     = var.slack_notification_webhook
      channel = "#sideprojects-operation-auto"
    }
  }
  window      = "5m"
  type        = "v1/insights/droplet/memory_utilization_percent"
  compare     = "GreaterThan"
  value       = 90
  enabled     = true
  description = "Memory utilization is running high"
}

# --------------------------------------
# SPACES
# --------------------------------------
resource "digitalocean_spaces_bucket" "sitemap_bucket" {
  name   = var.do_spaces_sitemap_bucket
  region = var.do_spaces_region
  acl    = "private"
  lifecycle_rule {
    enabled = true
    expiration {
      days = 30
    }
  }
}

# --------------------------------------
# PROJECT
# --------------------------------------

resource "digitalocean_project" "hexlet_basics_project" {
  name        = "Hexlet Basics"
  description = "A project to represent Hexlet Basics resources."
  purpose     = "Web Application"
  environment = "Production"
  resources   = [
    digitalocean_kubernetes_cluster.hexlet_basics_cluster_2.urn,
    digitalocean_database_cluster.postgres_db_cluster.urn,
    digitalocean_database_cluster.redis_db_cluster.urn,
  ]
}
