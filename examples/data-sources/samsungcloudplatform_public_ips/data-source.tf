# Find public ip list
data "samsungcloudplatform_public_ips" "my_scp_public_ips" {
}

output "contents" {
  value = data.samsungcloudplatform_public_ips.my_scp_public_ips.contents
}
