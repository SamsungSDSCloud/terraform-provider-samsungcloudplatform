data "scp_region" "region" {
}

data "scp_standard_image" "centos_image" {
  service_group = "COMPUTE"
  service       = "Virtual Server"
  region        = data.scp_region.region.location
  filter {
    name   = "image_name"
    values = ["CentOS 7.2"]
  }
}

resource "scp_virtual_server" "server_001" {
  name_prefix     = var.name
  admin_account   = var.id
  admin_password  = var.password
  cpu_count       = var.cpu
  memory_size_gb  = var.memory
  image_id        = data.scp_standard_image.centos_image.id
  vpc_id          = data.terraform_remote_state.vpc.outputs.id
  subnet_id       = data.terraform_remote_state.subnet.outputs.id

  delete_protection = false
  timezone          = "Asia/Seoul"
  contract_discount = "None"

  os_storage_name      = "hellodisk1"
  os_storage_size_gb   = 100
  os_storage_encrypted = false

  initial_script_content = "/test"

  security_group_ids = [
    data.terraform_remote_state.security_group.outputs.id
  ]
  use_dns = false

  external_storage {
    name            = var.ext_name
    product_name    = "SSD"
    storage_size_gb = 10
    encrypted       = false
  }
}
