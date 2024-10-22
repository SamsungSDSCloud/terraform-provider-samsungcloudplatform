data "samsungcloudplatform_vpc_peerings" "peerings" {
}

resource "samsungcloudplatform_vpc_peering_reject" "reject" {
  vpc_peering_id = "VPC_PEERING-XXXX"
}
