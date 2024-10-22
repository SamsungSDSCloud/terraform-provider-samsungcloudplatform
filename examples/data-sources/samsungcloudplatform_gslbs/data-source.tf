data "samsungcloudplatform_gslbs" "my_scp_gslbs" {
}

output "contents" {
  value = data.samsungcloudplatform_gslbs.my_scp_gslbs.contents
}
