data "twc_k8s_preset" "hexlet_basics_master" {
  cpu = 2
  type = "master"
}

resource "twc_k8s_cluster" "hexlet_basics" {
  name = "Hexlet basics k8s cluster"

  project_id = twc_project.hexlet_basics.id

  high_availability = false
  version = "v1.30.2+k0s.0"
  network_driver = "calico"
  ingress = true

  preset_id = data.twc_k8s_preset.hexlet_basics_master.id
}

data "twc_k8s_preset" "hexlet_basics_worker" {
  cpu = 2
  ram = 4096
  type = "worker"
}

resource "twc_k8s_node_group" "hexlet_basics" {
  cluster_id = twc_k8s_cluster.hexlet_basics.id
  name = "hexlet-basics-node-group"

  preset_id = data.twc_k8s_preset.hexlet_basics_worker.id

  node_count = 1
}
