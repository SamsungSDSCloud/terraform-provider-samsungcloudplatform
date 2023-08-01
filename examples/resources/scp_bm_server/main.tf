data "scp_region" "region" {
}

data "scp_standard_images" "centos_image" {
  service_group = "COMPUTE"
  service       = "Baremetal Server"
  region        = data.scp_region.region.location
  filter {
    name   = "image_name"
    values = ["CentOS 7.8 *"]
    use_regex = true
  }
}

resource "scp_bm_server" "server_001" {
  bm_server_name = var.name
  admin_account   = var.id
  admin_password  = var.password
  cpu_count       = var.cpu
  memory_size_gb  = var.memory
  image_id        = data.scp_standard_images.centos_image.standard_images[0].id
  vpc_id          = data.terraform_remote_state.vpc.outputs.id
  subnet_id       = data.terraform_remote_state.subnet.outputs.id
  ipv4 = "192.169.4.3"

  delete_protection = false
  contract_discount = "None"
  initial_script = ""
  use_dns = false
  public_ip_id = ""
  use_hyper_threading = "N"
  nat_enabled = false
  local_subnet_enabled = false
  local_subnet_ipv4 = ""
}
