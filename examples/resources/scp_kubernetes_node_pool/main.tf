data "scp_region" "region" {
}

data "scp_standard_image" "ubuntu_image" {
  service_group = "CONTAINER"
  service       = "Kubernetes Engine VM"
  region        = data.scp_region.region.location

  filter {
    name      = "image_name"
    values    = ["Ubuntu 18.04 (Kubernetes)-v1.24.8"]
    use_regex = false
  }
}

resource "scp_kubernetes_node_pool" "pool" {
  name               = var.name
  engine_id          = data.terraform_remote_state.engine.outputs.id
  image_id           = data.scp_standard_image.ubuntu_image.id
  desired_node_count = 2
  cpu_count          = 2
  memory_size_gb     = 4
  storage_size_gb    = 100

  availability_zone_name = ""
  auto_recovery      = false
  auto_scale         = false
  min_node_count     = null
  max_node_count     = null
}
