data "scp_region" "region" {
}

resource "scp_kubernetes_engine" "engine" {
  name               = var.name
  kubernetes_version = "v1.24.8"

  vpc_id            = data.terraform_remote_state.vpc.outputs.id
  subnet_id         = data.terraform_remote_state.subnet.outputs.id
  security_group_id = data.terraform_remote_state.security-group.outputs.id
  volume_id         = data.terraform_remote_state.file-storage.outputs.id

  // update optional field
  cloud_logging_enabled = false
  public_acl_ip_address = "123.123.123.123"
  load_balancer_id      = data.terraform_remote_state.load_balancer.outputs.id
  cifs_volume_id    = data.terraform_remote_state.file-storage.outputs.cifs_id
}



