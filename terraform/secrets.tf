resource "kubernetes_secret" "sparkpost_credentials" {
  metadata {
    name = "sparkpost-credentials"
  }

  data = {
    SPARKPOST_SMTP_USERNAME = var.sparkpost_smtp_username
    SPARKPOST_SMTP_PASSWORD = var.sparkpost_smtp_password
    GUARDIAN_SECRET_KEY = var.guardian_secret_key
  }
}

resource "kubernetes_secret" "github_credentials" {
  metadata {
    name = "github-credentials"
  }

  data = {
    GITHUB_CLIENT_ID     = var.github_client_id
    GITHUB_CLIENT_SECRET = var.github_client_secret
  }
}

resource "kubernetes_secret" "facebook_credentials" {
  metadata {
    name = "facebook-credentials"
  }

  data = {
    FACEBOOK_CLIENT_ID     = var.facebook_client_id
    FACEBOOK_CLIENT_SECRET = var.facebook_client_secret
  }
}

locals {
  postgres_db_user_pass = resource.digitalocean_database_user.postgres_db_user.password
  database_url = "postgres://${var.postgres_db_user}:${local.postgres_db_user_pass}@${data.digitalocean_database_cluster.postgres_db_data.host}:${data.digitalocean_database_cluster.postgres_db_data.port}/${var.postgres_db_name}"
}

resource "kubernetes_secret" "database_credentials" {
  metadata {
    name = "database-credentials"
  }

  data = {
    DB_HOSTNAME  = data.digitalocean_database_cluster.postgres_db_data.host
    DB_NAME      = var.postgres_db_name
    DB_USERNAME  = var.postgres_db_user
    DB_PORT  = data.digitalocean_database_cluster.postgres_db_data.port
    DB_PASSWORD  = data.digitalocean_database_cluster.postgres_db_data.password
    DATABASE_URL = local.database_url
    REDIS_URL    = data.digitalocean_database_cluster.redis_db_data.uri
  }
}

resource "kubernetes_secret" "hexlet_basics_secrets" {
  metadata {
    name = "hexlet-basics-secrets"
  }

  data = {
    SECRET_KEY_BASE      = var.secret_key_base
    ROLLBAR_ACCESS_TOKEN = var.rollbar_access_token
    ROLLBAR_CLIENT_ACCESS_TOKEN = var.rollbar_client_access_token
  }
}

resource "kubernetes_config_map" "hexlet_basics_data" {
  metadata {
    name = "hexlet-basics-data"
  }

  data = {
    MIX_ENV      = "prod"
    PORT         = "3000"
    NODE_ENV     = "production"
    RAILS_ENV    = "production"
    DB_POOL_SIZE = "10"
    DB_SSL_MODE  = "TRUE"
    RAILS_LOG_TO_STDOUT = true
    RAILS_SERVE_STATIC_FILES = true
    FORCE        = "11"
    APP_SCHEME   = var.app_scheme
    APP_HOST     = var.app_host
    APP_RU_HOST  = var.app_ru_host
    GOOGLE_TAG_MANAGER_KEY = var.google_tag_manager_key
  }

  depends_on = [resource.local_file.kubeconfig]
}

resource "kubernetes_secret" "docker-registry-auth" {
  metadata {
    name = "docker-config"
  }

  data = {
    ".dockerconfigjson" = file("${path.module}/docker-config.json")
  }

  type = "kubernetes.io/dockerconfigjson"
}
