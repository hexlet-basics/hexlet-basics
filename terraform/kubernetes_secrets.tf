resource "kubernetes_secret" "grafana_alloy_secrets_3" {
  depends_on = [twc_k8s_cluster.hexlet_basics_3]

  metadata {
    name = "grafana-alloy-secrets"
  }

  data = {
    receiver_url = var.victoriametrics.url
    receiver_username = var.victoriametrics.username
    receiver_password = var.victoriametrics.password
    db_host = var.postgres_db.host
    db_port = var.postgres_db.port
    db_name = var.postgres_db.name
    db_user = var.postgres_db.datadog_username
    db_password = var.postgres_db.datadog_password
  }
}

resource "kubernetes_secret" "environment_secrets_3" {
  depends_on = [twc_k8s_cluster.hexlet_basics_3]

  metadata {
    name = "hexlet-basics-environment-secrets"
  }

  data = var.environment_file
}

resource "kubernetes_secret" "docker_registry_auth_3" {
  depends_on = [twc_k8s_cluster.hexlet_basics_3]

  metadata {
    name = "docker-config"
  }

  data = {
    ".dockerconfigjson" = file("${path.module}/docker-config.json")
  }

  type = "kubernetes.io/dockerconfigjson"
}
