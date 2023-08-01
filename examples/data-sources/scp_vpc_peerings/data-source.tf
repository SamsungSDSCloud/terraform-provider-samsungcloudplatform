data "scp_vpc_peerings" "peering" {
}

output "contents" {
  value = data.scp_vpc_peerings.peering.contents
}
