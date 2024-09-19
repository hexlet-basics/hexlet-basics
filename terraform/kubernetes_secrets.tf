resource "kubernetes_secret" "datadog_secrets" {
  depends_on = [twc_k8s_cluster.hexlet_basics]

  metadata {
    name = "datadog-secret"
  }

  data = {
    api-key = var.datadog_api_key
  }
}

resource "kubernetes_secret" "environment_secrets" {
  depends_on = [twc_k8s_cluster.hexlet_basics]

  metadata {
    name = "hexlet-basics-environment-secrets"
  }

  data = var.environment_file
}
