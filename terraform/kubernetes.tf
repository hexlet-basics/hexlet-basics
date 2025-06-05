data "twc_k8s_preset" "hexlet_basics_master" {
  cpu = 2
  type = "master"
}

data "twc_k8s_preset" "hexlet_basics_worker" {
  cpu = 2
  ram = 4096
  type = "worker"
}

resource "twc_k8s_cluster" "hexlet_basics_3" {
  name = "Hexlet basics k8s cluster 3"

  project_id = twc_project.hexlet_basics.id

  high_availability = false
  version = "v1.33.1+k0s.0"
  network_driver = "calico"
  ingress = true

  preset_id = data.twc_k8s_preset.hexlet_basics_master.id

  lifecycle {
    ignore_changes = [
      preset_id,
    ]
  }
}

resource "twc_k8s_node_group" "hexlet_basics_3" {
  cluster_id = twc_k8s_cluster.hexlet_basics_3.id
  name = "default"

  preset_id = data.twc_k8s_preset.hexlet_basics_worker.id

  node_count = 3
  is_autoscaling = true
  max_size = 3
  min_size = 3

  # NOTE: баг таймвеб, выдаёт не актуальный preset_id. Создано с id = 443
  lifecycle {
    ignore_changes = [
      preset_id,
    ]
  }
}
