data "yandex_vpc_network" "hexlet" {
  network_id = local.data.terraform.yc_kz.network_id
}

data "yandex_vpc_route_table" "hexlet_a" {
  route_table_id = local.data.terraform.yc_kz.route_table_id
}

data "yandex_vpc_subnet" "db_subnet_a_1" {
  subnet_id = local.data.terraform.yc_kz.db_subnet_a_1_id
}

resource "yandex_vpc_subnet" "code_basics_a_1" {
  name           = "code-basics-subnet-a-1"
  v4_cidr_blocks = ["10.21.0.0/16"]
  zone           = local.data.terraform.yc_kz.zone
  network_id     = data.yandex_vpc_network.hexlet.id
  route_table_id = data.yandex_vpc_route_table.hexlet_a.id
}

resource "yandex_vpc_address" "code_basics_ingress_address_1" {
  name = "code-basics-ingress-adress-1"
  external_ipv4_address {
    zone_id = local.data.terraform.yc_kz.zone
  }
}

resource "yandex_vpc_security_group" "code_basics_postgresql" {
  name       = "code-basics-postgresql"
  network_id = data.yandex_vpc_network.hexlet.id

  ingress {
    description    = "Permit access to k8s nodes"
    protocol       = "TCP"
    port           = 6432
    v4_cidr_blocks = yandex_vpc_subnet.code_basics_a_1.v4_cidr_blocks
  }
}
