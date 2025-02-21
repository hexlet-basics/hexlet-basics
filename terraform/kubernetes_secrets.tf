resource "kubernetes_secret" "datadog_secrets" {
  depends_on = [twc_k8s_cluster.hexlet_basics]

  metadata {
    name = "datadog-secret"
  }

  data = {
    api-key = var.datadog_api_key
    db_host = var.postgres_db.host
    db_port = var.postgres_db.port
    db_name = var.postgres_db.name
    db_user = var.postgres_db.datadog_username
    db_password = var.postgres_db.datadog_password
  }
}

resource "kubernetes_secret" "environment_secrets" {
  depends_on = [twc_k8s_cluster.hexlet_basics]

  metadata {
    name = "hexlet-basics-environment-secrets"
  }

  data = var.environment_file
}

resource "kubernetes_secret" "docker_registry_auth" {
  depends_on = [twc_k8s_cluster.hexlet_basics]

  metadata {
    name = "docker-config"
  }

  data = {
    ".dockerconfigjson" = file("${path.module}/docker-config.json")
  }

  type = "kubernetes.io/dockerconfigjson"
}
