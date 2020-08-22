variable "project_name" {
  default = "hexlet-basics"
}

variable "gke_cluster_name" {
  default = "hexlet-basics-5"
}

variable "gke_cluster_name_2" {
  default = "hexlet-basics-6"
}

variable "billing_account" {
  default = "01730C-0A5BE7-686C51"
}

variable "org_id" {
  default = "431020544079"
}

variable "region" {
  default = "europe-west3"
}

variable "zone" {
  default = "europe-west3-a"
}

variable "additional_zones" {
  default = ["europe-west3-a", "europe-west3-c"]
}

variable "db_name" {}
variable "db_username" {}
variable "db_password" {}
variable "db_hostname" {}
variable "db_port" {}
variable "github_oauth_token" {}
variable "cloudflare_api_key" {}
variable "cloudflare_email" {}
variable "github_client_id" {}
variable "github_client_secret" {}
variable "facebook_client_id" {}
variable "facebook_client_secret" {}
variable "secret_key_base" {}
variable "slack_codebuild_webhook" {}
variable "rollbar_access_token" {}
variable "sparkpost_smtp_username" {}
variable "sparkpost_smtp_password" {}
variable "guardian_secret_key" {}
variable "digitalocean_token" {}
variable "app_host" {}
variable "app_ru_host" {}
variable "app_scheme" {}

variable "repositories" {
  type = map

  default = {
    "exercises_php"        = "github_hexlet-basics_exercises-php"
    "exercises_java"       = "github_hexlet-basics_exercises-java"
    "exercises_javascript" = "github_hexlet-basics_exercises-javascript"
    "exercises_python"     = "github_hexlet-basics_exercises-python"
    "hexlet_basics"        = "github_hexlet-basics_hexlet_5Fbasics"
  }
}
