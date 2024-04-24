data "scp_region" "region" {
}

data "scp_standard_image" "centos_image" {
  service_group = "COMPUTE"
  service       = "Virtual Server"
  region        = data.scp_region.region.location
  filter {
    name   = "image_name"
    values = ["CentOS 7.8"]
  }
}

resource "scp_virtual_server" "server_001" {
  virtual_server_name = var.name
  key_pair_id = data.terraform_remote_state.key_pair.outputs.id

  server_type = var.server-type
  image_id        = data.scp_standard_image.centos_image.id
  vpc_id          = data.terraform_remote_state.vpc.outputs.id
  subnet_id       = data.terraform_remote_state.subnet.outputs.id
  internal_ip_address = "192.169.4.17"

  delete_protection = false
  contract_discount = "None"

  os_storage_name      = "hellodisk1"
  os_storage_size_gb   = 100
  os_storage_encrypted = false

  initial_script_content = "/test"

  security_group_ids = [
    data.terraform_remote_state.security_group.outputs.id
  ]
  use_dns = false

  placement_group_id = data.terraform_remote_state.placement_group.outputs.id

  external_storage {
    name            = var.ext_name
    product_name    = "SSD"
    storage_size_gb = 10
    encrypted       = false
  }
}
