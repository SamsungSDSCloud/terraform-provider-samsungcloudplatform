data "samsungcloudplatform_endpoints" "my_scp_endpoints" {
}

output "contents" {
  value = data.samsungcloudplatform_endpoints.my_scp_endpoints.contents
}
