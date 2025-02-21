output "bounces_hostname" {
  value = resource.cloudflare_dns_record.bounces.name
}

output "hb_s3_access_key" {
  sensitive = true
  value     = yandex_iam_service_account_static_access_key.hb_s3_sa_static_key.access_key
}

output "hb_s3_secret_key" {
  sensitive = true
  value     = yandex_iam_service_account_static_access_key.hb_s3_sa_static_key.secret_key
}
