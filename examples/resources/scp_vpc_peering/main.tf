data "scp_vpcs" "vpcs" {
}

resource "scp_vpc_peering" "peering01" {
  approver_vpc_id         = "VPC-XXXX"
  firewall_enabled        = false
  requester_vpc_id        = "VPC-XXXX"
  vpc_peering_description = "Peering by Terraform"
}
