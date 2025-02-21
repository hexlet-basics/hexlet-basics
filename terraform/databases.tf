data "twc_database_preset" "postgresql" {
  location = var.location

  type = "postgres"
  cpu = 1
  ram = 2048

  price_filter {
    from = 400
    to   = 600
  }
}

resource "twc_database_cluster" "postgresql" {
  name = "Hexlet basics PostgreSQL cluster"

  project_id = twc_project.hexlet_basics.id
  preset_id = data.twc_database_preset.postgresql.id

  availability_zone = var.zone

  network {
    id = twc_vpc.hexlet_basics.id
  }

  # TODO: убрать после переноса и удалить IP
  is_external_ip = true

  type = "postgres16"
}

resource "twc_database_backup_schedule" "postgresql" {
  cluster_id = twc_database_cluster.postgresql.id

  copy_count = 7
  creation_start_at = "2024-09-19T00:00:00.000Z"
  interval = "day"
  enabled = true
}

resource "twc_database_instance" "hexlet_basics" {
  cluster_id = twc_database_cluster.postgresql.id

  name = var.postgres_db.name
}

resource "twc_database_user" "hexlet_basics_user" {
  cluster_id = twc_database_cluster.postgresql.id

  login = var.postgres_db.username
  password = var.postgres_db.password

  instance {
    instance_id = twc_database_instance.hexlet_basics.id
    privileges  = ["SELECT", "INSERT", "UPDATE", "DELETE", "ALTER", "REFERENCES", "CREATE", "DROP", "INDEX"]
  }
}

resource "twc_database_user" "hexlet_basics_user_readonly" {
  cluster_id = twc_database_cluster.postgresql.id

  login = var.postgres_db.username_readonly
  password = var.postgres_db.password

  instance {
    instance_id = twc_database_instance.hexlet_basics.id
    privileges  = ["SELECT"]
  }
}

resource "twc_database_user" "datadog_user_readonly" {
  cluster_id = twc_database_cluster.postgresql.id

  login = var.postgres_db.datadog_username
  password = var.postgres_db.datadog_password

  instance {
    instance_id = twc_database_instance.hexlet_basics.id
    privileges  = ["SELECT"]
  }
}
