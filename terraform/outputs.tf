output "postgresql_fqdn_current_master" {
  value = "c-${yandex_mdb_postgresql_cluster.code_basics_cluster.id}.rw.mdb.yandexcloud.net"
}

output "postgresql_fqdn_least_lagging_replica" {
  value = "c-${yandex_mdb_postgresql_cluster.code_basics_cluster.id}.ro.mdb.yandexcloud.net"
}

output "kube_cluster_id" {
  value = data.yandex_kubernetes_cluster.kube_cluster_1.id
}

output "kubernetes_nodes_subnet_id" {
  value = yandex_vpc_subnet.code_basics_b_1.id
}

output "ingress_ipv4" {
  value = yandex_vpc_address.code_basics_ingress_address_1.external_ipv4_address.0.address
}

output "cert_code_basics_com_name" {
  value = "yc-certmgr-cert-id-${yandex_cm_certificate.code_basics_com.id}"
}
