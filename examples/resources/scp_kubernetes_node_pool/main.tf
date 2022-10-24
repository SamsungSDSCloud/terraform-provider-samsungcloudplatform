data "scp_region" "region" {
}

data "scp_standard_image" "ubuntu_image" {
  service_group = "CONTAINER"
  service       = "Kubernetes Engine VM"
  region        = data.scp_region.region.location

  filter {
    name      = "image_name"
    values    = ["Ubuntu 18.04 *"]
    use_regex = true
  }
}

resource "scp_kubernetes_node_pool" "pool" {
  name               = var.name
  engine_id          = data.terraform_remote_state.engine.outputs.id
  image_id           = data.scp_standard_image.ubuntu_image.id
  desired_node_count = 0
  storage_size_gb    = 100
  cpu_count          = 2
  memory_size_gb     = 4

  auto_recovery      = true
  auto_scale         = true
  min_node_count     = 1
  max_node_count     = 4
}
