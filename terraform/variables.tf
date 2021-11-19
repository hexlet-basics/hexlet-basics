# Define K8s Region
variable "cluster_region" {
  type = string
  description = "K8s Cluster Region"
}
variable "cluster_name" {
  type = string
  description = "K8s Cluster Name"
}

variable "cluster_node_name" {
  type = string
  description = "K8s Cluster Node Name"
}

variable "cluster_node_size" {
  type = string
  description = "K8s Node Size"
}

variable "postgres_db_node_size" {
  type = string
  description = "K8s DB Size"
}

variable "postgres_version" {
  type = string
  description = "Postgres version"
}

variable "postgres_db_cluster_name" {
  type = string
  description = "Postgres DB Cluster name"
}

variable "postgres_db_name" {
  type = string
  description = "Hexlet Basics database name"
}

variable "postgres_db_user" {
  type = string
  description = "Hexlet Basics database user"
}

variable "redis_db_cluster_name" {
  type = string
  description = "Postgres DB name"
}

variable "redis_version" {
  type = string
  description = "Redis version"
}

variable "redis_db_node_size" {
  type = string
  description = "Redis db version"
}

variable "write_kubeconfig" {
  type        = bool
  default     = true
}

variable "rel_path_to_kubeconfig" {
  type        = string
  default     = "../.kube/config"
  description = "Path to kubeconfig file"
}


# Credentials
variable "sparkpost_smtp_username" {
  type        = string
  description = "Sparkpost SMTP Username"
}

variable "sparkpost_smtp_password" {
  type        = string
  description = "Sparkpost SMTP Password"
}
variable "guardian_secret_key" {
  type        = string
  description = "Sparkpost Guardian secret key"
}

variable "github_client_id" {
  type        = string
  description = "Github client id"
}

variable "github_client_secret" {
  type        = string
  description = "Github client secret"
}

variable "facebook_client_id" {
  type        = string
  description = "Facebook client ID"
}

variable "facebook_client_secret" {
  type        = string
  description = "Facebook client secret"
}

variable "rollbar_access_token" {
  type        = string
  description = "Rollbar access token"
}

variable "rollbar_client_access_token" {
  type        = string
  description = "Rollbar client access token"
}

variable "secret_key_base" {
  type        = string
  description = "Hexlet Basics secret key base"
}

variable "app_scheme" {
  type        = string
  description = "Hexlet Basics app scheme"
}

variable "app_host" {
  type        = string
  description = "Hexlet Basics app host"
}

variable "app_ru_host" {
  type        = string
  description = "Hexlet Basics app ru-host"
}

variable "google_tag_manager_key" {
  type        = string
  description = "Google Tag Manager key"
}

# Cloudflare
variable "cloudflare_email" {
  type        = string
  description = "Cloudflare email"
}

variable "cloudflare_api_key" {
  type        = string
  description = "Cloudflare api key"
}

variable "domain" {
  type        = string
  description = "App domain"
}

variable "domain_ru" {
  type        = string
  description = "App domain"
}

variable "ip" {
  type        = string
  description = "App IP"
}
