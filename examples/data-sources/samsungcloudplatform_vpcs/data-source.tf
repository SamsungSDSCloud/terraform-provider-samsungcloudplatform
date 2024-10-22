data "samsungcloudplatform_vpcs" "vpcs" {
}

output "contents" {
  value = data.samsungcloudplatform_vpcs.vpcs.contents
}
