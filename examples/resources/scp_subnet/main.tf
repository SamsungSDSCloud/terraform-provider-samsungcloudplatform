# Find all vpcs for current project
data "scp_vpcs" "vpcs" {
}

resource "scp_subnet" "my_subnet" {
  vpc_id      = data.scp_vpcs.vpcs.contents[0].vpc_id
  name        = var.name
  type        = "PUBLIC"
  cidr_ipv4   = "192.169.4.0/24"
  description = var.description
}
