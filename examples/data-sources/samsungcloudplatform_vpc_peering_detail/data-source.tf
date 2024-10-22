data "samsungcloudplatform_vpc_peerings" "peering" {
}

data "samsungcloudplatform_vpc_peering_detail" "detail" {
   vpc_peering_id = data.samsungcloudplatform_vpc_peerings.peering.contents[0].vpc_peering_id
}

output "detail" {
  value = data.samsungcloudplatform_vpc_peering_detail.detail
}
