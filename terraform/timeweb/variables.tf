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
