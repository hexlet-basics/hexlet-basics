variable "access_key" {
  type = string
}

variable "secret_key" {
  type = string
}

variable "location" {
  type = string
  default = "ru-1"
}

variable "zone" {
  type = string
  default = "spb-3"
}

variable "postgres_db" {
  type = map(any)
}

variable "environment_file" {
  type    = map(any)
}

variable "datadog_api_key" {
  type = string
}

variable "do_spaces_access_id" {
  type = string
  description = "Digital Ocean Spaces access id"
}

variable "do_spaces_secret_key" {
  type = string
  description = "Digital Ocean Spaces secret key"
}

variable "do_spaces_region" {
  type = string
  description = "Digital Ocean Spaces region"
}

variable "do_spaces_sitemap_bucket" {
  type = string
  description = "Digital Ocean Spaces bucket name for sitemap"
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

variable "k8s_data" {
  type = object({
    ip1 = string
    ip2 = string
    ip3 = string
  })
}
