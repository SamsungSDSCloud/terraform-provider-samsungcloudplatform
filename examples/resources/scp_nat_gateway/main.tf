resource "scp_nat_gateway" "my_nat" {
  subnet_id    = data.terraform_remote_state.subnet.outputs.id
  public_ip_id = data.terraform_remote_state.public_ip.outputs.id
  description  = "NAT GW from Terraform"
}
