data "samsungcloudplatform_region" "region" {
}

data "samsungcloudplatform_virtual_servers" "virtual_server_list" {
  filter {
    name   = "virtual_server_state"
    values = ["RUNNING"]
  }
}

resource "samsungcloudplatform_kubernetes_engine" "engine" {
  name               = var.name
  kubernetes_version = "v1.31.8"

  vpc_id            = data.terraform_remote_state.vpc.outputs.id
  subnet_id         = data.terraform_remote_state.subnet.outputs.id
  security_group_id = data.terraform_remote_state.security-group.outputs.id
  volume_id         = data.terraform_remote_state.file-storage.outputs.id

  // update optional field
  cloud_logging_enabled = false
  public_acl_ip_address = "123.123.123.123"
  private_acl_resources {
    resource_id = data.samsungcloudplatform_virtual_servers.virtual_server_list.contents[0].virtual_server_id
    resource_type = "Virtual Server"
    resource_value = data.samsungcloudplatform_virtual_servers.virtual_server_list.contents[0].virtual_server_name
  }
  load_balancer_id      = data.terraform_remote_state.load_balancer.outputs.id
  cifs_volume_id    = data.terraform_remote_state.file-storage.outputs.cifs_id
}



