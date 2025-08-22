resource "yandex_mdb_postgresql_cluster" "code_basics_cluster" {
  name                = "code-basics-postgresql-cluster"
  environment         = "PRODUCTION"
  network_id          = data.yandex_vpc_network.hexlet.id
  deletion_protection = true
  security_group_ids  = [yandex_vpc_security_group.code_basics_postgresql.id]

  config {
    version = 17

    resources {
      resource_preset_id = "c3-c2-m4"
      disk_type_id       = "network-ssd"
      disk_size          = 40
    }

    backup_window_start {
      hours   = 1
      minutes = 0
    }

    access {
      web_sql = false
      data_lens = true
    }

    performance_diagnostics {
      enabled                      = true
      sessions_sampling_interval   = 60
      statements_sampling_interval = 600
    }
  }

  host {
    name      = "code_basics_pg_host_b"
    zone      = data.yandex_vpc_subnet.hexlet_b.zone
    subnet_id = data.yandex_vpc_subnet.hexlet_b.id
  }

  host {
    name      = "code_basics_pg_host_a"
    zone      = data.yandex_vpc_subnet.hexlet_a.zone
    subnet_id = data.yandex_vpc_subnet.hexlet_a.id
  }

  host {
    name      = "code_basics_pg_host_d"
    zone      = data.yandex_vpc_subnet.hexlet_d.zone
    subnet_id = data.yandex_vpc_subnet.hexlet_d.id
  }

  maintenance_window {
    type = "WEEKLY"
    day  = "MON"
    hour = 2
  }
}

resource "yandex_mdb_postgresql_user" "user" {
  cluster_id = yandex_mdb_postgresql_cluster.code_basics_cluster.id
  name       = local.data.secrets.code_basics.database.user.name
  password   = local.data.secrets.code_basics.database.user.password
  conn_limit = 100
}

resource "yandex_mdb_postgresql_database" "hexlet_basics" {
  cluster_id = yandex_mdb_postgresql_cluster.code_basics_cluster.id
  name       = local.data.secrets.code_basics.database.name
  owner      = yandex_mdb_postgresql_user.user.name

  lc_collate = "en_US.UTF-8"
  lc_type    = "en_US.UTF-8"

  extension {
    name = "pg_repack"
  }
}

# NOTE: Yandex не умеет в создание readonly юзеров, поэтому тут делаем пользователя, и уже ручками делаем ему доступ в readonly
resource "yandex_mdb_postgresql_user" "readonly_user" {
  cluster_id = yandex_mdb_postgresql_cluster.code_basics_cluster.id
  name       = local.data.secrets.code_basics.database.readonly_user.name
  password   = local.data.secrets.code_basics.database.readonly_user.password

  permission {
    database_name = local.data.secrets.code_basics.database.name
  }
}
