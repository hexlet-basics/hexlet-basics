data "yandex_kubernetes_cluster" "kube_cluster_1" {
  cluster_id = local.data.terraform.yc.kube_cluster_id
}

resource "yandex_kubernetes_node_group" "code_basics_node_group_2" {
  cluster_id = data.yandex_kubernetes_cluster.kube_cluster_1.id
  name       = "code-basics-node-group-2"
  version    = data.yandex_kubernetes_cluster.kube_cluster_1.master[0].version

  node_labels = {
    "group" = "codebasics"
  }

  instance_template {
    platform_id = "standard-v3"

    network_interface {
      nat = false
      subnet_ids = [
        yandex_vpc_subnet.code_basics_b_1.id,
      ]
    }

    resources {
      memory = 4
      cores  = 2
    }

    boot_disk {
      type = "network-ssd-nonreplicated"
      size = 93
    }

    container_runtime {
      type = "containerd"
    }
  }

  scale_policy {
    fixed_scale {
      size = 3
    }
  }

  deploy_policy {
    max_expansion   = 1
    max_unavailable = 0
  }

  maintenance_policy {
    auto_upgrade = true
    auto_repair  = true
  }
}
