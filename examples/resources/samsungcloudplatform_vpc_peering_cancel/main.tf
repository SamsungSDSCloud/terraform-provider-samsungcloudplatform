data "samsungcloudplatform_vpc_peerings" "peerings" {
}

resource "samsungcloudplatform_vpc_peering_cancel" "cancel" {
  vpc_peering_id = "VPC_PEERING-XXXX"
}
