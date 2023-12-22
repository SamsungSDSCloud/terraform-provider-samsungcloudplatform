# Find all vpcs for current project
data "scp_region" "my_region" {
}
data "scp_vpcs" "vpcs" {
}

###VPC_DNS
resource "scp_endpoint" "my_endpoint" {
  ip_address   = "192.166.0.6"
  name        = var.name
  type        = "VPC_DNS"
  object_id   = data.scp_vpcs.vpcs.contents[0].vpc_id
  description = "Vpc Endpoint generated from Terraform"
  region      = data.scp_region.my_region.location
  vpc_id      = data.scp_vpcs.vpcs.contents[0].vpc_id
}
