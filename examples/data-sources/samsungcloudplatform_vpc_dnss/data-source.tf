data "samsungcloudplatform_vpcs" "vpcs" {
}

data "samsungcloudplatform_vpc_dnss" "dnss" {
  vpc_id = data.samsungcloudplatform_vpcs.vpcs.contents[0].vpc_id
}

output "contents" {
  value = data.samsungcloudplatform_vpc_dnss.dnss.contents
}
