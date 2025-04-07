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

variable "victoriametrics" {
  type = object({
    url = string
    username = string
    password = string
  })
}

variable "k8s_data" {
  type = object({
    ip1 = string
    ip2 = string
    ip3 = string
  })
}

variable "yc" {
  type = object({
    cloud_id = string
    folder_id = string
    zone = string
    storage = map(any)
  })
}
