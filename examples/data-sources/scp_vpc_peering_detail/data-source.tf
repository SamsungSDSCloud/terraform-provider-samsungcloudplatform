data "scp_vpc_peerings" "peering" {
}

data "scp_vpc_peering_detail" "detail" {
   vpc_peering_id = data.scp_vpc_peerings.peering.contents[0].vpc_peering_id
}

output "detail" {
  value = data.scp_vpc_peering_detail.detail
}
