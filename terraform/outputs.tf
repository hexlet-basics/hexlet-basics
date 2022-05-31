output "cluster_name" {
  value = resource.digitalocean_kubernetes_cluster.hexlet_basics_cluster_2.name
}

output "postgres_db_host" {
  value = data.digitalocean_database_cluster.postgres_db_data.host
}

output "bounces_hostname" {
  value = resource.cloudflare_record.bounces.hostname
}

# output "main_ru_hostname" {
#   value = resource.cloudflare_record.main-ru.hostname
# }
