resource "scp_subnet" "my_subnet" {
  vpc_id      = data.terraform_remote_state.vpc.outputs.id
  name        = var.name
  type        = "PUBLIC"
  cidr_ipv4   = "192.169.4.0/24"
  description = "Subnet generated from Terraform"
}

