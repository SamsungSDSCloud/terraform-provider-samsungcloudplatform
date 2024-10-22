resource "samsungcloudplatform_load_balancer" "my_lb" {
  vpc_id      = data.terraform_remote_state.vpc.outputs.id
  name        = var.name
  description = "LoadBalancer generated from Terraform"
  size        = "SMALL"
  cidr_ipv4   = "192.168.102.0/24"
}
