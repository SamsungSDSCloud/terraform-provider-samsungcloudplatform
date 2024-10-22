data "samsungcloudplatform_subnets" "my_scp_subnets" {
}

output "contents" {
  value = data.samsungcloudplatform_subnets.my_scp_subnets.contents
}
