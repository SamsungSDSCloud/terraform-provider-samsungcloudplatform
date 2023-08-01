data "scp_vpc_peerings" "peerings" {
}

resource "scp_vpc_peering_cancel" "cancel" {
  vpc_peering_id = "VPC_PEERING-XXXX"
}
