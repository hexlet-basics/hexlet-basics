data "yandex_vpc_network" "hexlet" {
  network_id = local.data.terraform.yc.network_id
}

data "yandex_vpc_route_table" "hexlet_b" {
  route_table_id = local.data.terraform.yc.route_table_id
}

data "yandex_vpc_subnet" "hexlet_a" {
  subnet_id = local.data.terraform.yc.subnet_a_id
}

data "yandex_vpc_subnet" "hexlet_b" {
  subnet_id = local.data.terraform.yc.subnet_b_id
}

data "yandex_vpc_subnet" "hexlet_d" {
  subnet_id = local.data.terraform.yc.subnet_d_id
}

resource "yandex_vpc_subnet" "code_basics_b_1" {
  name           = "code-basics-subnet-b-1"
  v4_cidr_blocks = ["10.21.0.0/16"]
  zone           = local.data.terraform.yc.zone
  network_id     = data.yandex_vpc_network.hexlet.id
  route_table_id = data.yandex_vpc_route_table.hexlet_b.id
}

resource "yandex_vpc_address" "code_basics_ingress_address_1" {
  name = "code-basics-ingress-adress-1"
  external_ipv4_address {
    zone_id = local.data.terraform.yc.zone
  }
}

resource "yandex_vpc_security_group" "code_basics_postgresql" {
  name        = "code-basics-postgresql"
  network_id  = data.yandex_vpc_network.hexlet.id

  ingress {
    description       = "Permit access to k8s nodes"
    protocol          = "TCP"
    port              = 6432
    v4_cidr_blocks    = yandex_vpc_subnet.code_basics_b_1.v4_cidr_blocks
  }
}
