data "samsungcloudplatform_vpc_peerings" "peering" {
}

output "contents" {
  value = data.samsungcloudplatform_vpc_peerings.peering.contents
}
