resource "twc_vpc" "hexlet_basics" {
  name = "Hexlet basics network"
  subnet_v4 = "10.3.1.0/24"
  location = var.location
}

resource "twc_firewall" "databases" {
  name = "Hexlet basic DB firewall"

  # NOTE: timeweb terraform provider doesn't support linking database resources to firewall
  # therefore databases were added manually from control panel
  # link {
  #   id = resource.twc_database_cluster.postgresql.id
  #   type = "server"
  # }

  # link {
  #   id = resource.twc_database_cluster.redis.id
  #   type = "server"
  # }
}

resource "twc_firewall_rule" "own_subnet" {
  firewall_id = resource.twc_firewall.databases.id

  direction = "ingress"
  protocol = "tcp"
  cidr = twc_vpc.hexlet_basics.subnet_v4
}

resource "twc_firewall_rule" "timeweb_monitoring" {
  firewall_id = resource.twc_firewall.databases.id

  direction = "ingress"
  port = 10050
  protocol = "tcp"
  cidr = "92.53.116.0/24"
}


resource "twc_firewall_rule" "out_tcp" {
  firewall_id = resource.twc_firewall.databases.id

  direction = "egress"
  protocol = "tcp"
  cidr = "0.0.0.0/0"
}
