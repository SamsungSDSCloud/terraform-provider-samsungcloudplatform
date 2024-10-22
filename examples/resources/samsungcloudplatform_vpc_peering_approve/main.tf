data "samsungcloudplatform_vpc_peerings" "peerings" {
}

resource "samsungcloudplatform_vpc_peering_approve" "approve" {
  vpc_peering_id   = "VPC_PEERING-XXXX"
  firewall_enabled = false
}
