data "samsungcloudplatform_mysqls" "my_scp_mysqls" {
}

output "output_my_scp_mysqls" {
  value = data.samsungcloudplatform_mysqls.my_scp_mysqls
}
