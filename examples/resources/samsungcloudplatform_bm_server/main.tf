data "samsungcloudplatform_region" "region" {
}

data "samsungcloudplatform_standard_images" "centos_image" {
  service_group = "COMPUTE"
  service       = "Baremetal Server"
  region        = data.samsungcloudplatform_region.region.location
  filter {
    name   = "image_name"
    values = ["CentOS 7.8 *"]
    use_regex = true
  }
}

resource "samsungcloudplatform_bm_server" "server_001" {
  admin_account   = var.id
  admin_password  = var.password
  cpu_count       = var.cpu
  memory_size_gb  = var.memory
  image_id        = data.samsungcloudplatform_standard_images.centos_image.standard_images[0].id
    vpc_id          = data.terraform_remote_state.vpc.outputs.id
    subnet_id       = data.terraform_remote_state.subnet.outputs.id
  delete_protection = false
  contract_discount = "None"
  initial_script = ""

  servers {
    bm_server_name = var.name
    ipv4 = var.ipv4
    use_dns = var.use_dns
    use_hyper_threading = var.hyper_threading
    nat_enabled = var.nat_enabled
    public_ip_id = var.public_ip_id
    local_subnet_enabled = var.local_subnet_enabled
    local_subnet_id = var.local_subnet_id
    local_subnet_ipv4 = var.local_subnet_ipv4
    state = var.state
  }

  timeouts {
    create = "30m"
    update = "30m"
    delete = "30m"
  }
}
