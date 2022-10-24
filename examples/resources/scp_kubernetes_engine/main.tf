data "scp_region" "region" {
}

resource "scp_kubernetes_engine" "engine" {
  name               = var.name
  kubernetes_version = "v1.21.8"

  vpc_id            = data.terraform_remote_state.vpc.outputs.id
  subnet_id         = data.terraform_remote_state.subnet.outputs.id
  security_group_id = data.terraform_remote_state.security-group.outputs.id
  volume_id         = data.terraform_remote_state.file-storage.outputs.id

  cloud_logging_enabled = false
  load_balancer_id      = data.terraform_remote_state.load_balancer.outputs.id
}
