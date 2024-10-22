data "samsungcloudplatform_sqlservers" "my_scp_sqlservers" {
}

output "output_my_scp_sqlservers" {
  value = data.samsungcloudplatform_sqlservers.my_scp_sqlservers
}
