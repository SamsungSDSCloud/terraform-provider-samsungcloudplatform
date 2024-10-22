data "samsungcloudplatform_vpcs" "vpcs" {
}

resource "samsungcloudplatform_vpc_peering" "peering01" {
  approver_vpc_id         = "VPC-XXXX"
  firewall_enabled        = false
  requester_vpc_id        = "VPC-XXXX"
  vpc_peering_description = "Peering by Terraform"
}
