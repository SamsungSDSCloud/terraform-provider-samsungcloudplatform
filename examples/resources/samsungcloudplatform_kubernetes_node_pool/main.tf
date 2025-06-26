data "samsungcloudplatform_region" "region" {
}

data "samsungcloudplatform_standard_image" "ubuntu_image" {
  service_group = "CONTAINER"
  service       = "Kubernetes Engine VM"
  region        = data.samsungcloudplatform_region.region.location

  filter {
    name      = "image_name"
    values    = ["Ubuntu 22.04 (Kubernetes)-v1.31.8"]
    use_regex = false
  }
}

resource "samsungcloudplatform_kubernetes_node_pool" "pool" {
  name               = var.name
  engine_id          = data.terraform_remote_state.engine.outputs.id
  image_id           = data.samsungcloudplatform_standard_image.ubuntu_image.id
  desired_node_count = 2

  scale_name = var.scale_name
  storage_name = var.storage_name
  storage_size_gb = var.storage_size_gb

  // optional field
  availability_zone_name = null

  // update optional field
  auto_scale         = false
  min_node_count     = null
  max_node_count     = null
  auto_recovery      = false
  labels {
    key = "test"
    value = "test"
  }
  taints {
    effect = "NoSchedule"
    key = "test"
    value = "test"
  }
}
